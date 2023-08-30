package service

import (
	"context"
	"fmt"
	"userSegments/interanal/model"
	"userSegments/interanal/storage"
)

//type repository interface {
//	//All(ctx context.Context) ([]model.User, error)
//	//Create(ctx context.Context, req model.Crea) error
//}

type User interface {
	GetUserById(ctx context.Context, id int) (model.User, error)
	CreateUser(ctx context.Context, name string) (int, error)
}

type UserService struct {
	//repository repository
	userStorage storage.User
}

func NewUserService(storage storage.User) *UserService {
	return &UserService{userStorage: storage}
}

//func (s *UserService) All(ctx context.Context) ([]model.User, error) {
//	products, err := s.repository.All(ctx)
//	if err != nil {
//		//return nil, errors.Wrap(err, "repository.All")
//	}
//
//	return products, nil
//}

func (s *UserService) GetUserById(ctx context.Context, id int) (model.User, error) {
	user, err := s.userStorage.GetUserById(ctx, id)
	if user.Id == 0 {
		return user, fmt.Errorf("user does not exists") // todo error
	}

	return user, err
}

func (s *UserService) CreateUser(ctx context.Context, name string) (int, error) {
	return s.userStorage.CreateUser(ctx, name)
}
