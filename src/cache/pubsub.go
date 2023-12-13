package cache

import (
	"log"
	"context"
	"encoding/json"

	"realtime-chat/src/models"
)

func PublishMessage(room string, message *chat.Message) {
	ctx := context.Background()
	msg, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	RedisClient.Publish(ctx, room, msg)
}