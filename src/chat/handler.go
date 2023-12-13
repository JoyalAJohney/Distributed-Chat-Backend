package chat

import (
	"log"

	"github.com/matoous/go-nanoid"
	"github.com/gofiber/contrib/websocket"

	"realtime-chat/src/cache"
	"realtime-chat/src/models"
)

func WebSocketHandler(conn *websocket.Conn) {
	userID := generateUserID()
	user := &User{
		ID: userID,
		Connection: conn,
	}

	for {
		var message Message
		if err := conn.ReadJSON(&message); err != nil {
			break
		}
		message.Sender = userID
		cache.PublishMessage(message.Room, &message)
	}
}

func generateUserID() string {
	id, err := gonanoid.New()
	return id
}