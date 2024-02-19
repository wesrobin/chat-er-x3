package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wesrobin/chat-er-x3/api_types"
	"log"
	"net/http"
)

func messageHander(srv *ChatterBoxSrv) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var msg api_types.MessageRequest
		if err := ctx.BindJSON(&msg); err != nil {
			_ = ctx.Error(err)
			return
		}

		log.Print("received: ", msg)
		srv.FanOut(msg.Msg)

		ctx.Status(http.StatusOK)
	}
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
