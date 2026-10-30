package routes

import "github.com/gofiber/fiber/v3"

func Register(app *fiber.App) {
	v1 := app.Group("api/v1")

	RegisterAppConfigRoutes(v1.Group("/app"))
	RegisterAuthRoutes(v1.Group("/auth"))
	RegisterHabitRoutes(v1.Group("/habits"))
	RegisterCompletionRoute(v1.Group("/completions"))
	RegisterDashboardRoute(v1.Group("/dashboard"))
}
