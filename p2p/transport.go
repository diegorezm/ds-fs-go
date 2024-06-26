package p2p

// Interface that represents the remote node
type Peer interface {
	Close() error
}

// Communication between nodes (tcp,UDP, websockets,....)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
