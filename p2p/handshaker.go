package p2p

type HandshakeFunc func(any) error

func NOHandshakeFunc(any) error { return nil }
