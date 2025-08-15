package main

import (
	"context" // Provides Context for deadlines/cancelation across API boundaries.
	"fmt"
	"log"

	// "os"

	// "mirabilis-api/src/middlewares"
	"mirabilis-api/src/config"
	"mirabilis-api/src/controllers"
	"mirabilis-api/src/models"
	"mirabilis-api/src/repos"
	"mirabilis-api/src/routes"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"time"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// init() runs automatically when main is ran.
func init() {
	config.LoadEnvVariables()
	config.InitializeDB()
}

func test() {
	// Create a new MongoDB client configured with a connection URI.
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil { // Always check errors returned by driver calls.
		log.Fatal(err) // log.Fatal logs the error and exits the program.
	}

	// Create a Context with a 10-second timeout for the upcoming connect call.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context's resources are freed when main() returns.

	// Establish the connection to the MongoDB server using the client and context.
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	// Ensure we cleanly close the connection when main() finishes.
	defer client.Disconnect(ctx)

	// Select the database ("testdb") and collection ("users") you want to work with.
	collection := client.Database("testdb").Collection("users")

	// Insert a document into the "users" collection.
	// Using a map as the BSON document; bson.M or a struct are also common.
	res, err := collection.InsertOne(ctx, map[string]interface{}{
		"name": "Alice", // Field "name" with string value.
		"age":  25,      // Field "age" with numeric value.
	})
	if err != nil { // Check if the insert failed.
		log.Fatal(err)
	}

	// Print the auto-generated _id of the inserted document (typically a primitive.ObjectID).
	fmt.Println("Inserted document ID:", res.InsertedID)
}

// * Program entry function.
func main() {
	// Initialize Gin router
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	router.Use(cors.New(corsConfig))
	// Apply middleware
	// router.Use(middlewares.Logger())

	// Setup routes
	// handlers.SetupRoutes(r, client, cfg)

	router.GET("/api", func(ctx *gin.Context) {
		repo := repos.NewUserRepository()

		user := models.User{
			Name:      "Alice",
			Password:  "Password",
			Email:     "email@email.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		InsertedID, err := repo.Insert(user)
		if err != nil {
			log.Fatal("Insert error:", err)
		}

		// idStr := "689e754057a7319195552dfe"

		// Convert string to ObjectID
		// objID, err := primitive.ObjectIDFromHex(idStr)
		// if err != nil {
		// 	log.Fatalf("Invalid ObjectID: %v", err)
		// }

		// user1, err1 := repos.FindOneByID(objID)

		// if err1 != nil {
		// 	log.Fatal("FindAll error:", err)
		// }

		user.ID = InsertedID

		ctx.IndentedJSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "user has been created",
			// "users":   user1,
			"InsertedID": user,
		})
	})

	//* Controllers
	routes.Auth(router)

	//* Error controllers
	controllers.MethodNotAllowed(router)
	controllers.RouteNotFound(router)

	// Start server
	port := ":" + config.Envs("port")
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
