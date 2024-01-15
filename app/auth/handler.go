package auth

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/dto"
	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/models"
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
		http.Error(w, err.Error(), status)
		return
	}
	util.EncodeResponse(w, res, status)
}

func (h *handler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := util.ParseRequestBody[dto.SignUpRequest](w, r)
	res, status, err := h.authService.SignUpUser(body.Username, body.FullName, body.Password)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	util.EncodeResponse(w, res, status)
}

func (h *handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := r.Context().Value(util.UserContextKey).(models.User)
	util.EncodeResponse(
		w,
		dto.GetCurrentUserResponse{
			Username: user.Username,
			FullName: user.FullName,
		},
		http.StatusOK,
	)
}
