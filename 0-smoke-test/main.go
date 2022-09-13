package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var port = flag.Int("p", 10000, "port to listen on")

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		log.Fatalln("error:", err)
	}

	log.Println("🚀 listening on port", *port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error:", err)
			continue
		}

		go func(conn net.Conn) {
			addr := conn.RemoteAddr()
			log.Println(addr, "accepted connection")

			r := io.TeeReader(conn, conn)
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				b := scanner.Bytes()
				log.Printf("%v echo: %#q", addr, b)
			}
			if err := scanner.Err(); err != nil {
				log.Println(addr, "error:", err)
			}

			conn.Close()
			log.Println(addr, "closing connection")
		}(conn)
	}
}
