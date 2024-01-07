package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wesrobin/chat-er-x3/server"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	srv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Kill()

	srv.Run()

	<-stop
}
