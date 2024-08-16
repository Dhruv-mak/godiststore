package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/Dhruv-mak/godiststore/p2p"
)

type FileServerOpts struct {
	EncKey 		  []byte
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		store:          NewStore(storeOpts),
		FileServerOpts: opts,
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server loop stopped due to error or user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode((&msg)); err != nil {
				log.Println("decoding error: ", err)
				return
			}
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("handle message error: ", err)
			}

		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case *MessageStoreFile:
		return s.handleMessageStorefile(from, v)
	case *MessageGetFile:
		return s.handleMessageGetfile(from, v)
	}

	return nil
}

func (s *FileServer) handleMessageGetfile(from string, msg MessageGetFile) error {
	if !s.store.Has(msg.Key) {
		fmt.Printf("[%s] need to serve file (%s) but it does not exist on disk \n", s.Transport.Addr(), msg.Key)
	}

	fmt.Printf("[%s] serving file (%s) over the network\n",s.Transport.Addr(), msg.Key)
	fileSize, r, err := s.store.Read(msg.Key)
	if err != nil {
		return err
	}

	if rc, ok := r.(*io.ReadCloser); ok {
		fmt.Println("Closing readCloser")
		defer (*rc).Close()
	}

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) not found in peer map", from)
	}

	// First send the "incomingStream" byte to the peer and then we can send
	// the file size as an int64.
	peer.Send([]byte{p2p.IncomingStream})
	binary.Write(peer, binary.LittleEndian, fileSize)
	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] written (%d) bytes over the network to %s\n", s.Transport.Addr(), n, from)

	peer.CloseStream()

	return nil
}

func (s *FileServer) handleMessageStorefile(from string, msg MessagStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		fmt.Errorf("peer (%s) not found in peer map", from)
	}

	n, err := s.store.Write(msg.key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}

	fmt.Printf("written %d bytes to disk\n", n)
	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

type Payload struct {
	Key  string
	Data []byte
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

type MessageGetFile struct {
	Key string
}

func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Has(key) {
		fmt.Printf("[%s] serving file (%s) from local disk\n", s.Transport.Addr(), key)
		_, r, err :=  s.store.Read(key)
		return r, err
	}

	fmt.Printf("[%s] dont have file (%s) locally, fetching from network", s.Transport.Addr(), key)

	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 500)

	for _, peer := range s.peers {
		// First read the file size so we can limit the amount of bytes that we read
		// from the connection, so it will not keep hanging.
		var fileSize int64
		binary.Read(peer, binary.LittleEndian, &fileSize)
		n, err := s.store.Write(key, io.LimitReader(peer, fileSize))
		if err == nil {
			return nil, err
		}

		fmt.Printf("[%s] received %d bytes over the network from (%s)\n", s.Transport.Addr(), n, peer.RemoteAddr().String())

		peer.CloseStream()
	}

	_, r, err, := s.store.Read(key)
	return r, err
}

func (s *FileServer) stream(msg *Message) error {
	peers := []p2p.Peer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(...peers)
	return gob.NewEncoder(mw).Encode(msg)

}

func (s *FileSerever) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func (s *FileServer) Store(key string, r io.Reader) error {
	// 1. Store this file to disk
	// 2. broadcast this file to all known peers

	var (
		fileBuffer = new(bytes.Buffer)
		tee        = io.TeeReader(r, fileBuffer)
	)

	size, err := s.store.Write(key, r)
	if err != nil {
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size + 16,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 5)

	// TODO: use a multiwriter here
	for _, peer := range s.peers {
		peer.Send([]byte{p2p.IncomingStream})
		// n, err := io.Copy(peer, fileBuffer)
		// if err != nil {
		// 	return err
		// }
		n, err := copyEncrypt(s.EncKey, fileBuffer, peer)
		if err != nil {
			return errors
		}

		fmt.Println("received and returned bytes to disk: ", n)
	}
	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s", p.RemoteAddr().String())

	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
				// panic(err)
			}
		}(addr)

	}
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.loop()

	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
