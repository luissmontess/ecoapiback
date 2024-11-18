package handlers

import (
	"context"
	"ecoApi/config"
	"ecoApi/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func TopTen() ([]models.TopTenUser, error) {
	collection := config.MongoDatabase.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: sumar campos de medallas
		{
			{Key: "$addFields", Value: bson.D{
				{Key: "TotalMedals", Value: bson.D{
					{Key: "$add", Value: bson.A{"$MedalTrans", "$MedalEnergy", "$MedalConsume", "$MedalDesecho"}},
				}},
			}},
		},
		// Stage 2: sortear en orden descendiente
		{
			{Key: "$sort", Value: bson.D{{Key: "TotalMedals", Value: -1}}},
		},
		// Stage 3: limitar a 10
		{
			{Key: "$limit", Value: 10},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error during aggregation: %v", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error decoding aggregation results: %v", err)
	}

	var topTenUsers []models.TopTenUser
	for _, user := range users {
		var totalMedalCount = user.MedalConsume + user.MedalDesecho + user.MedalEnergy + user.MedalTrans
		topTenUsers = append(topTenUsers, models.TopTenUser{
			UserFBID: user.FBID,
			Username: user.Username,
			Email:    user.Email,
			Place:    totalMedalCount,
		})
	}

	return topTenUsers, nil
}

func UpdateTopTen() error {
	topTenUsers, err := TopTen()
	if err != nil {
		return fmt.Errorf("failed to get top ten users: %v", err)
	}

	topTenArray := models.TopTenArray{
		UserArray: topTenUsers,
	}

	collection := config.MongoDatabase.Collection("topten")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "UserArray", Value: topTenArray.UserArray},
	}}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update top ten array: %v", err)
	}

	return nil
}

func GetTopTen(c *fiber.Ctx) error {
	collection := config.MongoDatabase.Collection("topten")

	var toptenarray models.TopTenArray
	err := collection.FindOne(context.TODO(), bson.M{}).Decode(&toptenarray)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving top ten users",
		})
	}

	return c.JSON(toptenarray)
}