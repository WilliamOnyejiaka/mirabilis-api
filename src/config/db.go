package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func InitializeDB() {
	//* Create a Context with a 10-second timeout for the upcoming connect call.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	databaseUrl := Envs("databaseUrl")
	databaseName := Envs("databaseName")

	//* Create a new MongoDB client configured with a connection URI.
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(databaseUrl))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	//* Select database
	DB = Client.Database(databaseName)

	log.Println("✅ Connected to MongoDB")
}
