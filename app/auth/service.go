package auth

import (
	"fmt"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/models"
)

type service struct {
	authRepo interfaces.AuthRespository
}

func NewService(authRepo interfaces.AuthRespository) *service {
	return &service{
		authRepo: authRepo,
	}
}

func (s *service) GetUserByUsername(username string) (user models.User, status int, err error) {
	status = http.StatusOK

	user, err = s.authRepo.GetUserByUsername(username)
	if err != nil {
		err = fmt.Errorf("unable to find user: %s", err.Error())
		status = http.StatusNotFound
		return
	}

	return
}
