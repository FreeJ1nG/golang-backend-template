package auth

import (
	"fmt"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/dto"
	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/models"
)

type service struct {
	authRepo interfaces.AuthRespository
	authUtil interfaces.AuthUtil
}

func NewService(authRepo interfaces.AuthRespository, authUtil interfaces.AuthUtil) *service {
	return &service{
		authRepo: authRepo,
		authUtil: authUtil,
	}
}

func (s *service) SignInUser(username string, password string) (res dto.SignInResponse, status int, err error) {
	status = http.StatusOK
	user, err := s.authRepo.GetUserByUsername(username)
	if err != nil {
		err = fmt.Errorf("unable to authenticate user: user with username of %s not found", username)
		status = http.StatusNotFound
		return
	}
	if !user.ValidatePasswordHash(password) {
		err = fmt.Errorf("unable to authenticate user: invalid password")
		status = http.StatusUnauthorized
		return
	}
	jwtToken, refreshToken, err := s.authUtil.GenerateTokenPair(user)
	if err != nil {
		err = fmt.Errorf("unable to generate token: %s", err.Error())
		status = http.StatusInternalServerError
		return
	}
	res = dto.SignInResponse{
		Token:        jwtToken,
		RefreshToken: refreshToken,
	}
	return
}

func (s *service) SignUpUser(username string, firstName string, lastName string, password string) (res dto.SignUpResponse, status int, err error) {
	status = http.StatusOK
	_, err = s.authRepo.GetUserByUsername(username)
	if err == nil {
		status = http.StatusForbidden
		err = fmt.Errorf("unable to sign user up: user with username %s already exists", username)
		return
	}
	passwordHash, err := s.authUtil.HashPassword(password)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to hash password: %s", err.Error())
		return
	}
	user, err := s.authRepo.CreateUser(username, firstName, lastName, passwordHash)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to create user: %s", err.Error())
		return
	}
	jwtToken, refreshToken, err := s.authUtil.GenerateTokenPair(user)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to generate token: %s", err.Error())
		return
	}
	res = dto.SignUpResponse{
		Token:        jwtToken,
		RefreshToken: refreshToken,
	}
	return
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
