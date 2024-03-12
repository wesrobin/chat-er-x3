package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/wesrobin/chat-er-x3/api_types"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	ctx, cFn := context.WithCancel(context.Background())

	r := rand.Int31()
	userID := "user" + strconv.Itoa(int(r))

	c := webSocket(userID)
	defer cFn()
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Print(err.Error())
		}
	}(c)

	go listen(ctx, c)

	for {
		scanner := bufio.NewScanner(os.Stdin)

		var inp string
		for scanner.Scan() {
			tkn := scanner.Text()
			if tkn == "" {
				break
			}
			inp += tkn
		}

		if inp == "quit" {
			return
		}

		req := api_types.Message{
			User:   userID,
			Msg:    inp,
			SentAt: time.Now(),
		}
		requestBody, err := json.Marshal(req)
		if err != nil {
			log.Fatalf("Error encoding JSON: %v", err)
		}

		rsp, err := http.Post("http://127.0.0.1:8071/message", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Print(err.Error())
		}
		err = rsp.Body.Close()
		if err != nil {
			log.Print(err.Error())
		}
		log.Print("Response code: ", rsp.StatusCode)
	}
}

func listen(ctx context.Context, c *websocket.Conn) {
	for {
		if ctx.Err() != nil {
			return
		}
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		log.Print("Rsp", string(msg))
	}
}

func webSocket(user string) *websocket.Conn {
	head := make(http.Header)
	head.Set(api_types.HeaderUserID, user)
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8071/ws", head)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
