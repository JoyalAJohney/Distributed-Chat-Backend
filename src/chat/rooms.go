package chat

import (
	"context"
	"log"
	"sync"
	"time"

	"realtime-chat/src/cache"
	"realtime-chat/src/database"
	"realtime-chat/src/models"
	"realtime-chat/src/utils"
)

var roomMutex = &sync.Mutex{}
var subscribedRooms = make(map[string]bool)

func JoinRoom(room string, user *models.User) {
	key := "room:" + room
	if err := addUserToRoomInRedis(key, user); err != nil {
		log.Println("Failed to add user to room:", err)
		utils.SendErrorMessage(user.Connection, "Unable to join room")
		return
	}

	cache.SubscribeToRoom(room, func(room string, message *models.Message) {
		BroadcastToRoom(room, *message)
	})

	log.Printf("User %s joined room %s\n", user.ID, room)
}

func BroadcastToRoom(room string, message models.Message) {
	key := "room:" + room
	for _, userID := range getAllMembersInRoom(key) {
		// Get the websocket connection for the user from the local map
		if conn, exists := GetConnection(userID); exists {
			// Send the message to the user
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Error sending message to user %s: %v\n", userID, err)
				conn.Close()
				RemoveConnection(userID)
			}
		}
	}
}

func SendMessageToRoom(message models.Message, user *models.User) {
	key := "room:" + message.Room
	if !isUserInRoom(key, user) {
		utils.SendErrorMessage(user.Connection, "You are not a member of this room")
		return
	}
	cache.PublishMessage(message.Room, &message)
	saveMessageToDatabase(message)
}

func LeaveRoom(room string, user *models.User) {
	key := "room:" + room
	if err := removeUserFromRoomInRedis(key, user); err != nil {
		log.Println("Failed to remove user from room:", err)
		utils.SendErrorMessage(user.Connection, "Unable to leave room")
		return
	}

	cache.CheckAndUnsubscribeFromRoom(room)
}

func LeaveAllRooms(user *models.User) {
	for _, room := range cache.GetAllRooms() {
		isMember := isUserInRoom(room, user)
		if isMember {
			removeUserFromRoomInRedis(room, user)
		}
	}
}

// Helper methods
func saveMessageToDatabase(message models.Message) {
	var currentTime = time.Now()
	dbMessage := database.DBMessage{
		UserID:    message.Sender,
		RoomID:    message.Room,
		Message:   message.Content,
		Timestamp: &currentTime,
	}

	if err := database.DB.Create(&dbMessage).Error; err != nil {
		log.Println("Error saving message to database:", err)
		return
	}
}

func addUserToRoomInRedis(room string, user *models.User) error {
	ctx := context.Background()
	_, err := cache.RedisClient.SAdd(ctx, room, user.ID).Result()
	return err
}

func removeUserFromRoomInRedis(room string, user *models.User) error {
	ctx := context.Background()
	_, err := cache.RedisClient.SRem(ctx, room, user.ID).Result()
	return err
}

func isUserInRoom(room string, user *models.User) bool {
	ctx := context.Background()
	isMember, err := cache.RedisClient.SIsMember(ctx, room, user.ID).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return isMember
}

func getAllMembersInRoom(room string) []string {
	ctx := context.Background()
	members, err := cache.RedisClient.SMembers(ctx, room).Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	return members
}
