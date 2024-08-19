package p2p

import (
	"encoding/gob"
	"io"
)

// Decoder defines an interface for decoding RPC messages from an io.Reader.
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

// GOBDecoder is an implementation of the Decoder interface that uses GOB encoding.
type GOBDecoder struct{}

// Decode decodes an RPC message from the given io.Reader using GOB encoding.
func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

// DefaultDecoder is an implementation of the Decoder interface that uses a custom decoding logic.
type DefaultDecoder struct{}

// Decode decodes an RPC message from the given io.Reader using custom logic.
// It first reads a byte to determine if the message is part of a stream.
// If it is a stream, it sets the Stream field of the RPC message to true and returns.
// Otherwise, it reads the payload into the RPC message.
func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	peerBuf := make([]byte, 1)
	if _, err := r.Read(peerBuf); err != nil {
		return err
	}

	// In case of stream we are not decoding what is being sent over the network.
	// We are just setting Stream to true so we can handle that in our logic.
	if peerBuf[0] == IncomingStream {
		msg.Stream = true
		return nil
	}

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]

	return nil
}
