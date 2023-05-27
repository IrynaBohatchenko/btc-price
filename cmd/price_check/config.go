package main

import "time"

type Config struct {
	Port          string        `env:"PORT" envDefault:":8080"`
	ServerTimeout time.Duration `env:"SERVER_TIMEOUT" envDefault:"10s"`
	ClientTimeout time.Duration `env:"CLIENT_TIMEOUT" envDefault:"30s"`
}
