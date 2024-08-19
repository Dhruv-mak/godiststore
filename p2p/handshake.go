package p2p

// HandshakeFunc defines a function type for performing a handshake.
// It takes an argument of any type and returns an error if the handshake fails.
type HandshakeFunc func(any) error

// NOPHandshakeFunc is a no-operation handshake function that always succeeds.
// It takes an argument of any type and returns nil.
func NOPHandshakeFunc(any) error {
	return nil
}
