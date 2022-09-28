package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/tlgs/protohackers/internal/cli"
	"github.com/tlgs/protohackers/internal/server"
)

type BudgetChat struct{}

type Session struct {
	Username string
	Conn     net.Conn
	errc     chan error
}

type Message struct {
	Sender  string
	Content string
}

type Room struct {
	Ingress  chan Session
	Egress   chan string
	Messages chan Message
}

func Coordinator(room Room) {
	validUsername := regexp.MustCompile(`^[[:alnum:]]{1,16}$`)

	sessions := make(map[string]net.Conn)
	for {
		select {
		case s := <-room.Ingress:
			if _, exists := sessions[s.Username]; exists {
				s.errc <- fmt.Errorf("requested username is taken: " + s.Username)
				break
			} else if match := validUsername.MatchString(s.Username); !match {
				s.errc <- fmt.Errorf("invalid username: " + s.Username)
				break
			}

			var users []string
			msg := "* " + s.Username + " has entered the room\n"
			log.Print(msg)
			for name, conn := range sessions {
				conn.Write([]byte(msg))
				users = append(users, name)
			}

			sessions[s.Username] = s.Conn
			s.errc <- nil

			msg = "* The room contains: " + strings.Join(users, ", ") + "\n"
			s.Conn.Write([]byte(msg))

		case u := <-room.Egress:
			delete(sessions, u)

			msg := "* " + u + " has left the room\n"
			log.Print(msg)
			for _, conn := range sessions {
				conn.Write([]byte(msg))
			}

		case m := <-room.Messages:
			msg := "[" + m.Sender + "] " + m.Content + "\n"
			log.Print(msg)
			for name, conn := range sessions {
				if name != m.Sender {
					conn.Write([]byte(msg))
				}
			}
		}
	}
}

type ctxKey string

func (s BudgetChat) Setup() context.Context {
	room := Room{
		Ingress:  make(chan Session),
		Egress:   make(chan string),
		Messages: make(chan Message),
	}

	go Coordinator(room)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKey("room"), room)
	return ctx
}

func (s BudgetChat) Handle(ctx context.Context, conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Printf("accepted connection: %v", addr)
	defer func() {
		conn.Close()
		log.Printf("closed connection: %v", addr)
	}()

	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	r := ctx.Value(ctxKey("room")).(Room)

	var username string
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		username = scanner.Text()
		errc := make(chan error)
		r.Ingress <- Session{username, conn, errc}
		if err := <-errc; err != nil {
			log.Printf("%v (%v)", err, addr)
			return
		}

		defer func() {
			r.Egress <- username
		}()

	} else {
		log.Printf("no username provided (%v)", addr)
		return
	}

	for scanner.Scan() {
		r.Messages <- Message{username, scanner.Text()}
	}
}

func main() {
	config := cli.Parse()
	server.Run(BudgetChat{}, config.Port)
}