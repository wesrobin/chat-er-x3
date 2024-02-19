package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wesrobin/chat-er-x3/server"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL) // Not working?

	srv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Kill()

	err = srv.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	<-stop
}
