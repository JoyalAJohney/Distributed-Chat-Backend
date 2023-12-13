package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"

	"realtime-chat/src/chat"
	"realtime-chat/src/cache"
)

func main() {
	app := fiber.New()
	cache.InitRedis()

	app.Use("/ws", upgradeToWebSocket)
	app.Get("/ws/:id", websocket.New(chat.WebSocketHandler))

	log.Fatal(app.Listen(":3000"))
}

func upgradeToWebSocket(context *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(context) {
		context.Locals("allowed", true)
		return context.Next()
	}
	return fiber.ErrUpgradeRequired
}