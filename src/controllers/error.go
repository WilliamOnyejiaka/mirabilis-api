package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MethodNotAllowed(router *gin.Engine) {
	router.HandleMethodNotAllowed = true
	router.NoMethod(func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusMethodNotAllowed, gin.H{
			"error":   true,
			"message": "Method not allowed",
		})
	})
}

func RouteNotFound(router *gin.Engine) {
	router.NoRoute(func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Route was not found.Please check the documentations.",
		})
	})
}
