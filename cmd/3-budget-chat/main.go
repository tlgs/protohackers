package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type ingressRequest struct {
	username string
	ch       chan string
	errc     chan error
}

type message struct {
	sender, content string
}

type Coordinator struct {
	ingress   chan ingressRequest
	egress    chan string
	broadcast chan message
}

func (c *Coordinator) loop() {
	users := make(map[string]chan string)
	for {
		select {
		case req := <-c.ingress:
			// validate username
			if _, exists := users[req.username]; exists {
				req.errc <- fmt.Errorf("requested username is taken")
				break
			} else if !isValidUsername(req.username) {
				req.errc <- fmt.Errorf("invalid username")
				break
			}

			// announce user, and collect existing users
			var existingUsers []string
			for name, ch := range users {
				ch <- fmt.Sprintf("* %v has entered the room\n", req.username)
				existingUsers = append(existingUsers, name)
			}

			// register user and list pre-existing users' names
			users[req.username] = req.ch
			req.errc <- nil
			req.ch <- fmt.Sprintf("* The room contains: %v\n", strings.Join(existingUsers, ", "))

		case username := <-c.egress:
			delete(users, username)
			for _, ch := range users {
				ch <- fmt.Sprintf("* %v has left the room\n", username)
			}

		case m := <-c.broadcast:
			for name, ch := range users {
				if m.sender == name {
					continue
				}
				ch <- fmt.Sprintf("[%v] %v\n", m.sender, m.content)
			}
		}
	}
}

func isValidUsername(s string) bool {
	matched, _ := regexp.MatchString(`^[[:alnum:]]{1,16}$`, s)
	return matched
}

func handle(conn net.Conn, c Coordinator) {
	addr := conn.RemoteAddr()
	log.Println(addr, "accepted connection")

	// cleanup in case we bail early
	defer func() {
		conn.Close()
		log.Println(addr, "closed connection")
	}()

	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	// first scanned line is the requested username of the client
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	username := scanner.Text()

	// attempt to join the chat room
	inc := make(chan string)
	errc := make(chan error)
	c.ingress <- ingressRequest{username, inc, errc}
	if err := <-errc; err != nil {
		log.Printf("%v %s: %v", addr, err, username)
		return
	}

	// setup a goroutine
	sch := make(chan string)
	go func() {
		for scanner.Scan() {
			sch <- scanner.Text()
		}
		close(sch)
	}()

	for {
		select {
		case v := <-inc:
			conn.Write([]byte(v))
		case s, ok := <-sch:
			if ok {
				c.broadcast <- message{username, s}
			} else {
				c.egress <- username
				return
			}
		}
	}
}

var port = flag.Int("p", 10003, "port to listen on")

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		log.Fatalln("error:", err)
	}

	coordinator := &Coordinator{
		ingress:   make(chan ingressRequest),
		egress:    make(chan string),
		broadcast: make(chan message),
	}
	go coordinator.loop()

	log.Println("🚀 listening on port", *port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error:", err)
			continue
		}

		go handle(conn, *coordinator)
	}
}
