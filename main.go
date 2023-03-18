package main

type peer struct {
  	conns []net.conn
}

func CreatePeerAndConnect(address string) peer {
  	p := peer {
    	connections: []net.conn {}
  	}
  	conn := // dial the address on udp
  	p.conns = append(p.conns, conn)
  
  	return p
}

func (p *peer) listen("udp", :"8080") {
  // initialise udp listener
	for {
		conn, err := ln.Accept()
    	if err != nil {
      		fmt.Println("Error accepting connection: ", err.Error())
			break
    	}
    go p.handleConnection(conn)
  }
}

func (p *peer) handleConnection(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 1024)
	msg, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading message: ", err.Error())
		break
	}
	
	segments := msg.split(' ')
	switch segments[0] { // command for rpc
	case "new":
	newAddress := segments[1] // payload, i.e. the address the new peer is listening on
    // tell a random connection to connect to newAddress
	case ...
	}
}