package handlers

import (
	"context"
	"ecoApi/config"
	"ecoApi/models"
	"fmt"
	// "time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectionStatus(c *fiber.Ctx) error {
	fmt.Println("Well done young master")
    return c.SendString("Successful connection")
}

func CreatUser(c *fiber.Ctx) error {
	collection := config.MongoDatabase.Collection("users")
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse user JSON",
		})
	}

	if user.FBID == "" || user.Username == "" || user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields (uid, username, email)",
		})
	}

	user.UserID = primitive.NewObjectID()
	user.MedalTrans = 0
	user.MedalEnergy = 0
	user.MedalConsume = 0
	user.MedalDesecho = 0
	user.NotificationArray = []models.Notification{}
	user.Glass = 0
	user.Metal = 0
	user.Paper = 0
	user.Tetrapack = 0 
	user.Plastic = 0
	user.Cardboard = 0

	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to insert user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func GetUserByFBID(c *fiber.Ctx) error {
	uid := c.Params("uid")

	// filtro de busqueda
	filter := bson.M{"FBID": uid}

	uCollection := config.MongoDatabase.Collection("users")
	var user models.User

	err := uCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if len(user.NotificationArray) == 0 {
		user.NotificationArray = []models.Notification{}
	}


	return c.JSON(fiber.Map{
		"user": user,
	})
}

func GetUsers(c *fiber.Ctx) error {
	collection := config.MongoDatabase.Collection("users")

	filter := bson.D{}
	findOptions := options.Find()

	var users []models.User

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error retrieving users",
        })
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error decoding user data",
			})
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No users found",
		})
	}

	return c.JSON(users)
}


