package routes

import (
	"ecoApi/handlers"
	"github.com/gofiber/fiber/v2"

	// "ecoApi/middleware"
)

func SetCollaboratorRoutes(app *fiber.App) {
	CollaboratorRoutes := app.Group("/api/collaborator")
	CollaboratorRoutes.Post("/", handlers.CreateCollaborator)
	CollaboratorRoutes.Get("/", handlers.GetCollaborators)
	CollaboratorRoutes.Get("/:uid", handlers.GetCollaboratorByFBID)
}    