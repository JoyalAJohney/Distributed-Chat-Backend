package models

import (
	"github.com/gofiber/contrib/websocket"
)

type User struct {
	ID string
	Connection *websocket.Conn
}

type Message struct {
	Sender string `json:"sender"`
	Room string `json:"room"`
	Content string `json:"content"`
}