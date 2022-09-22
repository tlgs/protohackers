package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type userChannel struct {
	mux sync.Mutex
	m   map[string]chan string
}

type message struct {
	sender  string
	content string
}

func isValidUsername(username string) bool {
	if len(username) < 1 {
		return false
	}

	for _, c := range username {
		if !(('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || ('0' <= c && c <= '9')) {
			return false
		}
	}

	return true
}

func handle(conn net.Conn, ch chan message) {
	addr := conn.RemoteAddr()
	log.Println(addr, "accepted connection")

	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	in := scanner.Text()
	log.Println(addr, fmt.Sprintf("%q tried to connect", in))

	users.mux.Lock()
	_, ok := users.m[in]
	users.mux.Unlock()

	if ok {
		log.Panicln(addr, "user already exists with this username")

		conn.Close()
		log.Println(addr, "closed connection")
		return
	}

	if !isValidUsername(in) {
		log.Println(addr, "bad username:", in)

		conn.Close()
		log.Println(addr, "closed connection")
		return
	}

	var existingUsers []string
	users.mux.Lock()
	for k := range users.m {
		existingUsers = append(existingUsers, k)
	}
	users.mux.Unlock()

	msg := fmt.Sprintf("* The room contains: %s\n", strings.Join(existingUsers, ", "))
	conn.Write([]byte(msg))

	inCh := make(chan string, 5)
	users.mux.Lock()
	users.m[in] = inCh
	users.mux.Unlock()
	ch <- message{in, fmt.Sprintf("* %v has entered the room\n", in)}

	scannerChan := make(chan string)
	quit := make(chan bool)
	go func() {
		for scanner.Scan() {
			scannerChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Println(addr, "error:", err)
		}

		close(scannerChan)
		quit <- true
	}()

	for flag := true; flag; {
		select {
		case v := <-scannerChan:
			// not sure what's up with this...
			if len(v) == 0 {
				break
			}

			ch <- message{in, fmt.Sprintf("[%v] %v\n", in, v)}
		case v := <-inCh:
			conn.Write([]byte(v))
		case <-quit:
			flag = false
		}
	}

	users.mux.Lock()
	delete(users.m, in)
	users.mux.Unlock()

	conn.Close()
	log.Println(addr, "closed connection")

	ch <- message{in, fmt.Sprintf("* %v has left the room\n", in)}
	close(inCh)
}

var port = flag.Int("p", 10003, "port to listen on")
var users = userChannel{m: make(map[string]chan string)}

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		log.Fatalln("error:", err)
	}

	broadcast := make(chan message)
	go func() {
		for {
			m := <-broadcast
			users.mux.Lock()
			for username, ch := range users.m {
				if username == m.sender {
					continue
				}
				ch <- m.content
			}
			users.mux.Unlock()
		}
	}()

	log.Println("🚀 listening on port", *port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error:", err)
			continue
		}

		go handle(conn, broadcast)
	}
}
