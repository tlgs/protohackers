package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/tlgs/protohackers/internal/protohackers"
)

type UnusualDatabaseProgram struct{ *protohackers.Config }

func (UnusualDatabaseProgram) Setup() context.Context { return context.TODO() }

func (UnusualDatabaseProgram) Handle(_ context.Context, conn net.Conn) {
	defer conn.Close()

	udpConn, ok := conn.(*net.UDPConn)
	if !ok {
		log.Println("bad connection type")
		return
	}

	var store = map[string]string{"version": "Ken's Key-Value Store 1.0.0"}

	buf := make([]byte, 1024)
	for {
		n, addr, err := udpConn.ReadFrom(buf)
		if err != nil {
			log.Print(err)
			continue
		}

		request := string(buf[:n])
		log.Printf("%v %v", addr, request)

		k, v, insert := strings.Cut(request, "=")
		if insert {
			if k == "version" {
				continue
			}
			store[k] = v

		} else {
			response := fmt.Sprintf("%v=%v", k, store[k])

			_, err := udpConn.WriteTo([]byte(response), addr)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func main() {
	cfg := protohackers.NewConfig(10004)
	cfg.ParseFlags()

	protohackers.RunUDP(UnusualDatabaseProgram{cfg})
}
