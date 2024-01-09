package auth

import (
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"realtime-chat/src/config"
)

func AuthorizationMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		log.Println("No Authorization header provided")
		return fiber.ErrUnauthorized
	}

	if err := ValidateJWTToken(authHeader); err != nil {
		log.Println("Error validating JWT token:", err)
		return fiber.ErrUnauthorized
	}

	return ctx.Next()
}

func ValidateJWTToken(authHeader string) error {
	tokenString := strings.Split(authHeader, " ")[1]
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JwtSecret), nil
	})

	return err
}

func ParseJWTToken(authHeader string) (userID string, userName string) {
	tokenString := strings.Split(authHeader, " ")[1]
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = claims["userID"].(string)
		userName = claims["username"].(string)
	}

	return userID, userName
}
