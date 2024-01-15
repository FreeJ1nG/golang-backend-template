package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type util struct {
}

func NewUtil() *util {
	return &util{}
}

func (u *util) GenerateToken(user models.User) (signedToken string, err error) {
	now := time.Now()
	expiryInHours := viper.GetInt("JWT_EXPIRY_IN_HOURS")
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    viper.GetString("APP_NAME"),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryInHours) * time.Hour)),
			Subject:   user.Username,
		},
	)
	signedToken, err = token.SignedString([]byte(viper.GetString("JWT_SECRET_KEY")))
	return
}

func (u *util) HashPassword(password string) (passwordHash string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	passwordHash = string(hashedPassword)
	return
}

func (u *util) ExtractJwtToken(r *http.Request) (jwtToken string, err error) {
	authorization := r.Header.Get("Authorization")
	authSplit := strings.Split(authorization, " ")
	if len(authSplit) != 2 {
		err = fmt.Errorf("invalid authorization header format")
		return
	}
	prefix := authSplit[0]
	tokenString := authSplit[1]
	if prefix != "Bearer" {
		err = fmt.Errorf("jwt token not found on authorization header")
		return
	}
	jwtToken = tokenString
	return
}

func (u *util) ToJwtToken(tokenString string) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method)
		}
		return []byte(viper.GetString("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		err = fmt.Errorf("unable to parse token: %s", err.Error())
		return
	}
	return
}
