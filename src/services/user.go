package services

import (
	"log"
	"mirabilis-api/src/repos"
	"mirabilis-api/src/types"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo *repos.UserRepository
	base BaseService
}

func NewUserService() *UserService {
	return &UserService{
		repo: repos.NewUserRepository(),
		base: BaseService{},
	}
}

func (this *UserService) Profile(id string) types.ServiceResponse {
	//* Convert string to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong", nil)
	}

	user, err := this.repo.FindOneByID(objID)

	if err != nil {
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong", nil)
	}

	if user == nil {
		return this.base.ServiceResponse(http.StatusNotFound, true, "User was not found", nil)
	}

	return this.base.ServiceResponse(http.StatusOK, false, "User profile was retrieved successfully", user)
}
