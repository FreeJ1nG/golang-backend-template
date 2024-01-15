package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username     string `json:"username"`
	FullName     string `json:"fullName"`
	PasswordHash string `json:"passwordHash"`
}

func (u *User) ValidatePasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
