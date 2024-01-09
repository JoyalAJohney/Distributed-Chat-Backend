package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"realtime-chat/src/auth"
	"realtime-chat/src/cache"
	"realtime-chat/src/chat"
	"realtime-chat/src/database"
)

func main() {
	app := fiber.New()
	cache.InitRedis()
	database.InitPostgres()

	app.Post("/auth/signup", auth.SignUp)
	app.Post("/auth/login", auth.Login)

	// Secure websocket connection
	app.Use("/ws", upgradeToWebSocket)
	app.Get("/ws/chat", websocket.New(chat.WebSocketHandler))

	log.Fatal(app.Listen(":8080"))
}

// Authorize and Upgrate to websocket
func upgradeToWebSocket(context *fiber.Ctx) error {
	authHeader := context.Get("Authorization")
	if authHeader == "" {
		log.Println("No Authorization header provided")
		return fiber.ErrUnauthorized
	}

	// Validate JWT token
	if err := auth.ValidateJWTToken(authHeader); err != nil {
		log.Println("Error validating JWT token:", err)
		return fiber.ErrUnauthorized
	}

	userID, userName := auth.ParseJWTToken(authHeader)
	if websocket.IsWebSocketUpgrade(context) {
		context.Locals("allowed", true)
		context.Locals("userID", userID)
		context.Locals("userName", userName)
		return context.Next()
	}
	return fiber.ErrUpgradeRequired
}
