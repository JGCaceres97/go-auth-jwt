package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jgcaceres97/go-auth-jwt/src/settings"
)

func CreateJWT(Id *string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    *Id,
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
	})

	return claims.SignedString([]byte(*settings.JWTSecret))
}

func GetJwtIssuer(c *fiber.Ctx) *string {
	token, _ := getJWT(c)

	return &token.Claims.(*jwt.StandardClaims).Issuer
}

func getJWT(c *fiber.Ctx) (*jwt.Token, error) {
	cookie := c.Cookies("jwt")

	return jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, keyFunc)
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(*settings.JWTSecret), nil
}

func CheckJWT(c *fiber.Ctx) error {
	_, err := getJWT(c)

	return err
}
