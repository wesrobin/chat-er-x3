package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wesrobin/chat-er-x3/api_types"
)

const (
	kafkaAddr = "192.168.1.115:9092"
	topic     = "chat-log"
)

type ChatterBoxSrv struct {
	db *sql.DB

	httpSrv *http.Server

	kafkaCfg *sarama.Config

	connMu  sync.RWMutex
	clients map[string]*websocket.Conn

	shutdowns []func() error
}

func (srv *ChatterBoxSrv) Start(ctx context.Context) error {
	r := gin.Default()
	registerHandlers(srv, r)
	go func() {
		err := r.Run("0.0.0.0:8071")
		if err != nil {
			log.Fatalf("fatal err: %v", err)
		}
	}()

	log.Print("HTTP server listening on port 8080")

	go srv.ListenKafkaForever()

	return nil
}

func registerHandlers(srv *ChatterBoxSrv, r *gin.Engine) {
	h, _ := messageHander(srv)
	r.POST("/message", h)
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

func (srv *ChatterBoxSrv) ListenKafkaForever() {
	csmr, err := sarama.NewConsumer([]string{kafkaAddr}, srv.kafkaCfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer csmr.Close()

	partitionConsumer, err := csmr.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	log.Printf("Listening to kafka topic %v", topic)

	for {
		select {
		case e := <-partitionConsumer.Messages():
			var msg api_types.Message
			err := json.Unmarshal(e.Value, &msg)
			if err != nil {
				log.Printf("err: %v", err.Error())
				continue
			}
			srv.fanOut(msg)
		case err := <-partitionConsumer.Errors():
			log.Printf("err: %v\n", err.Error())
		}
	}
}

func (srv *ChatterBoxSrv) fanOut(msg api_types.Message) {
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
	for _, fn := range srv.shutdowns {
		err := fn()
		if err != nil {
			log.Print(err)
		}
	}
}

func New() (*ChatterBoxSrv, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	return &ChatterBoxSrv{
		db:        db,
		kafkaCfg:  cfg,
		shutdowns: []func() error{db.Close},
		clients:   map[string]*websocket.Conn{},
	}, nil
}
