package main

import (
	"fmt"
	"net"
	"os"
)

type peer struct {
	conn lAddr
}

func CreatePeer(lAddr) (*peer, error) {
	p := &peer{
		Port: lAddr,
	}

	lAddr, err := net.ResolveUDPAddr("udp", "8080")
	if err != nil {
		return nil, fmt.Errorf("error resolving local address: %w", err)
	}

	conn, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		return nil, fmt.Errorf("error binding UDP connection: %w", err)
	}

	return p, nil
}



func main() {
	firstNode, err := CreatePeer("127.0.0.1:8080")
	if err != nil {
		fmt.Errorf("%w", err)
		return
	}

	fmt.Println("First peer is listening on", firstNode.conn[0])

	secondNode, err := CreatePeer("127.0.0.1:8081")
	if err != nil {
		fmt.Errorf("%w", err)
		return
	}

	fmt.Println("Second peer is listening on", secondNode.conn[0])


}
