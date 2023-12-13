package chat

import (
	"log"

	"github.com/matoous/go-nanoid/v2"
	"github.com/gofiber/contrib/websocket"

	"realtime-chat/src/models"
)

func WebSocketHandler(conn *websocket.Conn) {
	userID := generateUserID()
	user := &models.User{
		ID: userID,
		Connection: conn,
	}

	for {
		var message models.Message
		if err := conn.ReadJSON(&message); err != nil {
			break
		}

		message.Sender = userID
		switch MessageType(message.Type) {
			case JoinRoom:
				JoinRoom(message.Room, user)
			case LeaveRoom:
				LeaveRoom(message.Room, user)
			case ChatMessage:
				BroadcastToRoom(message.Room, message)
			default:
				log.Println("Unknown message type")
		}
	}

	// User disconnected
	LeaveAllRooms(user)
}

func generateUserID() string {
	id, _ := gonanoid.New()
	return id
}