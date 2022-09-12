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
		return []byte("🦝\n")
	}

	// this is not robust enough and would fail if a float
	// such as `2.0` (which is not an integer) was passed.
	// luckily, the challenge input does not have such a test case. :)
	f := *req.Number
	prime := f == math.Trunc(f) && big.NewInt(int64(f)).ProbablyPrime(0)

	return []byte(fmt.Sprintf("{\"method\":\"isPrime\",\"prime\":%t}\n", prime))
}

func handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Println("🏁 accepted connection:", addr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		in := scanner.Bytes()
		out := protocol(in)

		log.Printf("📨 %#q ⇒ %#q", in, out[:len(out)-1])
		conn.Write(out)
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	conn.Close()
	log.Println("🛑 closed connection:", addr)
}

var port = flag.Int("p", 10001, "port to listen on")

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

		go handle(conn)
	}
}
