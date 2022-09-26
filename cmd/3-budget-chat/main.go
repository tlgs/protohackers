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

type egressRequest string

type message struct {
	sender, content string
}

type ChatRoom struct {
	serverPort int

	ingress   chan ingressRequest
	egress    chan egressRequest
	broadcast chan message
}

func NewChatRoom(port int) ChatRoom {
	return ChatRoom{
		serverPort: port,
		ingress:    make(chan ingressRequest),
		egress:     make(chan egressRequest),
		broadcast:  make(chan message),
	}
}

func (cr ChatRoom) ServeForever() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", cr.serverPort))
	if err != nil {
		return err
	}
	log.Println("🚀 listening on port", cr.serverPort)

	// spawn a goroutine to listen for incoming network connections
	incoming := make(chan net.Conn)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("error:", err)
				continue
			}

			incoming <- conn
		}
	}()

	validUsername := regexp.MustCompile(`^[[:alnum:]]{1,16}$`)
	users := make(map[string]chan string)

	// main chat room control loop
	for {
		select {
		case conn := <-incoming:
			go cr.handle(conn)

		case req := <-cr.ingress:
			// validate username
			if _, exists := users[req.username]; exists {
				req.errc <- fmt.Errorf("requested username is taken")
				break
			} else if match := validUsername.MatchString(req.username); !match {
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

		case username := <-cr.egress:
			// remove user and announce departure
			delete(users, string(username))
			for _, ch := range users {
				ch <- fmt.Sprintf("* %v has left the room\n", username)
			}

		case msg := <-cr.broadcast:
			for name, ch := range users {
				if msg.sender == name {
					continue
				}
				ch <- fmt.Sprintf("[%v] %v\n", msg.sender, msg.content)
			}
		}
	}
}

func (cr ChatRoom) handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Println(addr, "accepted connection")

	// cleanup in case we bail early
	defer func() {
		conn.Close()
		log.Println(addr, "closed connection")
	}()

	// greet user and capture requested username
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	username := scanner.Text()

	// attempt to join the chat room
	inc := make(chan string)
	errc := make(chan error)
	cr.ingress <- ingressRequest{username, inc, errc}
	if err := <-errc; err != nil {
		log.Printf("%v %s: %v", addr, err, username)
		return
	}

	// setup a goroutine to capture user messages
	sch := make(chan string)
	go func() {
		for scanner.Scan() {
			sch <- scanner.Text()
		}
		close(sch)
	}()

	// main loop: write message to user connection, or broadcast a message from user
	for {
		select {
		case v := <-inc:
			conn.Write([]byte(v))
		case s, ok := <-sch:
			if !ok {
				cr.egress <- egressRequest(username)
				return
			}
			cr.broadcast <- message{username, s}
		}
	}
}

var port = flag.Int("p", 10003, "port to listen on")

func main() {
	flag.Parse()

	if err := NewChatRoom(*port).ServeForever(); err != nil {
		log.Fatalln("error:", err)
	}
}
