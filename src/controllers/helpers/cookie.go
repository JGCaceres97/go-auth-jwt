package helpers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateJwtCookie(Id *string) (*fiber.Cookie, error) {
	token, err := CreateJWT(Id)
	if err != nil {
		return nil, err
	}

	return &fiber.Cookie{
		Expires:  time.Now().Add(10 * time.Minute).UTC(),
		HTTPOnly: true,
		Name:     "jwt",
		Value:    token,
	}, nil
}

func DeleteJwtCookie() *fiber.Cookie {
	return &fiber.Cookie{
		Expires:  time.Now().Add(-time.Minute),
		HTTPOnly: true,
		Name:     "jwt",
		Value:    "",
	}
}
