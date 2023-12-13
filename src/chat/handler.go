package chat

import (
	"log"

	"github.com/gofiber/contrib/websocket"

	"realtime-chat/src/models"
	"realtime-chat/src/utils"
)

func WebSocketHandler(conn *websocket.Conn) {
	userID := utils.GenerateUserID()
	user := &models.User{
		ID:         userID,
		Connection: conn,
	}
	log.Printf("User %s connected\n", userID)

	for {
		var message models.Message
		if err := conn.ReadJSON(&message); err != nil {
			break
		}

		log.Printf("Received message: %v\n", message)
		message.Sender = userID
		switch MessageType(message.Type) {
		case JoinRoomType:
			JoinRoom(message.Room, user)
		case LeaveRoomType:
			LeaveRoom(message.Room, user)
		case ChatMessageType:
			SendMessageToRoom(message, user)
		default:
			log.Println("Unknown message type")
		}
	}

	// User disconnected
	LeaveAllRooms(user)
}
