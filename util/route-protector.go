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

var TokenClaimsContextKey = ContextKey("token-claims")
var TokenContextKey = ContextKey("token")

type routeProtector struct {
	authUtil interfaces.AuthUtil
	authRepo interfaces.AuthRespository
}

func NewRouteProtector(authUtil interfaces.AuthUtil, authRepo interfaces.AuthRespository) *routeProtector {
	return &routeProtector{
		authUtil: authUtil,
		authRepo: authRepo,
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
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			EncodeErrorResponse(w, "unable to get token claims", http.StatusInternalServerError)
			return
		}
		claims, err := ConvertMapToTypeStruct[models.JwtClaims](mapClaims)
		if err != nil {
			EncodeErrorResponse(w, fmt.Sprintf("unable to convert map to type struct: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		tokenType := claims.TokenType
		if tokenType != "access" {
			EncodeErrorResponse(w, "invalid token type, must be access token", http.StatusForbidden)
			return
		}
		blacklistedToken, err := rp.authRepo.GetBlacklistedTokenForUser(claims.Subject)
		if err == nil && *blacklistedToken == tokenString {
			EncodeErrorResponse(w, "given token has been blacklisted, either from user or admin action", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), TokenClaimsContextKey, claims)
		ctx = context.WithValue(ctx, TokenContextKey, tokenString)
		f(w, r.WithContext(ctx))
	}
}
