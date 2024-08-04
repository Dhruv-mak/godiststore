package main

import (
	"fmt"
	"log"

	"github.com/Dhruv-mak/godiststore/p2p"
)

func Onpeer(peer p2p.Peer) error {
	// fmt.Println("doing some logic with the peer outside the onpeer func")
	peer.Close()
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":5000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        Onpeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Received message: %s\n", msg.Payload)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
