package auth

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/dto"
	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/util"
)

type handler struct {
	authService interfaces.AuthService
}

func NewHandler(authService interfaces.AuthService) *handler {
	return &handler{
		authService: authService,
	}
}

func (h *handler) SignInUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := util.ParseRequestBody[dto.SignInRequest](w, r)
	res, status, err := h.authService.SignInUser(body.Username, body.Password)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}
	util.EncodeSuccessResponse(w, res, status)
}

func (h *handler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := util.ParseRequestBody[dto.SignUpRequest](w, r)
	res, status, err := h.authService.SignUpUser(body.Username, body.FirstName, body.LastName, body.Password)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}
	util.EncodeSuccessResponse(w, res, status)
}

func (h *handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.Context().Value(util.UserContextKey).(string)
	user, status, err := h.authService.GetUserByUsername(username)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}
	util.EncodeSuccessResponse(
		w,
		dto.GetCurrentUserResponse{
			Id:        user.Id,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		http.StatusOK,
	)
}

func (h *handler) RefreshJwt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := util.ParseRequestBody[dto.RefreshTokenRequest](w, r)
	res, status, err := h.authService.RefreshToken(body.RefreshToken)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}
	util.EncodeSuccessResponse(w, res, status)
}
