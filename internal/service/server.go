package service

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

func Run(s Server) {
	port := s.Port()

	ctx := s.Setup()

	switch s.Protocol() {
	case TCP:
		ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("🚀 listening for TCP connections on port", port)

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			go s.Handle(ctx, conn)
		}

	case UDP:
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%v", port))
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("🚀 listening for UDP connections on port", port)

		s.Handle(ctx, conn)

	default:
		log.Fatal("unkown protocol")
	}
}
