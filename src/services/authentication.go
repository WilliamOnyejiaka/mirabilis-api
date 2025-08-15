package services

import (
	"mirabilis-api/src/repos"
	"mirabilis-api/src/types"
	"net/http"
	// "mirabilis-api/src/services"
)

type AuthenticationService struct {
	repo *repos.UserRepository
	base BaseService
}

func NewAuthenticationService() *AuthenticationService {
	return &AuthenticationService{
		repo: repos.NewUserRepository(),
		base: BaseService{},
	}
}

func (this *AuthenticationService) GetName() types.ServiceResponse { // public method
	return this.base.ServiceResponse(http.StatusOK,false,"","Hello")
}
