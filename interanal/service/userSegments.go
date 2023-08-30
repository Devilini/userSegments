package service

import (
	"context"
	"userSegments/interanal/model"
	"userSegments/interanal/storage"
)

//type repository interface {
//	//All(ctx context.Context) ([]model.User, error)
//	//Create(ctx context.Context, req model.Crea) error
//}

type UserSegments interface {
	GetUserSegments(ctx context.Context, id int) ([]model.Segment, error)
	//CreateUserSegments(ctx context.Context, name string) (int, error)
}

type UserSegmentsService struct {
	userSegmentsStorage storage.UserSegments
}

func NewUserSegmentsService(storage storage.UserSegments) *UserSegmentsService {
	return &UserSegmentsService{userSegmentsStorage: storage}
}

func (s *UserSegmentsService) GetUserSegments(ctx context.Context, id int) ([]model.Segment, error) {
	return s.userSegmentsStorage.GetUserSegments(ctx, id)
}
