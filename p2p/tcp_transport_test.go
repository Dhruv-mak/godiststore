package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":5000"
	tr := NewTCPTransport(listenAddr)

	assert.Equal(t, listenAddr, tr.listenAddress)

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
