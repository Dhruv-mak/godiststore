package p2p

const (
	// IncomingMessage indicates that the message is a regular RPC message.
	IncomingMessage = 0x1
	// IncomingStream indicates that the message is part of a stream.
	IncomingStream = 0x2
)

// RPC holds any arbitrary data that is being sent over the transport between two nodes in the network.
type RPC struct {
	// From is the identifier of the sender.
	From string
	// Payload contains the data being sent.
	Payload []byte
	// Stream indicates whether the message is part of a stream.
	Stream bool
}
