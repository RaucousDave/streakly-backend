package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App, h *handlers.Handler) {
	v1 := app.Group("api/v1")

	RegisterAppConfigRoutes(v1.Group("/app"))
	RegisterAuthRoutes(v1.Group("/auth"), h.Auth)
	RegisterHabitRoutes(v1.Group("/habits"), h.Habits)
	RegisterCompletionRoute(v1.Group("/completions"), h.Completion)
	RegisterDashboardRoute(v1.Group("/dashboard"), h.Dashboard)
}
