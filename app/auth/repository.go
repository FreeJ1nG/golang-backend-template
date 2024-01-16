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

func (r *repository) CreateUser(username string, firstName string, lastName string, passwordHash string) (user models.User, err error) {
	ctx := context.Background()
	_, err = r.mainDB.Exec(
		ctx,
		`INSERT INTO Account (username, first_name, last_name, password_hash) VALUES ($1, $2, $3, $4);`,
		username,
		firstName,
		lastName,
		passwordHash,
	)
	if err != nil {
		return
	}
	user = models.User{
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
	}
	return
}

func (r *repository) GetUserByUsername(username string) (user models.User, err error) {
	ctx := context.Background()

	err = pgxscan.Get(
		ctx,
		r.mainDB,
		&user,
		`SELECT * FROM Account WHERE username = $1;`,
		username,
	)

	if err != nil {
		err = fmt.Errorf("unable to get user by username: %s", err.Error())
		return
	}

	return
}
