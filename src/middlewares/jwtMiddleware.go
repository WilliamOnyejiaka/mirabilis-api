package middlewares

import (
	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
	"mirabilis-api/src/services"
	"net/http"
	"strings"
)

func JWTMiddleware(role []string, neededValues []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		//* Must start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Authorization header missing or invalid",
			})
			return
		}

		//* Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		tokenService := services.NewTokenService()
		data, err := tokenService.ParseToken("user", token, neededValues)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		
		ctx.Set("data", data)

		ctx.Next()
	}
}
