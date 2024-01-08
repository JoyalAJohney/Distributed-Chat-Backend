package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"realtime-chat/src/cache"
	"realtime-chat/src/chat"
	"realtime-chat/src/database"
)

func main() {
	app := fiber.New()
	cache.InitRedis()
	database.InitPostgres()

	app.Use("/ws", upgradeToWebSocket)
	app.Get("/ws/chat", websocket.New(chat.WebSocketHandler))

	log.Fatal(app.Listen(":8080"))
}

func upgradeToWebSocket(context *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(context) {
		context.Locals("allowed", true)
		return context.Next()
	}
	return fiber.ErrUpgradeRequired
}
