package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"phoenixManager/api"
	"phoenixManager/nats"
)

func main() {
	// connect bridge
	nats.Connect()

	// start web api
	api.StartApi()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Sprintln("shutting down")
}
