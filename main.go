package main

import (
	"time"

	"github.com/RaucousDave/streakly/database"
	"github.com/RaucousDave/streakly/routes"
	"github.com/gofiber/fiber/v3"
)

var V1 any

func main() {
	app := fiber.New()

	database.DbConnect()
	routes.Register(app)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	app.Listen(":8000")
}
