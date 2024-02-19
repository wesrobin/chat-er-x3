package server

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wesrobin/chat-er-x3/api_types"
	"log"
	"net/http"
	"sync"
)

type ChatterBoxSrv struct {
	db *sql.DB

	httpSrv *http.Server

	connMu  sync.RWMutex
	clients map[string]*websocket.Conn

	shutdowns []func() error
}

func (srv *ChatterBoxSrv) Start(ctx context.Context) error {
	r := gin.Default()
	registerHandlers(srv, r)
	err := r.Run("localhost:8080")
	if err != nil {
		return err
	}
	log.Print("HTTP server listening on port 8080")

	return nil
}

func registerHandlers(srv *ChatterBoxSrv, r *gin.Engine) {
	r.POST("/message", messageHander(srv))
	r.GET("/ws", websocketHandler(srv))
}

func (srv *ChatterBoxSrv) register(userID string, conn *websocket.Conn) {
	srv.connMu.Lock()
	srv.clients[userID] = conn
	srv.connMu.Unlock()

	log.Print("connected new user")
}

func (srv *ChatterBoxSrv) deregister(userID string) {
	srv.connMu.Lock()
	if c, ok := srv.clients[userID]; ok {
		c.Close()
		delete(srv.clients, userID)
		log.Print("removed ", userID)
	}
	srv.connMu.Unlock()
}

func (srv *ChatterBoxSrv) FanOut(msg api_types.Message) {
	srv.connMu.RLock()
	defer srv.connMu.RUnlock()

	log.Print("fanout ", msg)

	for uid, c := range srv.clients {
		err := c.WriteJSON(msg)
		if err != nil {
			log.Print(err)
			srv.deregister(uid)
			continue
		}
	}
}

func (srv *ChatterBoxSrv) Kill() {
	log.Printf("Kill called, %v shutdowns to run", len(srv.shutdowns))
	for _, fn := range srv.shutdowns {
		err := fn()
		if err != nil {
			log.Print(err)
		}
		log.Printf("kill run")
	}
}

func New() (*ChatterBoxSrv, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	return &ChatterBoxSrv{
		db:        db,
		shutdowns: []func() error{db.Close},
		clients:   map[string]*websocket.Conn{},
	}, nil
}
