package models

import "fmt"

type UserRole string

const (
	Admin  UserRole = "admin"
	Member UserRole = "member"
)

func ParseUserRole(s string) (res UserRole, err error) {
	switch s {
	case "admin":
		res, err = Admin, nil
	case "member":
		res, err = Member, nil
	default:
		err = fmt.Errorf("invalid user role: %s", s)
		return
	}
	return
}
