package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

type pricePoint struct {
	timestamp int32
	price     int32
}

func handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Println(addr, "🏁 accepted connection")

	var asset []pricePoint

	buf := make([]byte, 9)
	for {
		_, err := io.ReadFull(conn, buf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(addr, err)
			break
		}

		switch buf[0] {
		case byte('I'):
			log.Println(addr, "⬅️", buf[1:])

			pp := pricePoint{
				int32(binary.BigEndian.Uint32(buf[1:5])),
				int32(binary.BigEndian.Uint32(buf[5:])),
			}
			asset = append(asset, pp)

		case byte('Q'):
			minTime := int32(binary.BigEndian.Uint32(buf[1:5]))
			maxTime := int32(binary.BigEndian.Uint32(buf[5:]))

			total := 0.0
			n := 0
			for _, a := range asset {
				if minTime <= a.timestamp && a.timestamp <= maxTime {
					total += float64(a.price)
					n += 1
				}
			}

			var v int32
			if n > 0 {
				v = int32(total / float64(n))
			}

			out := make([]byte, 4)
			binary.BigEndian.PutUint32(out, uint32(v))

			log.Println(addr, "➡️", out)
			conn.Write(out)
		}
	}

	conn.Close()
	log.Println(addr, "🛑 closed connection")
}

var port = flag.Int("p", 10002, "port to listen on")

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("🚀 listening on port", *port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handle(conn)
	}
}
