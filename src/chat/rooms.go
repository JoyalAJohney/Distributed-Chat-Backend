package chat

import (
	"sync"

	"github.com/gofiber/contrib/websocket"

	"realtime-chat/src/models"
	"realtime-chat/src/cache"
)

var roomMutex = &sync.Mutex{}
// Shared resource
var roomsAndMembersMap = make(map[string]map[*websocket.Conn]bool)


func JoinRoom(room string, user *models.User) {
	roomMutex.Lock()
	if roomsAndMembersMap[room] == nil {
		roomsAndMembersMap[room] = make(map[*websocket.Conn]bool)
	}
	roomsAndMembersMap[room][user.Connection] = true
	roomMutex.Unlock()

	cache.SubscribeToRoom(room, func(msg models.Message) {
		BroadcastToRoom(room, msg)
	})
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