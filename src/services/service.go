package services

import (
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
