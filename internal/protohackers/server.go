package protohackers

import (
	"context"
	"fmt"
	"log"
	"net"
)

type Server interface {
	Configuration

	Protocol() string
	Setup() context.Context
	Handle(ctx context.Context, conn net.Conn)
}

func tcpService(s Server) {
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

func udpService(s Server) {
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

func Run(s Server) {
	switch s.Protocol() {
	case "tcp":
		tcpService(s)
	case "udp":
		udpService(s)
	default:
		log.Fatal("unknown protocol")
	}
}
