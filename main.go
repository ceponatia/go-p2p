package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

const maxConnections = 2

type peer struct {
	connections []*net.UDPConn
	mux         sync.Mutex
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func CreatePeer(listenAddress, remoteAddress string) (*peer, error) {
	p := &peer{
		connections: make([]*net.UDPConn, 0, maxConnections),
	}

	localAddr, err := net.ResolveUDPAddr("udp", listenAddress)
	if err != nil {
		return nil, fmt.Errorf("error resolving local address: %w", err)
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("error binding UDP connection: %w", err)
	}

	p.AddConnection(conn)

	if remoteAddress != "" {
		err = p.CreatePeerAndConnect(remoteAddress)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *peer) CreatePeerAndConnect(address string) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.connections) >= maxConnections {
		return fmt.Errorf("maximum number of connections reached")
	}

	for _, conn := range p.connections {
		if conn.RemoteAddr().String() == address {
			return nil
		}
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("error resolving remote address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		return fmt.Errorf("error dialing UDP connection: %w", err)
	}

	p.AddConnection(conn)

	return nil
}

func (p *peer) AddConnection(conn *net.UDPConn) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.connections = append(p.connections, conn)
}

func (p *peer) SendMessage(connIndex int, message []byte) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	_, err := p.connections[connIndex].Write(message)
	return err
}

func (p *peer) ReceiveMessage(connIndex int, buf []byte) (int, *net.UDPAddr, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	n, remoteAddr, err := p.connections[connIndex].ReadFromUDP(buf)
	return n, remoteAddr, err
}

func main() {
	firstNode, err := CreatePeer("127.0.0.1:8080", "")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println("First peer is listening on", firstNode.connections[0].LocalAddr())

	secondNode, err := CreatePeer("127.0.0.1:8081", "127.0.0.1:8080")
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
}
