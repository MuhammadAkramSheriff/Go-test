package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-test/models"
)

// @Summary Get user by ID
// @Description Retrieve user details using ID
// @Tags Users
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} models.UserResponse
// @Router /api/user/{id} [get]
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := models.UserResponse{
		ID:   id,
		Name: "John Doe",
	}
	return c.JSON(user)
}