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
		log.Fatalln(err)
	}

	log.Println("listening on port", *port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func(conn net.Conn) {
			addr := conn.RemoteAddr()
			log.Println("🏁 accepted connection:", addr)

			r := io.TeeReader(conn, conn)
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				b := scanner.Bytes()
				log.Printf("📨 %q", b)
			}
			if err := scanner.Err(); err != nil {
				log.Println(err)
			}

			log.Println("🛑 closing connection:", addr)
			conn.Close()
		}(conn)
	}
}
