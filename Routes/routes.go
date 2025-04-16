package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-test/handlers"
)

func Register(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/user/:id", handlers.GetUser)
}