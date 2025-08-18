package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mirabilis-api/src/services"
	"mirabilis-api/src/types"
	"net/http"
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
	type formBody struct {
		Name     string `form:"name" binding:"required"`
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required,min=6"`
	}

	var form formBody

	if err := ctx.ShouldBind(&form); err != nil {
		errorsMap := ParseValidationErrors(err)
		for _, msg := range errorsMap {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": msg})
			return
		}
	}

	image, _, err := ctx.Request.FormFile("image")
	var service = services.NewAuthenticationService()
	var serviceResult types.ServiceResponse

	if err != nil {
		serviceResult = service.SignUp(form.Name, form.Email, form.Password, ctx, nil)
	} else {
		serviceResult = service.SignUp(form.Name, form.Email, form.Password, ctx, image)
		defer image.Close()
	}

	ctx.IndentedJSON(serviceResult.StatusCode, serviceResult.JSON)
}

func Login(ctx *gin.Context) {
	password, _ := (ctx.Get("password"))
	email, _ := ctx.Get("username")

	if email.(string) == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Email is empty",
		})
		return
	}

	if password.(string) == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Password is empty",
		})
		return
	}

	var service = services.NewAuthenticationService()

	serviceResult := service.Login(email.(string), password.(string))

	ctx.IndentedJSON(serviceResult.StatusCode, serviceResult.JSON)
}
