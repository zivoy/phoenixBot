package nats

import (
	"errors"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

type Connector struct {
	nc *nats.Conn
}

var Gateway = &Connector{}

func Connect() {
	uri := os.Getenv("NATS_URI")
	if uri == "" {
		log.Print("WARNING: NATS_URI not provided, not connecting")
		return
	}

	var err error
	Gateway.nc, err = nats.Connect(uri, nats.Timeout(5*time.Second), nats.Name("Manager API"), nats.PingInterval(1*time.Minute))
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	log.Println("Connected to NATS at:", Gateway.nc.ConnectedUrl())
}

func Disconnect() error {
	if Gateway == nil || Gateway.nc == nil {
		return nil
	}
	return Gateway.nc.Drain()
}

func (c *Connector) verifyConnected() error {
	if c != nil {
		return nil
	}

	//log.Print("Not connected to NATS")
	return errors.New("NATS not connected")
}
