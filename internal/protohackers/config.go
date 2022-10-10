package protohackers

import "flag"

type Configuration interface {
	Port() int

	ParseFlags()
}

// Config is the standard implementation of the Configuration interface
type Config struct {
	port int
}

func NewConfig(defaultPort int) *Config {
	return &Config{port: defaultPort}
}

func (cfg *Config) Port() int {
	return cfg.port
}

func (cfg *Config) ParseFlags() {
	flag.IntVar(&cfg.port, "port", cfg.port, "port to listen on")
	flag.Parse()
}
