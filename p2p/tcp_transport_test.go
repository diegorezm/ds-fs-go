package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: ":4000",
		HandShaker: NOHandshakeFunc,
		Decoder:    DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.opts.ListenAddr, ":4000")
	assert.Nil(t, tr.ListenAndAccept())
}
