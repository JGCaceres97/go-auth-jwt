package controllers

import (
	"errors"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jgcaceres97/go-auth-jwt/src/database"
	"github.com/jgcaceres97/go-auth-jwt/src/models"
	"github.com/jgcaceres97/go-auth-jwt/src/settings"
	"gorm.io/gorm"
)

func checkJWT(c *fiber.Ctx) (*jwt.Token, error) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, keyFunc)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(*settings.JWTSecret), nil
}

func GetUsers(c *fiber.Ctx) error {
	_, err := checkJWT(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var users []*models.User

	database.DB.Select("id", "name", "email", "createdAt").Find(&users)
	return c.Status(fiber.StatusOK).JSON(&users)
}

func GetUser(c *fiber.Ctx) error {
	token, err := checkJWT(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var id *string
	paramId := c.Params("id")

	if paramId != "" {
		id = &paramId
	} else {
		id = &token.Claims.(*jwt.StandardClaims).Issuer
	}

	user := &models.User{Id: *id}

	query := database.DB.Select("id", "name", "email", "createdAt").First(&user)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		responseError := Error{
			Status:  fiber.StatusNotFound,
			Message: &ErrUserNotFound,
		}

		return c.Status(fiber.StatusNotFound).JSON(&responseError)
	}

	return c.Status(fiber.StatusOK).JSON(&user)
}

func UpdateUser(c *fiber.Ctx) error {
	token, err := checkJWT(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		responseError := Error{
			Status:  fiber.StatusBadRequest,
			Message: &ErrParsingData,
		}

		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&responseError)
	}

	var id *string
	paramId := c.Params("id")

	if paramId != "" {
		id = &paramId
	} else {
		id = &token.Claims.(*jwt.StandardClaims).Issuer
	}

	user := &models.User{Id: *id}
	rows := database.DB.Model(&user).Updates(models.User{Name: data["name"], Email: data["email"]}).RowsAffected
	if rows == 0 {
		responseError := Error{
			Status:  fiber.StatusNotFound,
			Message: &ErrUserNotFound,
		}

		return c.Status(fiber.StatusNotFound).JSON(&responseError)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func DeleteUser(c *fiber.Ctx) error {
	token, err := checkJWT(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var id *string
	paramId := c.Params("id")

	if paramId != "" {
		id = &paramId
	} else {
		id = &token.Claims.(*jwt.StandardClaims).Issuer
	}

	user := &models.User{Id: *id}
	rows := database.DB.Delete(&user).RowsAffected
	if rows == 0 {
		responseError := Error{
			Status:  fiber.StatusNotFound,
			Message: &ErrUserNotFound,
		}

		return c.Status(fiber.StatusNotFound).JSON(&responseError)
	}

	return c.SendStatus(fiber.StatusOK)
}
