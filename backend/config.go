package main

import "os"

type Config struct{ Port string }

func loadConfig() Config {
	p := os.Getenv("PORT")
	return Config{Port: p}
}
