package internal

import (
	"context"
	"database/sql"

	"littlerollingsushi.com/example/entity"
)

const (
	GetUserByEmailQuery = "SELECT first_name, last_name, email, crypted_password FROM user WHERE email = ?"
)

type GetUserByEmailGateway struct {
	sql *sql.DB
}

func NewGetUserByEmailGateway(sql *sql.DB) *GetUserByEmailGateway {
	return &GetUserByEmailGateway{sql: sql}
}

func (g *GetUserByEmailGateway) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	user := entity.User{}
	err := g.sql.QueryRowContext(ctx, GetUserByEmailQuery, email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.CryptedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserNotFound
		}

		return user, err
	}

	return user, nil
}
