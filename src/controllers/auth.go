package controllers

import (
	"mirabilis-api/src/services"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// * 1️⃣ Define your custom messages **before** using them
var validationMessages = map[string]map[string]string{
	"Email": {
		"required": "Email is required",
		"email":    "Email must be a valid email address",
	},
	"Password": {
		"required": "Password is required",
		"min":      "Password must be at least 6 characters",
	},
	"Name": {
		"required": "Name is required",
	},
}

func ParseValidationErrors(err error) map[string]string {
	res := make(map[string]string)
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()

			// Check if we have a custom message
			if msgs, ok := validationMessages[field]; ok {
				if msg, exists := msgs[tag]; exists {
					res[field] = msg
					continue
				}
			}

			// Fallback message
			res[field] = "Invalid value"
		}
	} else {
		// Not a validation error
		res["error"] = err.Error()
	}

	return res
}

func SignUp(ctx *gin.Context) {
	type jsonBody struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var body jsonBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		errorsMap := ParseValidationErrors(err)
		for _, msg := range errorsMap {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": msg})
			return
		}
	}

	if len(body.Password) < 5 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "password length should be greater than 4",
		})
		return
	}

	service := services.NewAuthenticationService()

	serviceResult := service.SignUp(body.Name, body.Email, body.Password)

	ctx.IndentedJSON(serviceResult.StatusCode, serviceResult.JSON)
}

// func Key(ctx *gin.Context) {
// 	password, _ := ctx.Get("password")
// 	email, _ := ctx.Get("username")

// 	if email.(string) == "" {
// 		ctx.IndentedJSON(500, gin.H{
// 			"error":   true,
// 			"message": "email empty",
// 		})
// 		return
// 	}

// 	if password.(string) == "" {
// 		ctx.IndentedJSON(500, gin.H{
// 			"error":   true,
// 			"message": "password empty",
// 		})
// 		return
// 	}

// 	serviceResult := services.Key(email.(string), password.(string), ctx)

// 	if serviceResult.Error {
// 		ctx.IndentedJSON(serviceResult.StatusCode, gin.H{
// 			"error":   serviceResult.Error,
// 			"message": serviceResult.Message,
// 		})
// 		return
// 	}

// 	ctx.IndentedJSON(serviceResult.StatusCode, gin.H{
// 		"error": serviceResult.Error,
// 		"data":  serviceResult.Data,
// 	})
// }
