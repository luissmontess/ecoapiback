package main

import (
	"ecoApi/config"
	"ecoApi/routes"
	"ecoApi/handlers"

	"ecoApi/firebase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron/v3"
	"github.com/joho/godotenv"
	"fmt"
	"os"

	"log"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // permite cualquier origen, solo pruebas
		AllowMethods: "GET,POST,PUT,DELETE", // peticiones permitidas
	}))

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	c := cron.New()

	config.ConnectMongoDB()

	firebase.InitFirebase()

	routes.SetUserRoutes(app)
	routes.SetCollaboratorRoutes(app)
	routes.SetCourseRoutes(app)
	routes.SetRecollectRoutes(app)
	routes.SetTopTenRoutes(app)

	_, err := c.AddFunc("40 10 * * 6", func() {
		err := handlers.UpdateTopTen()
		if err != nil {
			log.Printf("Failed to update top ten users: %v", err)
		} else {
			fmt.Println("Successfully updated top ten users at 12 PM on Sunday")
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule UpdateTopTen: %v", err)
	}

	// empezar cron
	c.Start()

	log.Fatal(app.Listen("0.0.0.0:" + port))
}