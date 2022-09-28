package server

import (
	"context"
	"fmt"
	"log"
	"net"
)

type PHServer interface {
	Setup() context.Context
	Handle(ctx context.Context, conn net.Conn)
}

func Run(s PHServer, port int) {
	ctx := s.Setup()
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("🚀 listening on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.Handle(ctx, conn)
	}
}
