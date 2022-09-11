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
		return []byte("x\n")
	}

	f := *req.Number
	if f == math.Trunc(f) && big.NewInt(int64(f)).ProbablyPrime(0) {
		return []byte(`{"method":"isPrime","prime":true}` + "\n")
	}

	return []byte(`{"method":"isPrime","prime":false}` + "\n")
}

func handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Println("accepted connection:", addr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		b := scanner.Bytes()
		out := protocol(b)

		log.Println(string(b), "⇒", string(out[:len(out)-1]))
		conn.Write(out)
	}

	conn.Close()
	log.Println("closed connection:", addr)
}

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

		go handle(conn)
	}
}
