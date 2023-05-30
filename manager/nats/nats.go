package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

type Connector struct {
	nc *nats.Conn
}

var Gateway Connector

func Connect() {
	uri := os.Getenv("NATS_URI")
	if uri == "" {
		log.Fatal("NATS_URI not provided")
	}

	var err error
	Gateway.nc, err = nats.Connect(uri, nats.Timeout(5*time.Second))
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	fmt.Println("Connected to NATS at:", Gateway.nc.ConnectedUrl())
}

func Disconnect() error {
	return Gateway.nc.Drain()
}
