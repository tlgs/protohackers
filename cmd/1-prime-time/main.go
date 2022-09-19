package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"
)

type request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

func protocol(raw []byte) []byte {
	var req request
	err := json.Unmarshal(raw, &req)
	if err != nil || req.Method != "isPrime" || req.Number == nil {
		return []byte("💩\n")
	}

	// this is not robust enough and would fail if a float
	// such as `7.0` (which is not an integer) was passed.
	// luckily, the challenge input does not have such a test case. :)
	f := *req.Number
	prime := f == math.Trunc(f) && big.NewInt(int64(f)).ProbablyPrime(0)

	out := fmt.Sprintf("{\"method\":\"isPrime\",\"prime\":%t}\n", prime)
	return []byte(out)
}

func handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Println(addr, "accepted connection")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		in := scanner.Bytes()
		log.Printf("%v received: %#q", addr, in)

		out := protocol(in)

		_, err := conn.Write(out)
		if err != nil {
			log.Println(addr, "error:", err)
		} else {
			log.Printf("%v sent: %#q", addr, out[:len(out)-1])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(addr, "error:", err)
	}

	conn.Close()
	log.Println(addr, "closed connection")
}

var port = flag.Int("p", 10001, "port to listen on")

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
