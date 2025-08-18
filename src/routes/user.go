package routes

import (
	"mirabilis-api/src/controllers"
	"mirabilis-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	group := router.Group("/api/v1/users")
	
	jwtMiddleware :=  middlewares.JWTMiddleware([]string{"user"},[]string{"id"})

	group.GET("/profile",jwtMiddleware,controllers.Profile)
	group.GET("/login", middlewares.GetBasicAuthorization, controllers.Login)

}
