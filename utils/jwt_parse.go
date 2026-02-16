package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type SignedDetails struct {
	Id 			string
	Firstname 	string
	Lastname  	string
	Password  	string
	Email     	string
	Role      	string
	jwt.RegisteredClaims
}

func GenerateJwt (id uuid.UUID, firstname string, lastname string, password string, email string, role string) (string, string, error) {
	
	if err := godotenv.Load(); err != nil {
		return "", "", errors.New("Failed to load env as you want!")
	}
	token_not_refresh := os.Getenv("JWT_SECRET_KEY")

	signed_details_not_refresh := &SignedDetails{
		Id: id.String(),
		Firstname: firstname,
		Lastname: lastname,
		Password: password,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer: token_not_refresh,
			IssuedAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			
		},
	}

	token_not_refresh_final := jwt.NewWithClaims(jwt.SigningMethodHS256, signed_details_not_refresh)
	token_not_refresh_final_1, err := token_not_refresh_final.SignedString([]byte(token_not_refresh))
	if err != nil {
		return "", "", errors.New("Failed to signed the data of the json web token!" + err.Error())
	}

	token_refresh := os.Getenv("JWT_SECRET_KEY_REFRESH_TOKEN")

	signed_details_refresh := &SignedDetails{
		Id: id.String(),
		Firstname: firstname,
		Lastname: lastname,
		Password: password,
		Email: email,
		Role: role, 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			Issuer: token_refresh,
			IssuedAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	token_refresh_final := jwt.NewWithClaims(jwt.SigningMethodHS256, signed_details_refresh)
	token_refresh_final_1, err  := token_refresh_final.SignedString([]byte(token_refresh))
	if err != nil {
		return "", "", errors.New("Failed to signed the data of json web token!" + err.Error())
	}

	return token_not_refresh_final_1, token_refresh_final_1, nil

}

func ValidateToken (tokenAuth string) (*SignedDetails, error) {

	claims := &SignedDetails{}

	if err := godotenv.Load(); err != nil {
		return nil, errors.New("Failed to load env as you want!")
	}
	token := os.Getenv("JWT_SECRET_KEY")
	token_refresh_env := os.Getenv("JWT_SECRET_KEY_REFRESH")

	token_not_refresh, err := jwt.ParseWithClaims(tokenAuth, claims, func(t *jwt.Token) (any, error) {
		return []byte(token), nil
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	token_refresh, err := jwt.ParseWithClaims(tokenAuth, claims, func(t *jwt.Token) (any, error) {
		return []byte(token_refresh_env), nil
	})

	if _, ok := token_not_refresh.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("Failed to convert data!")
	}

	if _, ok := token_refresh.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("Failed to convert data!")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("expired is true!")
	}

	return claims, nil

}