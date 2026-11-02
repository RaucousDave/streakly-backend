package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/RaucousDave/streakly/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterDashboardRoute(r fiber.Router, h *handlers.DashbaordHandler) {
	r.Use(middleware.RequireAuth)
	r.Get("/", h.GetDashboard)
}
