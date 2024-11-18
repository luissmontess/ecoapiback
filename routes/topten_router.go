package routes

import (
	"ecoApi/handlers"
	"github.com/gofiber/fiber/v2"

	// "ecoApi/middleware"
)

func SetTopTenRoutes(app *fiber.App) {
	TopTenRoutes := app.Group("/api/topten")
	TopTenRoutes.Get("/", handlers.GetTopTen)
}