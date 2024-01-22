package models

import "github.com/golang-jwt/jwt/v4"

type JwtClaims struct {
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}
