package handlers

import (
	"context"
	"ecoApi/config"
	"ecoApi/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"fmt"

	"time"
)

func CreateRecollect(c *fiber.Ctx) error {
	var recollect models.Recollect
	collection := config.MongoDatabase.Collection("recollects")

	if err := c.BodyParser(&recollect); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse recollect JSON",
		})
	}

	if recollect.CollaboratorFBID == "" || recollect.Longitude == -1 || recollect.Latitude == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Longitude must be between greater than 0",
		})
	}

	recollect.RecollectID = primitive.NewObjectID()
	recollect.Cardboard = 0
	recollect.Glass = 0
	recollect.Paper = 0
	recollect.Plastic = 0
	recollect.Tetrapack = 0
	recollect.Metal = 0
	recollect.DonationArray = []models.Donation{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, recollect)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create recollect",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(recollect)
}

func GetRecollectByID(c *fiber.Ctx) error {
	uid := c.Params("recollect_id")
	recollectObjectID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid recollect ID",
		})
	}

	filter := bson.M{"RecollectID": recollectObjectID}

	var recollect models.Recollect
	rCollection := config.MongoDatabase.Collection("recollects")

	err = rCollection.FindOne(context.TODO(), filter).Decode(&recollect)
	if err != nil {
		fmt.Printf("Error finding recollect with ID %v: %v\n", uid, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "recollect not found",
			"details": err.Error(),
		})
	}

	return c.JSON(recollect)
}


func GetActiveRecollects(c *fiber.Ctx) error {
	currentTime := primitive.NewDateTimeFromTime(time.Now())

	filter := bson.M{
		"StartTime": bson.M{"$lte": currentTime}, // parametrizacion por tiempo
		"EndTime":   bson.M{"$gte": currentTime}, 
	}

	recollectCollection := config.MongoDatabase.Collection("recollects")
	var recollects []models.Recollect

	cursor, err := recollectCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving active recollects",
		})
	}

	if err := cursor.All(context.TODO(), &recollects); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding active recollects",
		})
	}

	return c.JSON(recollects)
}

func GetCollaboratorRecollects(c *fiber.Ctx) error {
	uid := c.Params("uid")
	filter := bson.M{"CollaboratorFBID" : uid}
	var recollects []models.Recollect
	
	recollectCollection := config.MongoDatabase.Collection("recollects")
	recollectCursor, err := recollectCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error retrieving recollects",
        })
	}
	defer recollectCursor.Close(context.TODO())

	if err := recollectCursor.All(context.TODO(), &recollects); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding recollects",
		})
	}

	return c.JSON(recollects)
}

func GetActiveCollaboratorRecollects(c *fiber.Ctx) error {
	uid := c.Params("uid")
	currentTime := primitive.NewDateTimeFromTime(time.Now())

	filter := bson.M{
		"CollaboratorFBID": uid,
		"StartTime":        bson.M{"$lte": currentTime},
		"EndTime":          bson.M{"$gte": currentTime},
	}

	recollectCollection := config.MongoDatabase.Collection("recollects")
	var recollects []models.Recollect

	cursor, err := recollectCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving active recollects for collaborator",
		})
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &recollects); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding active recollects for collaborator",
		})
	}

	return c.JSON(recollects)
}

func AddtoRecollect(c *fiber.Ctx) error {
	type Request struct {
		UserFBID  string `json:"UserFBID"`
		Cardboard int    `json:"Cardboard"`
		Glass     int    `json:"Glass"`
		Tetrapack int    `json:"Tetrapack"`
		Plastic   int    `json:"Plastic"`
		Paper     int    `json:"Paper"`
		Metal     int    `json:"Metal"`
	}

	recollectID := c.Params("recollect_id")
	if recollectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Recollect ID is required",
		})
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	recollectObjectID, err := primitive.ObjectIDFromHex(recollectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid recollect ID",
		})
	}

	recollectCollection := config.MongoDatabase.Collection("recollects")
	userCollection := config.MongoDatabase.Collection("users")

	// checar por existencia de usuario
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"FBID": req.UserFBID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var recollect models.Recollect
	err = recollectCollection.FindOne(context.TODO(), bson.M{"RecollectID": recollectObjectID}).Decode(&recollect)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Recollect not found",
		})
	}

	totalKg := recollect.Cardboard + recollect.Glass + recollect.Paper + recollect.Tetrapack + recollect.Plastic + recollect.Metal
	if(totalKg) >= recollect.Limit {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Weight Limit has been reached",
		})
	}

	// update de datos
	filter := bson.M{"RecollectID": recollectObjectID}
	update := bson.M{
		"$inc": bson.M{
			"Cardboard": req.Cardboard,
			"Glass":     req.Glass,
			"Tetrapack": req.Tetrapack,
			"Plastic":   req.Plastic,
			"Paper":     req.Paper,
			"Metal":     req.Metal,
		},
		"$push": bson.M{
			"DonationArray": bson.M{
				"UserFBID":  user.FBID,
				"Username":  user.Username,
				"Cardboard": req.Cardboard,
				"Glass":     req.Glass,
				"Tetrapack": req.Tetrapack,
				"Plastic":   req.Plastic,
				"Paper":     req.Paper,
				"Metal":     req.Metal,
			},
		},
	}

	_, err = recollectCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update recollect",
		})
	}

	notification := models.Notification{
		NotificationType: "RECOLLECT",
		Message:   "Gracias por tu contribucion a la recolecta!",
	}

	userFilter := bson.M{"FBID": req.UserFBID}
	userUpdate := bson.M{
		"$push": bson.M{"NotificationArray": notification},
		"$inc": bson.M{
			"MedalConsume": 1, // incrementos
			"MedalDesecho": 1, 
		},
	}

	result, err := userCollection.UpdateOne(context.TODO(), userFilter, userUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add notification to user and increment medals",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No matching user found to update notification and medals",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Recollect updated, donator added, medals incremented, and notification created successfully",
	})
}
