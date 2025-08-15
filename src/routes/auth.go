package routes

import (
	"mirabilis-api/src/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.Engine) {
	group := router.Group("/api/v1/auth")

	group.POST("/sign-up", controllers.SignUp)
	// userGroup.GET("/key", middlewares.GetBasicAuthorization, handlers.Key)

	group.GET("/test", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "test route",
		})
	})
}
