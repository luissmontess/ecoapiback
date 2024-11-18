package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func ConnectMongoDB() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	MongoClient = client
	MongoDatabase = MongoClient.Database(dbName)
	
	log.Println("Ou yea, conectamos a mongo!")
}
