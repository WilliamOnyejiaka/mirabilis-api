package services

import (
	"log"
	"mime/multipart"
	"mirabilis-api/src/models"
	"mirabilis-api/src/repos"
	"mirabilis-api/src/types"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
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

func (this *AuthenticationService) SignUp(name string, email string, password string, ctx *gin.Context, image any) types.ServiceResponse {
	userExists, err := this.repo.FindOneByEmail(email)

	if err != nil {
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong", nil)
	}

	if userExists != nil {
		return this.base.ServiceResponse(http.StatusBadRequest, true, "Email already exists", nil)
	}

	var result map[string]any

	if image != nil {
		folder, resourceType := "mirabilis-cdn/profile-pictures", "image"
		cloudinaryService := NewCloudinaryService()
		resultChan := cloudinaryService.UploadFile(ctx, image.(multipart.File), folder, resourceType)
		result = <-resultChan // Wait for upload to finish
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		log.Println(err.Error())
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Failed to hash password", nil)
	}

	user := models.User{
		Name:      name,
		Password:  string(passwordHash),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if image != nil {
		uploadResult, uploadErr := result["data"].(*uploader.UploadResult), result["error"].(bool)
		if uploadErr {
			return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong, failed to upload image", nil)
		}

		user.ImageURL = uploadResult.URL
		user.ImagePublicID = uploadResult.PublicID
	}

	insertedID, err := this.repo.Insert(user)

	if err != nil {
		return this.base.ServiceResponse(http.StatusInternalServerError, true, "Something went wrong", nil)
	}

	user.ID = insertedID
	tokenService := NewTokenService()
	token, err := tokenService.CreateToken([]string{"user"}, map[string]any{
		"id": insertedID.Hex(),
	})
	user.Password = ""
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
		token, err := tokenService.CreateToken([]string{"user"}, map[string]any{
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
