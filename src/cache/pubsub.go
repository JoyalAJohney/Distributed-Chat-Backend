package cache

import (
	"log"
	"context"
	"encoding/json"

	"realtime-chat/src/models"
)

func SubscribeToRoom(room string, handleMessage	func(msg models.Message)) {
	pubsub := RedisClient.Subscribe(context.Background(), room)
	go func() {
		channel := pubsub.Channel()
		for message := range channel {
			var chatMessage models.Message
			json.Unmarshal([]byte(message.Payload), &chatMessage)
			handleMessage(chatMessage)
		}
	}()
}

func PublishMessage(room string, message *models.Message) {
	ctx := context.Background()
	msg, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	RedisClient.Publish(ctx, room, msg)
}