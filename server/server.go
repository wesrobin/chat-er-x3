package server

import (
	"database/sql"
	"log"
	"net"
)

type ChatterBoxSrv struct {
	db *sql.DB

	shutdowns []func() error
}

func (srv *ChatterBoxSrv) Kill() {
	for _, fn := range srv.shutdowns {
		err := fn()
		if err != nil {
			log.Print(err)
		}
	}
	srv.db.Close()
}

func (srv *ChatterBoxSrv) Run() error {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	// Could also defer but I like dis
	srv.shutdowns = append(srv.shutdowns, listener.Close)

	log.Print("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handle(conn)
	}
}

// TODO: Client connects to the server, opens socket, connecting them to an event listener
// Client can send and receive messages via this socket
// Upon receipt of a message, the server persists it and emits an event
// Event is consumed in the socket loop, sending the message to the connected client
// Extension:
// Support https
// Microservices, managed by k8s

func handle(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Print(err)
			return
		}

		log.Print(string(buffer[:n]))

		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Print(err)
			return
		}
	}
}

func New() (*ChatterBoxSrv, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	return &ChatterBoxSrv{db: db}, nil
}
