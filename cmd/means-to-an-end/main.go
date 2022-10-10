package main

import (
	"context"
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/tlgs/protohackers/internal/service"
)

type MeansToAnEnd struct{ *service.Config }

func (s MeansToAnEnd) Setup() context.Context {
	return context.TODO()
}

func (s MeansToAnEnd) Handle(_ context.Context, conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Printf("accepted connection (%v)", addr)

	defer func() {
		conn.Close()
		log.Printf("closed connection (%v)", addr)
	}()

	asset := make(map[int32]int32)
	buf := make([]byte, 9)
	for {
		if _, err := io.ReadFull(conn, buf); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("%v (%v)", err, addr)
			break
		}

		fst := int32(binary.BigEndian.Uint32(buf[1:5]))
		snd := int32(binary.BigEndian.Uint32(buf[5:]))

		switch buf[0] {
		case 'I':
			asset[fst] = snd
			log.Printf("insert: %v %v (%v)", fst, snd, addr)

		case 'Q':
			var total, n, avg int
			for ts, price := range asset {
				if fst <= ts && ts <= snd {
					total += int(price)
					n += 1
				}
			}

			if n > 0 {
				avg = total / n
			}

			out := make([]byte, 4)
			binary.BigEndian.PutUint32(out, uint32(avg))

			if _, err := conn.Write(out); err != nil {
				log.Printf("%v (%v)", err, addr)
			} else {
				log.Printf("query: %v %v ⇒ %v (%v)", fst, snd, out, addr)
			}

		default:
			log.Printf("unknown operation type %v (%v)", buf[0], addr)
		}
	}
}

func main() {
	cfg := service.NewConfig(service.TCP, 10002)
	cfg.ParseFlags()

	service.Run(MeansToAnEnd{cfg})
}
