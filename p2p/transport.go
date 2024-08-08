package p2p

import "net"

// Peer is the interface that represents the remote node in the network
type Peer interface {
	RemoteAddr() net.Addr
	Send([]byte) error
	Close() error
}

// Transport is the interface that handles the communication
// between the nodes in the network. This can be of the
// form (TCP, UDP, Websockets, etc)
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
