package server

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wesrobin/chat-er-x3/api_types"
	"log"
	"net/http"
)

func messageHander(srv *ChatterBoxSrv) (gin.HandlerFunc, error) {
	prod, err := sarama.NewSyncProducer([]string{kafkaAddr}, srv.kafkaCfg)
	if err != nil {
		return nil, err
	}
	srv.shutdowns = append(srv.shutdowns, prod.Close)

	return func(ctx *gin.Context) {
		var msg api_types.Message
		if err := ctx.BindJSON(&msg); err != nil {
			_ = ctx.Error(err)
			return
		}

		log.Print("received: ", msg)

		jsonMsg, _ := json.Marshal(msg)
		event := sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(jsonMsg)}
		_, _, err := prod.SendMessage(&event)
		if err != nil {
			log.Printf("err: %v", err)
			_ = ctx.Error(err)
			return
		}

		ctx.Status(http.StatusOK)
	}, nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(srv *ChatterBoxSrv) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return
		}
		userID := ctx.GetHeader(api_types.HeaderUserID)

		srv.register(userID, conn)
	}
}
