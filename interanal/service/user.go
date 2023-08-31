package service

import (
	"context"
	"userSegments/interanal/model"
	"userSegments/interanal/storage"
)

type User interface {
	GetUserById(ctx context.Context, id int) (model.User, error)
	CreateUser(ctx context.Context, name string) (int, error)
}

type UserService struct {
	userStorage storage.User
}

func NewUserService(storage storage.User) *UserService {
	return &UserService{userStorage: storage}
}

func (s *UserService) GetUserById(ctx context.Context, id int) (model.User, error) {
	user, err := s.userStorage.GetUserById(ctx, id)
	if err != nil {
		return user, err
	}

	return user, err
}

func (s *UserService) CreateUser(ctx context.Context, name string) (int, error) {
	return s.userStorage.CreateUser(ctx, name)
}
