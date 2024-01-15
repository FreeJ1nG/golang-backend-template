package interfaces

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/golang-jwt/jwt/v4"
)

type AuthRespository interface {
	GetUserByUsername(username string) (user models.User, err error)
}

type AuthService interface {
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
