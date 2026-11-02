package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(r fiber.Router, h *handlers.AuthHandler) {

	r.Post("/login", h.Login)
	r.Post("/sign-up", h.Register)
	r.Post("/google", h.GoogleLogin)
}
