package services

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secretKey string
}

func NewTokenService() *TokenService {
	return &TokenService{
		secretKey: "helloworld",
	}
}

func (this *TokenService) CreateToken(role string, data any) (string, error) {
	//* Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"role": role,
		"exp":  time.Now().Add(100 * 365 * 24 * time.Hour).Unix(),
	})

	//* Sign the token
	tokenString, err := token.SignedString([]byte(this.secretKey)) // ✅ cast to []byte
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, nil
}

func (this *TokenService) ParseToken(role string,tokenString string, neededValues []string) (map[string]any, error) {
	//* Parse and validate the token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//* Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(this.secretKey), nil // ✅ correct for HS256
	})

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error parsing token")
	}

	//* Prepare output map
	data := make(map[string]any)

	//* Extract claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		//* Check if 'exp' claim exists and is valid
		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if time.Now().After(expTime) {
				return nil, fmt.Errorf("token has expired")
			}
		} else {
			return nil, fmt.Errorf("missing or invalid 'exp' claim")
		}

		if tokenRole, ok := claims["role"].(string); ok {
			if tokenRole != role {
				return nil, fmt.Errorf("user not authorized")
			}
		} else {
			return nil, fmt.Errorf("missing or invalid 'role' claim")
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