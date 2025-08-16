package services

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"mirabilis-api/src/types"
)

type BaseService struct{}

// ServiceResponse creates a JSON response string
func (this *BaseService) ServiceResponse(
	statusCode int,
	isError bool,
	message string,
	data any, // any type, can be nil
) types.ServiceResponse {

	response := types.ServiceResponse{
		StatusCode: statusCode,
		JSON: map[string]any{
			"error":   isError,
			"message": message,
			"data":    data,
		},
	}

	return response
}

func (this *BaseService) HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return string(passwordHash), nil
}

func (this *BaseService) ComparePassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return false
	}
	return true
}
