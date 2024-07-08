package db

import (
	"context"
	"fmt"

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

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) GetUserByName(ctx context.Context, name string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(ctx, "SELECT id, username, password FROM schema.users WHERE username = $1;", name).
		Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, user User) error {
	err := s.db.QueryRow(ctx, "insert into schema.users (username, password) values ($1,$2) returning id", user.Name, user.Password).
		Scan(&user.Id)
	fmt.Println(user)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}
