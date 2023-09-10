package controllers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jgcaceres97/go-auth-jwt/src/controllers/helpers"
	"github.com/jgcaceres97/go-auth-jwt/src/database"
	"github.com/jgcaceres97/go-auth-jwt/src/models"
	"gorm.io/gorm"
)

func GetUsers(c *fiber.Ctx) error {
	var users []*models.User

	database.DB.Select("id", "name", "email", "createdAt").Find(&users)
	return c.Status(fiber.StatusOK).JSON(&users)
}

func GetUser(c *fiber.Ctx) error {
	var id *string
	paramId := c.Params("id")

	if paramId != "" {
		id = &paramId
	} else {
		id = helpers.GetJwtIssuer(c)
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
		id = helpers.GetJwtIssuer(c)
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
	var id *string
	paramId := c.Params("id")

	if paramId != "" {
		id = &paramId
	} else {
		id = helpers.GetJwtIssuer(c)
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
