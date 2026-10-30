package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/RaucousDave/streakly/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterCompletionRoute(r fiber.Router) {
	r.Use(middleware.RequireAuth)
	r.Post("/", handlers.SyncHabits)
}
