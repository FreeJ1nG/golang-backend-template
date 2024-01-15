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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := rp.authUtil.ToJwtToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "unable to get token claims", http.StatusInternalServerError)
			return
		}
		username := claims["sub"].(string)
		user, status, err := rp.authService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		f(w, r.WithContext(ctx))
	}
}
