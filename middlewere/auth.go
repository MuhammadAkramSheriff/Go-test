package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key") // Must match what you used in user service

func AuthMiddleware(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Optionally pass user info to next handlers
		c.Locals("username", claims["username"])
		return c.Next()
	} else {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}
}
