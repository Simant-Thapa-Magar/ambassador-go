package middlewares

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func IsAuthenticatedUser(c *fiber.Ctx) error {

	cookie := c.Cookies("ambassadorJWT")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.Next()
}

func GetAuthenticatedUserId(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("ambassadorJWT")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	payload := token.Claims.(*jwt.StandardClaims)

	uId, e := strconv.Atoi(payload.Subject)

	if e != nil {
		return 0, e
	}

	return uint(uId), nil
}
