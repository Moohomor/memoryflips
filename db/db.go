package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Id       int
	Name     string
	Password string
}

type Database interface {
	GetUserByID(ctx context.Context, id int) (*User, error)
}

type Service struct {
	db *pgxpool.Pool
}

func (s *Service) GetUserByName(ctx context.Context, name string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(ctx, "SELECT id, username, password FROM users WHERE username = $1", name).
		Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
