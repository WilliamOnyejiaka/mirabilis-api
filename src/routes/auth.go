package routes

import (
	"github.com/gin-gonic/gin"
	"mirabilis-api/src/controllers"
	"mirabilis-api/src/middlewares"
	"net/http"
)

func AuthRoute(router *gin.Engine) {
	group := router.Group("/api/v1/auth")

	group.POST("/sign-up", controllers.SignUp)
	group.GET("/login", middlewares.GetBasicAuthorization, controllers.Login)

	group.GET("/test", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "test route",
		})
	})
}
