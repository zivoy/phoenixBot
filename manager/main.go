package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"phoenixManager/files"
	"syscall"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	files.StopCleaning()
	api.Shutdown(ctx)
	log.Fatal(nats.Disconnect())
}
