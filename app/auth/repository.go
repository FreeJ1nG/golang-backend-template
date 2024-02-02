package auth

import (
	"context"
	"fmt"

	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	mainDB *pgxpool.Pool
}

func NewRepository(mainDB *pgxpool.Pool) *repository {
	return &repository{
		mainDB: mainDB,
	}
}

func (r *repository) CreateUser(username string, firstName string, lastName string, passwordHash string, role models.UserRole) (user models.User, err error) {
	ctx := context.Background()

	row := r.mainDB.QueryRow(
		ctx,
		`INSERT INTO account (username, first_name, last_name, password_hash, role) VALUES ($1, $2, $3, $4, $5) RETURNING *;`,
		username,
		firstName,
		lastName,
		passwordHash,
		role,
	)

	err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.PasswordHash, &user.Role)
	return
}

func (r *repository) GetUserByUsername(username string) (user models.User, err error) {
	ctx := context.Background()

	err = pgxscan.Get(
		ctx,
		r.mainDB,
		&user,
		`SELECT * FROM account WHERE username = $1;`,
		username,
	)

	if err != nil {
		err = fmt.Errorf("unable to get user by username: %s", err.Error())
		return
	}

	return
}
