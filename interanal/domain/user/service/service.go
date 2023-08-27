package service

import (
	"context"
	"userSegments/interanal/domain/user/storage"

	"userSegments/interanal/domain/user/model"
)

//type repository interface {
//	//All(ctx context.Context) ([]model.User, error)
//	//Create(ctx context.Context, req model.Crea) error
//}

type User interface {
	//CreateProduct(ctx context.Context, name string) (int, error)
	GetUserById(ctx context.Context, id int) (model.User, error)
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
	return s.userStorage.GetById(ctx, id)
	//return model.User{
	//	1,
	//	"req.Name",
	//}, nil
}

//func (s *UserService) CreateProduct(ctx context.Context, req model.CreateProduct) (model.User, error) {
//	// cache
//
//	err := s.repository.Create(ctx, req)
//	if err != nil {
//		return model.User{}, err
//	}
//
//	return model.NewProduct(
//		req.ID,
//		req.Name,
//		req.Description,
//		req.ImageID,
//		nil,
//	), nil
//}
