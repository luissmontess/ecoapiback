package routes

import (
	"ecoApi/handlers"
	"github.com/gofiber/fiber/v2"

	// "ecoApi/middleware"
)

func SetUserRoutes(app *fiber.App) {
	UserRoutes := app.Group("/api/user")
	UserRoutes.Post("/", handlers.CreatUser)
	UserRoutes.Get("/", handlers.GetUsers)
	UserRoutes.Get("/:uid", handlers.GetUserByFBID)
}