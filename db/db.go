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

type Word struct {
	Id  int
	Rus string
	Eng string
}

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) GetWord(ctx context.Context, user User) (*Word, error) {
	word := &Word{}

	// check if all word are learned
	rows, err := s.db.Query(ctx, "SELECT * FROM schema.words WHERE eng NOT IN (SELECT word FROM schema.learned WHERE user_id = $1) LIMIT 1;", user.Id)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}

	err = s.db.QueryRow(ctx, "SELECT id, eng, rus FROM schema.words WHERE eng NOT IN (SELECT word FROM schema.learned WHERE user_id = $1) ORDER BY random() LIMIT 1;", user.Id).
		Scan(&word.Id, &word.Eng, &word.Rus)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (s *Service) GetWordById(ctx context.Context, id int) (*Word, error) {
	word := &Word{}
	err := s.db.QueryRow(ctx, "SELECT id, eng, rus FROM schema.words WHERE id = $1;", id).
		Scan(&word.Id, &word.Eng, &word.Rus)
	if err != nil {
		return nil, err
	}
	return word, nil
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
	err := s.db.QueryRow(ctx, "INSERT INTO schema.users (username, password) VALUES ($1,$2) RETURNING id", user.Name, user.Password).
		Scan(&user.Id)
	fmt.Println(user)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func (s *Service) MarkWordAsLearned(ctx context.Context, user User, word Word) error {
	_, err := s.db.Exec(ctx, "INSERT INTO schema.learned(word, user_id) VALUES ($1, $2);", word.Eng, user.Id)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}
