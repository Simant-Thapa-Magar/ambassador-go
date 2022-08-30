package middlewares

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

const SECRET = "secret"
const COOKIE = "ambassadorJWT"

type ClaimsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func GenerateToken(id uint, scope string) (string, error) {
	payload := ClaimsWithScope{}
	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()
	payload.Scope = scope
	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(SECRET))
}

func IsAuthenticatedUser(c *fiber.Ctx) error {

	cookie := c.Cookies(COOKIE)

	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	payload := token.Claims.(*ClaimsWithScope)

	isAdminRoute := strings.Contains(c.Path(), "api/admin")

	if isAdminRoute && payload.Scope != "admin" {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.Next()
}

func GetAuthenticatedUserId(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies(COOKIE)

	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	payload := token.Claims.(*ClaimsWithScope)

	uId, e := strconv.Atoi(payload.Subject)

	if e != nil {
		return 0, e
	}

	return uint(uId), nil
}
