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

// TCPPeer implements peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr string
	HandShaker HandshakeFunc
	Decoder    Decoder
	OnPeer     func(Peer) error
}

// Mutex will protect the peers
type TCPTransport struct {
	opts     TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
	mutex    sync.RWMutex
	peers    map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		opts:  opts,
		rpcch: make(chan RPC),
	}
}

// Consume implements the transport interface. Returns read-only channel
// for reading the incoming messages.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
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

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	peer := NewTCPPeer(conn, true)
	fmt.Printf("Connection: %-v\n", peer)
	if err := t.opts.HandShaker(peer); err != nil {
		handleError(err)
		conn.Close()
		return
	}
	if t.opts.OnPeer != nil {
		if err = t.opts.OnPeer(peer); err != nil {
			handleError(err)
			conn.Close()
			return
		}
	}
	// read loop
	msg := RPC{}
	for {
		err := t.opts.Decoder.Decode(conn, &msg)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Op == "read" && opErr.Err.Error() == "use of closed network connection" {
					handleError(opErr.Err)
					conn.Close()
					return
				}
			}
			fmt.Printf("TCP READ ERROR: %s\n", err)
			continue
		}

		msg.From = conn.RemoteAddr()
		fmt.Printf("message: %+v\n", msg)
		t.rpcch <- msg
	}
}

// TODO: this is horrible and i should change it
func handleError(err error, message ...string) {
	if len(message) > 0 {
		fmt.Printf("%s: %s\n", message[0], err)
	} else {
		fmt.Printf("ERROR: %s\n", err)
	}
}
