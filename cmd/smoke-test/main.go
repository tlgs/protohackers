package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/tlgs/protohackers/internal/protohackers"
)

type SmokeTest struct{ *protohackers.Config }

func (SmokeTest) Setup() context.Context { return context.TODO() }

func (SmokeTest) Handle(_ context.Context, conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Println(err)
	}
}

func main() {
	cfg := protohackers.NewConfig(10000)
	cfg.ParseFlags()

	protohackers.RunTCP(SmokeTest{cfg})
}
