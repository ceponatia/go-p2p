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
