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

type Message struct {
	User   string    `json:"user"`
	Msg    string    `json:"message"`
	SentAt time.Time `json:"sent_at"`
}
