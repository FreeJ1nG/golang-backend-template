package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	mainDB *pgxpool.Pool
	rdb    *redis.Client
}

func NewRepository(mainDB *pgxpool.Pool, rdb *redis.Client) *repository {
	return &repository{
		mainDB: mainDB,
		rdb:    rdb,
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

func (r *repository) SetTokenBlacklistForUser(username string, tokenToInvalidate string, timeUntilExpiration time.Duration) {
	ctx := context.Background()
	r.rdb.Set(ctx, fmt.Sprintf("blacklisted:token:%s", username), tokenToInvalidate, timeUntilExpiration)
}

func (r *repository) GetBlacklistedTokenForUser(username string) (res *string, err error) {
	ctx := context.Background()
	result, err := r.rdb.Get(ctx, fmt.Sprintf("blacklisted:token:%s", username)).Result()
	if err == redis.Nil {
		res = nil
		err = fmt.Errorf("key does not exist for user %s", username)
		return
	}
	res = &result
	return
}
