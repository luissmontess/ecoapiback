package routes

import (
	"ecoApi/handlers"
	"github.com/gofiber/fiber/v2"

	// "ecoApi/middleware"
)

func SetCourseRoutes(app *fiber.App) {
	CourseRoutes := app.Group("/api/course")
	CourseRoutes.Post("/", handlers.CreateCourse)
	CourseRoutes.Get("/", handlers.GetActiveCourses)
	CourseRoutes.Get("/:course_id", handlers.GetCourseByID)
	CourseRoutes.Get("/collaborator_courses/:uid", handlers.GetCollaboratorCourses)
	CourseRoutes.Patch("/add_assistant/:course_id", handlers.AddUserToCourse)
}