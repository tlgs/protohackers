package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Println(addr, "accepted connection")

	asset := make(map[int32]int32)
	buf := make([]byte, 9)
	for {
		_, err := io.ReadFull(conn, buf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(addr, "error:", err)
			break
		}

		log.Println(addr, "received:", buf)

		fst := int32(binary.BigEndian.Uint32(buf[1:5]))
		snd := int32(binary.BigEndian.Uint32(buf[5:]))

		switch buf[0] {
		case 'I':
			asset[fst] = snd

		case 'Q':
			var total, n int
			for ts, price := range asset {
				if fst <= ts && ts <= snd {
					total += int(price)
					n += 1
				}
			}

			var avg int
			if n > 0 {
				avg = total / n
			}

			out := make([]byte, 4)
			binary.BigEndian.PutUint32(out, uint32(avg))

			_, err := conn.Write(out)
			if err != nil {
				log.Println(addr, "error:", err)
			} else {
				log.Println(addr, "sent:", out)
			}
		}
	}

	conn.Close()
	log.Println(addr, "closed connection")
}

var port = flag.Int("p", 10002, "port to listen on")

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

		go handle(conn)
	}
}
