package interfaces

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/dto"
	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/golang-jwt/jwt/v4"
)

type AuthRespository interface {
	CreateUser(username string, firstName string, lastName string, passwordHash string) (user models.User, err error)
	GetUserByUsername(username string) (user models.User, err error)
}

type AuthService interface {
	SignInUser(username string, password string) (res dto.SignInResponse, status int, err error)
	SignUpUser(username string, firstName string, lastName string, password string) (res dto.SignUpResponse, status int, err error)
	GetUserByUsername(username string) (user models.User, status int, err error)
	RefreshToken(refreshToken string) (res dto.RefreshTokenResponse, status int, err error)
}

type AuthHandler interface {
	SignInUser(w http.ResponseWriter, r *http.Request)
	SignUpUser(w http.ResponseWriter, r *http.Request)
	GetCurrentUser(w http.ResponseWriter, r *http.Request)
	RefreshJwt(w http.ResponseWriter, r *http.Request)
}

type AuthUtil interface {
	GenerateTokenPair(user models.User) (signedJwtToken string, signedRefreshToken string, err error)
	HashPassword(password string) (passwordHash string, err error)
	ExtractJwtToken(r *http.Request) (jwtToken string, err error)
	ToJwtToken(tokenString string) (token *jwt.Token, err error)
}
