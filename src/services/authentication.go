package services

import (
	"log"
	"mirabilis-api/src/models"
	"mirabilis-api/src/repos"
	"mirabilis-api/src/types"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
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

func (this *AuthenticationService) SignUp(name string, email string, password string) types.ServiceResponse {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		log.Println(err.Error())
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Failed to hash password", nil)
	}

	userExists, err := this.repo.FindOneByEmail(email)

	if err != nil {
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong", nil)
	}

	if userExists != nil {
		return this.base.ServiceResponse(http.StatusBadRequest, true, "Email already exists", nil)
	}

	user := models.User{
		Name:      name,
		Password:  string(passwordHash),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	insertedID, err := this.repo.Insert(user)

	if err != nil {
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong", nil)
	}

	user.ID = insertedID
	tokenService := NewTokenService()
	token, err := tokenService.CreateToken("user", map[string]any{
		"id": insertedID.Hex(),
	})

	return this.base.ServiceResponse(http.StatusOK, false, "User has signed up successfully", map[string]any{
		"user":  user,
		"token": token,
	})
}

func (this *AuthenticationService) Login(email string, password string) types.ServiceResponse {
	user, err := this.repo.FindOneByEmail(email)

	if err != nil {
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong", nil)
	}

	if user == nil {
		return this.base.ServiceResponse(http.StatusNotFound, true, "User was not found", nil)
	}

	validPassword := this.base.ComparePassword(user.Password, password)

	if validPassword == true {
		tokenService := NewTokenService()
		token, err := tokenService.CreateToken("user", map[string]any{
			"id": user.ID.Hex(),
		})

		if err != nil {
			log.Println(err)
			return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong", nil)
		}

		return this.base.ServiceResponse(http.StatusOK, false, "User has logged in successfully", map[string]any{
			"user":  user,
			"token": token,
		})
	}

	return this.base.ServiceResponse(http.StatusBadRequest, true, "Invalid Password", nil)
}
