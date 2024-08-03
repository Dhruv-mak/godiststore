package p2p

// Peer is the interface that represents the remote node in the network
type Peer interface {
}

// Transport is the interface that handles the communication
// between the nodes in the network. This can be of the
// form (TCP, UDP, Websockets, etc)
type Transport interface {
	ListenAndAccept() error
}
