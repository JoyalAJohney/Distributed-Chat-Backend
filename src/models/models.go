package models

import (
	"github.com/gofiber/contrib/websocket"
)

type User struct {
	ID         string
	Connection *websocket.Conn
}

type MessageHandlerCallbackType func(room string, message *Message)

type Message struct {
	Sender  string `json:"sender"`
	Room    string `json:"room"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ErrorMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
