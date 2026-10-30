package routes

import (
	"github.com/RaucousDave/streakly/handlers"
	"github.com/gofiber/fiber/v3"
)

func RegisterAppConfigRoutes(r fiber.Router){
	r.Get("/config", handlers.GetAppConfig)
}
