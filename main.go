package main

import (
	"fmt"
	"net"
)

type peer struct {
	connections []*net.UDPConn
}

func CreatePeerAndConnect(address string) (*peer, error) {
	// Create the first node
	p := &peer{
		connections: []*net.UDPConn{},
	}

	// Resolve remote address
	remoteAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("error resolving remote address: %w", err)
	}

	// Dial UDP connection
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("error dialing UDP connection: %w", err)
	}

	// Add connection to peer
	p.connections = append(p.connections, conn)

	return p, nil
}

func main() {

	firstNode := "localhost:8080"

	secondNode, err := CreatePeerAndConnect(firstNode)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Instantiate test message
	message := []byte("Hello World!")

	// Send message to first node
	_, err = secondNode.connections[0].Write(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	// First node may need time to receive the message
	// time.Sleep(1 * time.Second)
}
