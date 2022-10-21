package main

import (
	"bufio"
	"bytes"
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

// This is a custom implementation of bufio.ScanLines so that the last
// non-empty line of input is not returned. Also, not stripping `\r`.
func customScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	return 0, nil, nil
}

type MobInTheMiddle struct{ *protohackers.Config }

func (MobInTheMiddle) Setup() context.Context { return context.TODO() }

func (MobInTheMiddle) Handle(_ context.Context, downstream net.Conn) {
	addr := downstream.RemoteAddr()
	log.Printf("accepted connection: %v", addr)

	upstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		log.Print(err)
		downstream.Close()
		return
	}

	// The following _patching proxy_ implementation follows a couple of ideas discussed in:
	// <https://stackoverflow.com/q/25090690/5818220>
	// <https://gist.github.com/jbardin/821d08cb64c01c84b81a>
	once := sync.Once{}
	patchedCopy := func(src io.ReadCloser, dst io.WriteCloser) {
		defer once.Do(
			func() {
				src.Close()
				dst.Close()
				log.Printf("closed connection: %v", addr)
			},
		)

		scanner := bufio.NewScanner(src)
		scanner.Split(customScanLines)

		for scanner.Scan() {
			parts := make([]string, 0)
			for _, part := range strings.Split(scanner.Text(), " ") {
				s := boguscoin.ReplaceAllString(part, "7YWHMfk9JZe0LM0g1ZauHuiSxhI")
				parts = append(parts, s)
			}

			out := strings.Join(parts, " ") + "\n"
			dst.Write([]byte(out))
		}
	}

	go patchedCopy(upstream, downstream)
	patchedCopy(downstream, upstream)
}

func main() {
	cfg := protohackers.NewConfig(10005)
	cfg.ParseFlags()

	protohackers.RunTCP(MobInTheMiddle{cfg})
}
