package controllers

import (
	"github.com/gin-gonic/gin"
	"mirabilis-api/src/services"
)

func Profile(ctx *gin.Context) {
	// Get data from context
	val, _ := ctx.Get("data")

	data, _ := val.(map[string]interface{})
	// if !ok {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid context data"})
	// 	return
	// }

	id, _ := data["id"].(string)
	// if !ok {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
	// 	return
	// }

	service := services.NewUserService()
	serviceResult := service.Profile(id)

	ctx.IndentedJSON(serviceResult.StatusCode, serviceResult.JSON)
}
