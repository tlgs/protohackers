package protohackers

import (
	"context"
	"fmt"
	"log"
	"net"
)

type Server interface {
	Configuration

	Setup() context.Context
	Handle(ctx context.Context, conn net.Conn)
}

// RunTCP implements the setup, listen, and accept loop steps for a
// Server intended to use the TCP protocol.
// Each TCP connection is handled by an individual goroutine.
func RunTCP(s Server) {
	port := s.Port()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("🚀 listening for TCP connections on port %v", port)

	ctx := s.Setup()
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.Handle(ctx, conn)
	}
}

// RunUDP implements the setup, and listen steps for a Server intended
// to use the UDP protocol.
//
// Because UDP is a connectionless protocol, each packet needs to be handled
// individually and it (probably) doesn't make much sense to spin up a goroutine
// for each arriving piece of data.
// As it stands, the current design is limited to handle each
// received packet sequentially in the program's main thread.
// If a future problem necessitates it, it should be easy to create a
// pool of N goroutines that share the *UDPConn object and can handle the
// incoming data concurrently.
func RunUDP(s Server) {
	port := s.Port()

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("🚀 listening for UDP connections on port %v", port)

	ctx := s.Setup()
	s.Handle(ctx, conn)
}
