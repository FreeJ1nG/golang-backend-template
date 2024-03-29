package util

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/golang-jwt/jwt/v4"
)

type ContextKey string

var UserContextKey = ContextKey("user")

type routeProtector struct {
	authUtil    interfaces.AuthUtil
	authService interfaces.AuthService
}

func NewRouteProtector(authUtil interfaces.AuthUtil, authService interfaces.AuthService) *routeProtector {
	return &routeProtector{
		authUtil:    authUtil,
		authService: authService,
	}
}

func (rp *routeProtector) Wrapper(f http.HandlerFunc, adminOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := rp.authUtil.ExtractJwtToken(r)
		if err != nil {
			EncodeErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := rp.authUtil.ToJwtToken(tokenString)
		if err != nil {
			EncodeErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			EncodeErrorResponse(w, "unable to get token claims", http.StatusInternalServerError)
			return
		}
		tokenType := claims["typ"].(string)
		if tokenType != "access" {
			EncodeErrorResponse(w, "invalid token type, must be access token", http.StatusForbidden)
			return
		}
		parsedRole, err := models.ParseUserRole(claims["role"].(string))
		if err != nil {
			EncodeErrorResponse(w, "invalid token, role does not exist in payload", http.StatusBadRequest)
			return
		}
		fmt.Println(parsedRole)
		if adminOnly && parsedRole != models.Admin {
			EncodeErrorResponse(w, "unauthorized access to endpoint, must be admin", http.StatusUnauthorized)
			return
		}
		username := claims["sub"].(string)
		ctx := context.WithValue(r.Context(), UserContextKey, username)
		f(w, r.WithContext(ctx))
	}
}
