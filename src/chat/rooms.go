package chat

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"

	"realtime-chat/src/cache"
	"realtime-chat/src/models"
	"realtime-chat/src/utils"
)

var (
	roomMutex          = &sync.Mutex{}
	roomsAndMembersMap = make(map[string]map[*websocket.Conn]bool)

	subscriptionsMutex  = &sync.Mutex{}
	activeSubscriptions = make(map[string]bool)
)

func JoinRoom(room string, user *models.User) {
	roomMutex.Lock()
	if roomsAndMembersMap[room] == nil {
		roomsAndMembersMap[room] = make(map[*websocket.Conn]bool)
	}
	roomsAndMembersMap[room][user.Connection] = true
	roomMutex.Unlock()

	subscriptionsMutex.Lock()
	if !activeSubscriptions[room] {
		activeSubscriptions[room] = true
		cache.SubscribeToRoom(room, func(msg models.Message) {
			BroadcastToRoom(room, msg)
		}, RoomCleanup)
	}
	subscriptionsMutex.Unlock()

	log.Printf("User %s joined room %s\n", user.ID, room)
	log.Printf("Room %s has %d members\n", room, len(roomsAndMembersMap[room]))
}

func BroadcastToRoom(room string, message models.Message) {
	roomMutex.Lock()
	defer roomMutex.Unlock()
	for conn := range roomsAndMembersMap[room] {
		if err := conn.WriteJSON(message); err != nil {
			conn.Close()
			delete(roomsAndMembersMap[room], conn)
		}
	}
}

func SendMessageToRoom(message models.Message, user *models.User) {
	if !isUserInRoom(message.Room, user) {
		utils.SendErrorMessage(user.Connection, "You are not a member of this room")
		return
	}
	cache.PublishMessage(message.Room, &message)
}

func isUserInRoom(room string, user *models.User) bool {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	if roomMembers, exists := roomsAndMembersMap[room]; exists {
		_, userExists := roomMembers[user.Connection]
		return userExists
	}
	return false
}

func RoomCleanup(room string) {
	subscriptionsMutex.Lock()
	delete(activeSubscriptions, room)
	subscriptionsMutex.Unlock()

	roomMutex.Lock()
	if len(roomsAndMembersMap[room]) == 0 {
		delete(roomsAndMembersMap, room)
	}
	roomMutex.Unlock()
}

func LeaveRoom(room string, user *models.User) {
	roomMutex.Lock()
	delete(roomsAndMembersMap[room], user.Connection)
	roomMutex.Unlock()
}

func LeaveAllRooms(user *models.User) {
	roomMutex.Lock()
	for room := range roomsAndMembersMap {
		delete(roomsAndMembersMap[room], user.Connection)
	}
	roomMutex.Unlock()
}
