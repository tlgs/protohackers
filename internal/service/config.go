package service

import "flag"

type Protocol int

const (
	TCP Protocol = iota
	UDP
)

type Configuration interface {
	Protocol() Protocol
	Port() int

	ParseFlags()
}

// Config is the standard implementation of the Configuration interface
type Config struct {
	protocol Protocol
	port     int
}

func NewConfig(protocol Protocol, defaultPort int) *Config {
	return &Config{protocol: protocol, port: defaultPort}
}

func (cfg *Config) Port() int {
	return cfg.port
}

func (cfg *Config) Protocol() Protocol {
	return cfg.protocol
}

func (cfg *Config) ParseFlags() {
	flag.IntVar(&cfg.port, "port", cfg.port, "port to listen on")
	flag.Parse()
}
