package p2p

import "net"

// Peer is an interface that represents the remote node.
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport is anything that handles the communication
// between the nodes in the network. This can be of the
// form (TCP, UDP, websockets, ...)

// Transport defines an interface for network transport mechanisms in a peer-to-peer system.
type Transport interface {
	// Addr returns the address of the transport as a string.
	Addr() string

	// Dial establishes a connection to the given address.
	Dial(string) error

	// ListenAndAccept starts listening for incoming connections and accepts them.
	ListenAndAccept() error

	// Consume returns a channel from which RPC messages can be received.
	Consume() <-chan RPC

	// Close shuts down the transport and releases any resources associated with it.
	Close() error
}
