package storage

import (
	"context"
	"fmt"
	"userSegments/interanal/model"
)

type UserSegments interface {
	GetUserSegments(ctx context.Context, id int) ([]model.Segment, error)
	//CreateUser(ctx context.Context, name string) (int, error)
}

type UserSegmentsStorage struct {
	client PostgresClient
}

func NewUserSegmentsStorage(client PostgresClient) UserSegmentsStorage {
	return UserSegmentsStorage{client: client}
}

func (s *UserSegmentsStorage) GetUserSegments(ctx context.Context, userId int) ([]model.Segment, error) {
	query := `SELECT id, slug FROM segments inner join user_segments on user_segments.segment_id = segments.id WHERE user_id=$1`
	//query := `SELECT id, slug FROM segments`
	var segments []model.Segment
	rows, err := s.client.Query(ctx, query, userId)
	if err != nil {
		return segments, fmt.Errorf("unable to query user: %w", err)
	}

	segments = []model.Segment{}
	for rows.Next() {
		segment := model.Segment{}
		err := rows.Scan(&segment.Id, &segment.Slug)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

//func (s *UserStorage) CreateUser(ctx context.Context, name string) (int, error) {
//	var id int
//	query := "INSERT INTO users (name) values ($1) RETURNING id" // todo table name
//	row := s.client.QueryRow(ctx, query, name)
//	if err := row.Scan(&id); err != nil {
//		return 0, err
//	}
//
//	return id, nil
//}
