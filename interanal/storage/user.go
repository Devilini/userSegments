package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"userSegments/interanal/apperror"
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
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id=$1", usersTable)
	var user model.User
	err := s.client.QueryRow(ctx, query, id).Scan(&user.Id, &user.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, apperror.NotFoundError("not found user")
		}
	}

	return user, err
}

func (s *UserStorage) CreateUser(ctx context.Context, name string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name) values ($1) RETURNING id", usersTable)
	row := s.client.QueryRow(ctx, query, name)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
