package internal

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"littlerollingsushi.com/example/entity"
)

const (
	insertUserQuery = "INSERT INTO user (first_name, last_name, email, crypted_password, created_at) VALUES (?, ?, ?, ?, ?)"
)

type InsertUserGateway struct {
	sql *sql.DB
}

func NewInsertUserGateway(sql *sql.DB) *InsertUserGateway {
	return &InsertUserGateway{sql: sql}
}

func (g *InsertUserGateway) InsertUser(ctx context.Context, user entity.User) error {
	_, err := g.sql.ExecContext(ctx, insertUserQuery, user.FirstName, user.LastName, user.Email, user.CryptedPassword, time.Now().UTC())
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok {
			if me.Number == errNoDuplicateRecord {
				return ErrUserAlreadyExist
			}
		}

		return err
	}

	return nil
}
