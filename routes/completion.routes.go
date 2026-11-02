package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/RaucousDave/streakly/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterCompletionRoute(r fiber.Router, h *handlers.CompletionHandler) {
	r.Use(middleware.RequireAuth)
	r.Post("/", h.SyncHabits)
}
