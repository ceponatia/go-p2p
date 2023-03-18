package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

type peer struct {
	connections []*net.UDPConn
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Creates a peer and either listens for incoming connections or connects to second peer
func CreatePeer(listenAddress, remoteAddress string) (*peer, error) {
	p := &peer{
		connections: []*net.UDPConn{},
	}

	// Resolve local address
	localAddr, err := net.ResolveUDPAddr("udp", listenAddress)
	if err != nil {
		return nil, fmt.Errorf("error resolving local address: %w", err)
	}

	// Bind UDP connection to localAddr
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("error binding UDP connection: %w", err)
	}

	// Add connection to peer
	p.connections = append(p.connections, conn)

	// If remoteAddress is not empty, connection to remote peer
	if remoteAddress != "" {
		err = p.CreatePeerAndConnect(remoteAddress)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *peer) CreatePeerAndConnect(address string) error {
	remoteAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		handleError(err)
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		handleError(err)
	}

	p.connections = append(p.connections, conn)
	return nil
}

func (p *peer) SendMessage(connIndex int, message []byte) error {
	_, err := p.connections[connIndex].Write(message)
	return err
}

func (p *peer) ReceiveMessage(connIndex int, buf []byte) (int, *net.UDPAddr, error) {
	n, remoteAddr, err := p.connections[connIndex].ReadFromUDP(buf)
	return n, remoteAddr, err
}

func main() {
	firstNode, err := CreatePeer("localhost:8080", "")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println("First peer is listening on", firstNode.connections[0].LocalAddr())

	secondNode, err := CreatePeer("localhost:8081", "localhost:8080")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println("Second peer is listening on", secondNode.connections[0].LocalAddr())

	go func() {
		buf := make([]byte, 1024)
		for {
			n, remoteAddr, err := firstNode.ReceiveMessage(0, buf)
			if err != nil {
				handleError(err)
				continue
			}

			message := string(buf[:n])
			fmt.Printf("First node received message from %s: %s\n", remoteAddr.String(), message)
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, remoteAddr, err := secondNode.ReceiveMessage(0, buf)
			if err != nil {
				handleError(err)
				continue
			}

			message := string(buf[:n])
			fmt.Printf("Second node received message from %s: %s\n", remoteAddr.String(), message)
		}
	}()

	time.Sleep(2 * time.Second)

	err = firstNode.SendMessage(0, []byte("Hello from first node"))
	if err != nil {
		handleError(err)
		return
	}

	err = secondNode.SendMessage(0, []byte("Hello from second node"))
	if err != nil {
		handleError(err)
		return
	}

	time.Sleep(2 * time.Second)
}
