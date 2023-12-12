package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"

	"realtime-chat/src/cache"
)

func main() {
	app := fiber.New()
	cache.InitRedis()

	app.Use("/ws", upgradeToWebSocket)
	app.Get("/ws/:id", websocket.New(handleWebSocket))

	log.Fatal(app.Listen(":3000"))
}

func upgradeToWebSocket(context *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(context) {
		context.Locals("allowed", true)
		return context.Next()
	}
	return fiber.ErrUpgradeRequired
}

func handleWebSocket(context *websocket.Conn) {
	log.Println(context.Locals("allowed"))
	log.Println(context.Params("id"))
	log.Println(context.Query("v"))
	log.Println(context.Cookies("session"))

	var (
		messageType int
		message []byte
		err error
	)

	for {
		messageType, message, err = context.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("received: %s", message)
		log.Printf("message type: %d", messageType)
		err = context.WriteMessage(messageType, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}