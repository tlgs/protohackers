/*
This is a fairly straightforward exercise, especially because
a possible solution is part of the `net` package documentation:
<https://pkg.go.dev/net@go1.19.1#example-Listener>

And apparently its also given as an example on the challenges's
help page: <https://protohackers.com/help>
*/
package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var port = 10000

func handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Println("accepted connection:", addr)

	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Println(err)
	}

	log.Println("closing connection:", addr)
	conn.Close()
}

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("listening on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}

		go handle(conn)
	}
}
