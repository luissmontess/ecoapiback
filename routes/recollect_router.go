package routes

import (
	"ecoApi/handlers"
	"github.com/gofiber/fiber/v2"

	// "ecoApi/middleware"
)

func SetRecollectRoutes(app *fiber.App) {
	RecollectRoutes := app.Group("/api/recollect")
	RecollectRoutes.Post("/", handlers.CreateRecollect)
	RecollectRoutes.Get("/", handlers.GetActiveRecollects)
	RecollectRoutes.Get("/:recollect_id", handlers.GetRecollectByID)
	RecollectRoutes.Get("/collaborator_recollects/:uid", handlers.GetCollaboratorRecollects)
	RecollectRoutes.Patch("/add_to_recollect/:recollect_id", handlers.AddtoRecollect)
}