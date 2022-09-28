package cli

import "flag"

type Config struct {
	Port int
}

var port = flag.Int("port", 10000, "port to listen on")

func Parse() Config {
	flag.Parse()
	return Config{Port: *port}
}
