package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/tlgs/protohackers/internal/cli"
	"github.com/tlgs/protohackers/internal/server"
)

type SmokeTest struct{}

func (s SmokeTest) Setup() context.Context {
	return context.TODO()
}

func (s SmokeTest) Handle(_ context.Context, conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Println(err)
	}
}

func main() {
	config := cli.Parse()
	server.Run(SmokeTest{}, config.Port)
}
