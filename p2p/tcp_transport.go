package p2p

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type TCPPeer struct {
	conn net.Conn
	// if we dial a connection => outBound == true
	// if we accept a connection => outBound == false
	outBound bool
}

func NewTCPPeer(conn net.Conn, outBound bool) *TCPPeer {
	return &TCPPeer{conn: conn, outBound: outBound}
}

type TCPTransportOpts struct {
	ListenAddr string
	HandShaker HandshakeFunc
	Decoder    Decoder
}

// Mutex will protect the peers
type TCPTransport struct {
	opts     TCPTransportOpts
	listener net.Listener
	mutex    sync.RWMutex
	peers    map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		opts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	ln, err := net.Listen("tcp", t.opts.ListenAddr)
	if err != nil {
		return err
	}
	t.listener = ln
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stdout, []any{"TCP accept error: %s\n", err}...)
		}
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("Connection: %-v\n", peer)
	if err := t.opts.HandShaker(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}
	// read loop
	msg := &Message{}
	for {
		if err := t.opts.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		fmt.Printf("message: %+v\n", msg)
	}
}
