package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"userSegments/interanal/model"
)

type User interface {
	GetUserById(ctx context.Context, id int) (model.User, error)
	CreateUser(ctx context.Context, name string) (int, error)
}

type UserStorage struct {
	client *pgxpool.Pool
}

func NewUserStorage(client *pgxpool.Pool) UserStorage {
	return UserStorage{client: client}
}

func (s *UserStorage) GetUserById(ctx context.Context, id int) (model.User, error) {
	query := `SELECT id, name FROM users WHERE id=$1`
	var user model.User
	err := s.client.QueryRow(ctx, query, id).Scan(&user.Id, &user.Name)
	if err != nil {
		return user, fmt.Errorf("unable to query user: %w", err)
	}

	return user, nil
}

func (s *UserStorage) CreateUser(ctx context.Context, name string) (int, error) {
	var id int
	query := "INSERT INTO users (name) values ($1) RETURNING id" // todo table name
	row := s.client.QueryRow(ctx, query, name)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
