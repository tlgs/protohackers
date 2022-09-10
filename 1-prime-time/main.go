package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

// using json.RawMessage is a workaround;
// json.Number happily slurps quoted numbers (https://github.com/golang/go/issues/34472),
// which this challenge does not accept as a valid number.
type request struct {
	Method string          `json:"method"`
	Number json.RawMessage `json:"number"`
}

func (m request) parseNumber() (any, error) {
	d, err := strconv.ParseInt(string(m.Number), 10, 64)
	if err == nil {
		return d, nil
	}

	f, err := strconv.ParseFloat(string(m.Number), 64)
	if err == nil {
		return f, nil
	}

	return nil, fmt.Errorf("could not parse %v as an Int64 or Float64", m.Number)
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

// Simple primality test
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

func protocol(raw []byte) []byte {
	var req request

	err := json.Unmarshal(raw, &req)
	if err != nil {
		// malformed JSON, sending malformed response
		return []byte(`{"reason":"malformed JSON"}` + "\n")
	}

	if req.Method != "isPrime" || req.Number == nil {
		// missing fields, sending malformed response
		return []byte(`{"reason":"missing fields"}` + "\n")
	}

	v, err := req.parseNumber()
	if err != nil {
		// number field is not a valid number, sending malformed response
		return []byte(`{"reason":"number field is not a valid number"}` + "\n")
	}

	var resp response
	switch n := v.(type) {
	case float64:
		resp = response{"isPrime", false}
	case int64:
		resp = response{"isPrime", isPrime(int(n))}
	}

	out, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return append(out, byte('\n'))
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
