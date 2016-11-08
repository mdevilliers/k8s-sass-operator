package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Printf("front-end starting...")

	go startHttpServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("front-end exiting...")
}

func startHttpServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "front-end alive!")
	})

	log.Printf("front-end running...")

	http.ListenAndServe(":3000", nil)
}
