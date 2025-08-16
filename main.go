package main

import (
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
)

// init() runs automatically when main is ran.
func init() {
	config.LoadEnvVariables()
	config.InitializeDB()
}

// * Program entry function.
func main() {
	if config.Envs("environment") == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode) // 👈 switch to release mode
	}

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
	routes.AuthRoute(router)
	routes.UserRoute(router)

	//* Error controllers
	controllers.MethodNotAllowed(router)
	controllers.RouteNotFound(router)

	// Start server
	port := ":" + config.Envs("port")
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
