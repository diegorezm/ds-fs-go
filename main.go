package main

import (
	"fmt"
	"log"

	"github.com/diegorezm/ds-fs-go/p2p"
)

func onPeer(p p2p.Peer) error {
	p.Close()
	return nil
}
func main() {
	addr := ":3000"
	tr := p2p.NewTCPTransport(p2p.TCPTransportOpts{
		ListenAddr: addr,
		HandShaker: p2p.NOHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
		OnPeer:     onPeer,
	})
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("msg: %+v\n", msg)
		}
	}()
	log.Printf("server running on http://localhost%s\n", addr)
	select {}
}
