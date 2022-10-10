package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/tlgs/protohackers/internal/service"
)

type SmokeTest struct{ *service.Config }

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
	cfg := service.NewConfig(service.TCP, 10000)
	cfg.ParseFlags()

	service.Run(SmokeTest{cfg})
}
