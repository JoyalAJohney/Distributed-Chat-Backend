package cache

import (
	"context"
	"encoding/json"
	"log"

	"realtime-chat/src/models"
)

func SubscribeToRoom(room string, handleMessage func(msg models.Message), cleanupCallback func(room string)) {
	// Establish a TCP connection that listens for messages
	pubsub := RedisClient.Subscribe(context.Background(), room)

	go func() {
		channel := pubsub.Channel()
		for message := range channel {
			log.Printf("Received message from channel: %s\n", message.Payload)
			var chatMessage models.Message
			err := json.Unmarshal([]byte(message.Payload), &chatMessage)
			if err != nil {
				log.Printf("Error decoding message from channel: %v\n", err)
				continue
			}
			handleMessage(chatMessage)
		}

		// cleanup after leaving the room
		cleanupCallback(room)
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
	log.Printf("Published message to channel: %s\n", room)
}
