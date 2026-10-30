package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/RaucousDave/streakly/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterHabitRoutes(r fiber.Router) {
	r.Use(middleware.RequireAuth)
	r.Post("/", handlers.AddHabit)
	r.Get("/", handlers.GetAllHabits)
	r.Patch("/:id/archive", handlers.ArchiveHabits)
	r.Post("/:id/freeze", handlers.FreezeHabit)
}
