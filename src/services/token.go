package services

import (
	"fmt"
	"log"
	"mirabilis-api/src/config"
	"time"

	"golang.org/x/exp/slices"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secretKey string
}

func NewTokenService() *TokenService {
	return &TokenService{
		secretKey: config.Envs("secretKey"),
	}
}

func (this *TokenService) CreateToken(roles []string, data any) (string, error) {
	//* Create a new token with claims

	var years time.Duration = 100
	var days time.Duration = 365
	var hours time.Duration = 24

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data":  data,
		"roles": roles,
		"exp":   time.Now().Add(years * days * hours * time.Hour).Unix(), //* This token lasts for 100years
	})

	//* Sign the token
	tokenString, err := token.SignedString([]byte(this.secretKey)) //* ✅ cast to []byte
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, nil
}

func (this *TokenService) ParseToken(roles []string, tokenString string, neededValues []string) (map[string]any, error) {
	//* Parse and validate the token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//* Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(this.secretKey), nil // ✅ correct for HS256
	})

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error parsing token")
	}

	//* Prepare output map
	data := make(map[string]any)

	//* Extract claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		//* Check if 'exp' claim exists and is valid
		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if time.Now().After(expTime) {
				return nil, fmt.Errorf("Token has expired")
			}
		} else {
			return nil, fmt.Errorf("Missing or invalid 'exp' claim")
		}

		if !slices.Contains(roles, "any") { //? Check if this can be a separate function
			lookup := make(map[string]bool)
			for _, s := range roles {
				lookup[s] = true
			}

			if tokenRoles, ok := claims["roles"].([]any); ok {
				found := false
				for _, s := range tokenRoles {
					if lookup[s.(string)] {
						found = true
						break
					}
				}

				if !found {
					return nil, fmt.Errorf("User not authorized")
				}
			} else {
				return nil, fmt.Errorf("Missing or invalid 'role' claim")
			}
		}

		//* Extract 'data' claim
		if jwtData, ok := claims["data"].(map[string]interface{}); ok {
			for _, key := range neededValues {
				if val, exists := jwtData[key]; exists {
					data[key] = val
				}
			}
		}
	}

	return data, nil
}
