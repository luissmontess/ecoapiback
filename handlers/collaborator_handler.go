package handlers

import (
	"context"
	"ecoApi/config"
	"ecoApi/models"
	"time"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCollaborator(c *fiber.Ctx) error {
	collection := config.MongoDatabase.Collection("collaborators")
	var collaborator models.Collaborator

	if err := c.BodyParser(&collaborator); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse collaborator JSON",
		})
	}

	if collaborator.FBID == "" || collaborator.Username == "" || collaborator.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields (uid, username, email)",
		})
	}

	collaborator.CollaboratorID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.TODO(), collaborator)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to insert collaborator",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(collaborator)
}

func GetCollaboratorByFBID(c *fiber.Ctx) error {
	uid := c.Params("uid")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// search filter
    filter := bson.M{"FBID": uid}

	cCollection := config.MongoDatabase.Collection("collaborators")
	var collaborator models.Collaborator


	err := cCollection.FindOne(ctx, filter).Decode(&collaborator)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "collaborator not found",
		})
	}


	return c.JSON(fiber.Map{
		"collaborator": collaborator,
	})
}

func GetCollaborators(c *fiber.Ctx) error {
	// db collection
	collaboratorsCollection := config.MongoDatabase.Collection("collaborators")

	// arreglo de colaboradores
	var collaborators []models.Collaborator

	// encontrar todos
	cursor, err := collaboratorsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve collaborators",
		})
	}
	defer cursor.Close(context.TODO())

	// decode con cursor
	if err := cursor.All(context.TODO(), &collaborators); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding collaborators",
		})
	}

	return c.JSON(collaborators)
}