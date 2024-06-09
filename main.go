package main

import (
	"log"

	"github.com/diegorezm/ds-fs-go/p2p"
)

func main() {
	addr := ":3000"
	tr := p2p.NewTCPTransport(addr)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	log.Printf("server running on http://localhost:%s\n", addr)
	select {}
}
