package middlewares

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetBasicAuthorization(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")

	if authHeader == "" {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Authorization header missing",
		})
		context.Abort()
		return
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Basic Authorization header required",
		})
		context.Abort()
		return
	}

	authValue := strings.TrimPrefix(authHeader, "Basic ")
	authBytes, err := base64.StdEncoding.DecodeString(authValue)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Invalid Authorization header",
		})
		context.Abort()
		return
	}

	authStr := string(authBytes)
	credentials := strings.SplitN(authStr, ":", 2)
	if len(credentials) != 2 {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Invalid Authorization header",
		})
		context.Abort()
		return
	}

	username := credentials[0]
	password := credentials[1]

	context.Set("username", username)
	context.Set("password", password)
	context.Next()
}
