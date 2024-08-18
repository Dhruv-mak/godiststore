package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/Dhruv-mak/godiststore/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tCPTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: onpeer func
	}
	tcpTransport := p2p.NewTCPTransport(tCPTransportOpts)
	fileServerOpts := FileServerOpts{
		EncKey:            newEncryptionKey(),
		StorageRoot:       listenAddr + "_store",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}
	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	s3 := makeServer(":5000", ":3000", ":4000")

	
	go func() {log.Fatal(s1.Start())}()
	time.Sleep(1 * time.Second)
	go func() {log.Fatal(s2.Start())}()
	

	time.Sleep(2 * time.Second)

	go s3.Start()
	time.Sleep(2 * time.Second)

	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("picture_%d.jpg", i)
		data := bytes.NewReader([]byte("my big data file here!"))
		s2.Store(key, data)
	
		if err := s3.store.Delete(key)
	
		time.Sleep(5 * time.Millisecond); 
		if err != nil {
			log.Fatal(err)
		}
	
		r, err := s3.Get("coolPicture.jpg")
		if err != nil {
			log.Fatal(err)
		}
	
		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}

}
