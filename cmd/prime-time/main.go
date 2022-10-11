package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"math"
	"net"

	"github.com/tlgs/protohackers/internal/protohackers"
)

type PrimeTime struct{ *protohackers.Config }

func (PrimeTime) Setup() context.Context { return context.TODO() }

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func isPrime(n int) bool {
	if n <= 3 {
		return n > 1
	} else if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}

func (PrimeTime) Handle(_ context.Context, conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Printf("accepted connection: %v", addr)

	defer func() {
		conn.Close()
		log.Printf("closed connection: %v", addr)
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		in := scanner.Bytes()

		var out []byte
		var req Request
		err := json.Unmarshal(in, &req)
		if err != nil || req.Method != "isPrime" || req.Number == nil {
			out = []byte("💩")
		} else {
			// this is not robust enough and would fail if a float
			// such as `7.0` (which is not an integer) was passed.
			// luckily, the challenge input does not have such a test case. :)
			f := *req.Number
			primality := f == math.Trunc(f) && isPrime(int(f))

			out, _ = json.Marshal(Response{"isPrime", primality})
		}
		out = append(out, byte('\n'))

		if _, err = conn.Write(out); err != nil {
			log.Printf("%v (%v)", err, addr)
		} else {
			log.Printf("%#q ⇒ %#q (%v)", in, out[:len(out)-1], addr)
		}
	}
}

func main() {
	cfg := protohackers.NewConfig(10001)
	cfg.ParseFlags()

	protohackers.RunTCP(PrimeTime{cfg})
}
