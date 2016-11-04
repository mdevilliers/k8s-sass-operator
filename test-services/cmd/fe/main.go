package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Printf("front-end running...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("shutdown signal received, exiting...")
}
