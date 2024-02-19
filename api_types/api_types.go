package api_types

import (
	"time"
)

const (
	HeaderUserID = "x-wes-user-id"
)

type ConnRequest struct {
	User   string `json:"user"`
	RoomID int64  `json:"room_id"`
}

type MessageRequest struct {
	User string  `json:"user"`
	Msg  Message `json:"msg"`
}

type MessageResponse struct {
	User string  `json:"user"`
	Msg  Message `json:"msg"`
}

type Message struct {
	Msg    string    `json:"message"`
	SentAt time.Time `json:"sent_at"`
}
