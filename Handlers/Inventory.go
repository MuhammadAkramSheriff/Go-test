package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetInventory(c *fiber.Ctx) error {
	username := c.Locals("username")
	return c.JSON(fiber.Map{
		"message":  "Inventory data",
		"user":     username,
	})
}
