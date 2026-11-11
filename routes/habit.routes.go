package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/RaucousDave/streakly/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterHabitRoutes(r fiber.Router, h *handlers.HabitHandler) {
	r.Use(middleware.RequireAuth)
	r.Post("/", h.AddHabit)
	r.Get("/", h.GetAllHabits)
	r.Get("/archived", h.GetArchivedHabits)
	r.Patch("/:id/archive", h.ArchiveHabits)
	r.Patch("/:id", h.EditHabit)
	r.Post("/:id/freeze", h.FreezeHabit)
}
