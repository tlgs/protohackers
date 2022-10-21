package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"

	"github.com/tlgs/protohackers/internal/protohackers"
)

var boguscoin = regexp.MustCompile(`^7[a-zA-Z0-9]{25,34}$`)

type MobInTheMiddle struct{ *protohackers.Config }

func (MobInTheMiddle) Setup() context.Context { return context.TODO() }

func (MobInTheMiddle) Handle(_ context.Context, downstream net.Conn) {
	addr := downstream.RemoteAddr()
	log.Printf("accepted connection: %v", addr)

	upstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		log.Printf("%s (%v)", err, addr)
		downstream.Close()
		return
	}

	// The following relay implementation follows a couple of ideas discussed in:
	// <https://stackoverflow.com/q/25090690/5818220>
	// <https://gist.github.com/jbardin/821d08cb64c01c84b81a>
	once := sync.Once{}
	relay := func(dst io.WriteCloser, src io.ReadCloser) {
		defer once.Do(func() { src.Close(); dst.Close(); log.Printf("closed connection: %v", addr) })

		for r := bufio.NewReader(src); ; {
			msg, err := r.ReadString('\n')
			if err != nil {
				return
			}

			tokens := make([]string, 0, 8)
			for _, raw := range strings.Split(msg[:len(msg)-1], " ") {
				t := boguscoin.ReplaceAllString(raw, "7YWHMfk9JZe0LM0g1ZauHuiSxhI")
				tokens = append(tokens, t)
			}

			out := strings.Join(tokens, " ") + "\n"
			if _, err = dst.Write([]byte(out)); err != nil {
				log.Printf("%s (%v)", err, addr)
			}
		}
	}

	go relay(downstream, upstream)
	relay(upstream, downstream)
}

func main() {
	cfg := protohackers.NewConfig(10005)
	cfg.ParseFlags()

	protohackers.RunTCP(MobInTheMiddle{cfg})
}
