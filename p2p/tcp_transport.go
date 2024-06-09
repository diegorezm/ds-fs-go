package p2p

import (
	"fmt"
	"net"
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

// Mutex will protect the peers
type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handShaker    HandshakeFunc
	mutex         sync.RWMutex
	peers         map[net.Addr]Peer
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{listenAddress: addr, handShaker: NOHandshakeFunc}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	ln, err := net.Listen("tcp", t.listenAddress)
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
			fmt.Println("TCP accept error: %s\n", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("Connection: %+v\n", peer)
}
