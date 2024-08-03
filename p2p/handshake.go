package p2p

// HandshakeFunc... ?
type HandshakeFunc func(any) error

func NOPHandshakeFunc(any) error {
	return nil
}
