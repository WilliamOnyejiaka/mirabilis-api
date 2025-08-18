package routes

import (
	"mirabilis-api/src/controllers"
	"mirabilis-api/src/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {
	group := router.Group("/api/v1/auth")

	const MaxFileSize = 4 << 20 // 4MB in bytes
	validateImage := middlewares.ValidateFile(
		[]string{"image/jpeg", "image/png", "image/jpg"},
		MaxFileSize,
		"image",
		false,
	)

	group.POST("/sign-up", validateImage, controllers.SignUp)
	group.GET("/login", middlewares.GetBasicAuthorization, controllers.Login)

	group.GET("/test", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "test route",
		})
	})
}
