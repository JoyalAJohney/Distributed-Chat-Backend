package chat

import (
	"sync"

	"github.com/gofiber/contrib/websocket"

	"realtime-chat/src/cache"
	"realtime-chat/src/models"
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
