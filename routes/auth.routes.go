package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(r fiber.Router) {

	r.Post("/login", handlers.Login)
	r.Post("/sign-up", handlers.Register)
	r.Post("/google", handlers.GoogleLogin)
}
