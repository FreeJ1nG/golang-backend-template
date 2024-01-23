package util

import (
	"context"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/interfaces"
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

func (rp *routeProtector) Wrapper(f http.HandlerFunc) http.HandlerFunc {
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
		username := claims["sub"].(string)
		ctx := context.WithValue(r.Context(), UserContextKey, username)
		f(w, r.WithContext(ctx))
	}
}
