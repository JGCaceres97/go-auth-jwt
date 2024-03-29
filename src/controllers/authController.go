package controllers

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jgcaceres97/go-auth-jwt/src/controllers/helpers"
	"github.com/jgcaceres97/go-auth-jwt/src/database"
	"github.com/jgcaceres97/go-auth-jwt/src/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrParsingData           = "Error parsing the data."
	ErrUserAlreadyExists     = "User already exists."
	ErrUserNotFound          = "User not found."
	ErrUserIncorrectPassword = "Incorrect password."
	ErrEncryptingPassword    = "Error encrypting the given password."
	ErrGeneratingJWT         = "Error generating JWT."
)

type Error struct {
	Status  uint16  `json:"status"`
	Message *string `json:"message"`
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		responseError := Error{
			Status:  fiber.StatusBadRequest,
			Message: &ErrParsingData,
		}

		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&responseError)
	}

	var user *models.User

	query := database.DB.Where("email = ?", data["email"]).First(&user)
	if !errors.Is(query.Error, gorm.ErrRecordNotFound) {
		responseError := Error{
			Status:  fiber.StatusBadRequest,
			Message: &ErrUserAlreadyExists,
		}

		return c.Status(fiber.StatusBadRequest).JSON(&responseError)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		responseError := Error{
			Status:  fiber.StatusInternalServerError,
			Message: &ErrEncryptingPassword,
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(&responseError)
	}

	user = &models.User{
		Id:        uuid.NewString(),
		Name:      data["name"],
		Email:     data["email"],
		Password:  &password,
		CreatedAt: time.Now(),
	}

	database.DB.Create(&user)

	return c.Status(fiber.StatusCreated).JSON(&user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		responseError := Error{
			Status:  fiber.StatusInternalServerError,
			Message: &ErrParsingData,
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(&responseError)
	}

	var user *models.User

	query := database.DB.Where("email = ?", data["email"]).First(&user)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		responseError := Error{
			Status:  fiber.StatusNotFound,
			Message: &ErrUserNotFound,
		}

		return c.Status(fiber.StatusNotFound).JSON(&responseError)
	}

	if err := bcrypt.CompareHashAndPassword(*user.Password, []byte(data["password"])); err != nil {
		responseError := Error{
			Status:  fiber.StatusBadRequest,
			Message: &ErrUserIncorrectPassword,
		}

		return c.Status(fiber.StatusBadRequest).JSON(&responseError)
	}

	cookie, err := helpers.CreateJwtCookie(&user.Id)
	if err != nil {
		responseError := Error{
			Status:  fiber.StatusInternalServerError,
			Message: &ErrGeneratingJWT,
		}

		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&responseError)
	}

	c.Cookie(cookie)
	return c.SendStatus(fiber.StatusOK)
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(helpers.DeleteJwtCookie())

	return c.SendStatus(fiber.StatusOK)
}

func CheckAuth(c *fiber.Ctx) error {
	err := helpers.CheckJWT(c)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
