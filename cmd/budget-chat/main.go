package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/tlgs/protohackers/internal/protohackers"
)

type BudgetChat struct{ *protohackers.Config }

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
				if _, err := conn.Write([]byte(msg)); err != nil {
					log.Printf("%s (%v)", err, conn.RemoteAddr())
				}

				users = append(users, name)
			}

			sessions[s.Username] = s.Conn
			s.errc <- nil

			msg = "* The room contains: " + strings.Join(users, ", ") + "\n"
			if _, err := s.Conn.Write([]byte(msg)); err != nil {
				log.Printf("%s (%v)", err, s.Conn.RemoteAddr())
			}

		case u := <-room.Egress:
			delete(sessions, u)

			msg := "* " + u + " has left the room\n"
			log.Print(msg)
			for _, conn := range sessions {
				if _, err := conn.Write([]byte(msg)); err != nil {
					log.Printf("%s (%v)", err, conn.RemoteAddr())
				}
			}

		case m := <-room.Messages:
			msg := "[" + m.Sender + "] " + m.Content + "\n"
			log.Print(msg)
			for name, conn := range sessions {
				if name == m.Sender {
					continue
				}

				if _, err := conn.Write([]byte(msg)); err != nil {
					log.Printf("%s (%v)", err, conn.RemoteAddr())
				}
			}
		}
	}
}

type ctxKey string

func (BudgetChat) Setup() context.Context {
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

func (BudgetChat) Handle(ctx context.Context, conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Printf("accepted connection: %v", addr)
	defer func() {
		conn.Close()
		log.Printf("closed connection: %v", addr)
	}()

	msg := "Welcome to budgetchat! What shall I call you?\n"
	if _, err := conn.Write([]byte(msg)); err != nil {
		log.Printf("%s (%v)", err, addr)
	}

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
	cfg := protohackers.NewConfig(10003)
	cfg.ParseFlags()

	protohackers.RunTCP(BudgetChat{cfg})
}
