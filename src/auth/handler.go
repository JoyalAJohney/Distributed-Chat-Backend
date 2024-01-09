package auth

import (
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"realtime-chat/src/config"
	"realtime-chat/src/database"
	"realtime-chat/src/models"
)

func SignUp(ctx *fiber.Ctx) error {
	var request models.AuthRequest
	if err := ctx.BodyParser(&request); err != nil {
		log.Println("Error parsing request body for Signup:", err)
		return fiber.ErrBadRequest
	}

	if request.Username == "" || request.Password == "" {
		log.Println("Username or password is empty")
		return fiber.ErrBadRequest
	}

	var userExists database.DBUser
	if err := database.DB.Where("name = ?", request.Username).First(&userExists).Error; err == nil {
		log.Println("User already exists")
		return fiber.ErrConflict
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return fiber.ErrInternalServerError
	}

	user := database.DBUser{
		Name:     request.Username,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Error adding user to database: %v", err)
		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func Login(ctx *fiber.Ctx) error {
	var request models.AuthRequest
	if err := ctx.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body for Login: %v", err)
		return fiber.ErrBadRequest
	}

	if request.Username == "" || request.Password == "" {
		log.Println("Username or password is empty")
		return fiber.ErrBadRequest
	}

	var user database.DBUser
	if err := database.DB.Where("name = ?", request.Username).First(&user).Error; err != nil {
		log.Printf("Error finding user in database: %v", err)
		return fiber.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		log.Printf("Password mismatch: %v", err)
		return fiber.ErrUnauthorized
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	token, err := GenerateJWT(userID, user.Name)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(fiber.Map{
		"status": "success",
		"token":  token,
	})
}

func GenerateJWT(userID string, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString([]byte(config.Config.JwtSecret))
}
