package server

import (
	"log"
	"os"
	"strconv"
	"sync"
)

var once sync.Once
var conf ServerConfig

const defaultPort = "8080"

type ServerConfig struct {
	Port string
}

func Config() ServerConfig {
	once.Do(func() {
		port := os.Getenv("WALLET_SERVER_PORT")
		if port == "" {
			log.Printf("WALLET_SERVER_PORT env variable not specified, defaulting to %s", defaultPort)
			port = defaultPort
		}
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Printf("failed to parse port: %s, defaulting to %s", port, defaultPort)
			port = defaultPort
		}
		conf = ServerConfig{
			Port: port,
		}
	})
	return conf
}
