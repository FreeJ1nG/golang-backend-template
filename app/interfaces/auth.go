package interfaces

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/dto"
	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/golang-jwt/jwt/v4"
)

type AuthRespository interface {
	CreateUser(username string, fullName string, passwordHash string) (user models.User, err error)
	GetUserByUsername(username string) (user models.User, err error)
}

type AuthService interface {
	SignInUser(username string, password string) (res dto.SignInResponse, status int, err error)
	SignUpUser(username string, fullName string, password string) (res dto.SignUpResponse, status int, err error)
	GetUserByUsername(username string) (user models.User, status int, err error)
}

type AuthHandler interface {
}

type AuthUtil interface {
	GenerateToken(user models.User) (signedToken string, err error)
	HashPassword(password string) (passwordHash string, err error)
	ExtractJwtToken(r *http.Request) (jwtToken string, err error)
	ToJwtToken(tokenString string) (token *jwt.Token, err error)
}
