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
)

var request struct {
	Username string
	Password string
}

func SignUp(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(&request); err != nil {
		log.Println("Error parsing request body for Signup:", err)
		return fiber.ErrBadRequest
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return fiber.ErrInternalServerError
	}

	user := database.DBUser{
		Name:     request.Username,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user); err != nil {
		log.Println("Error adding user to database:", err)
		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func Login(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(&request); err != nil {
		log.Println("Error parsing request body for Login:", err)
		return fiber.ErrBadRequest
	}

	var user database.DBUser
	if err := database.DB.Where("name = ?", request.Username).First(&user).Error; err != nil {
		log.Println("Error finding user in database:", err)
		return fiber.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		log.Println("Password mismatch:", err)
		return fiber.ErrUnauthorized
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	token, err := GenerateJWT(userID, user.Name)
	if err != nil {
		log.Println("Error generating JWT:", err)
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
