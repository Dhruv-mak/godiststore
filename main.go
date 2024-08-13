package main

import (
	"bytes"
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

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(2 * time.Second)

	go s2.Start()
	time.Sleep(2 * time.Second)

	data := bytes.NewReader([]byte("my big data file here!"))
	s2.StoreData("myprivatedata", data)

	select {}
}
