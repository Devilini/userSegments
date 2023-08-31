package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"userSegments/interanal/model"
)

type Segment interface {
	GetSegmentById(ctx context.Context, id int) (model.Segment, error)
	GetSegmentBySlug(ctx context.Context, slug string) (model.Segment, error)
	CreateSegment(ctx context.Context, slug string) (int, error)
	DeleteSegmentBySlug(ctx context.Context, slug string) (int, error)
}

type SegmentStorage struct {
	client *pgxpool.Pool
}

func NewSegmentStorage(client *pgxpool.Pool) SegmentStorage {
	return SegmentStorage{client: client}
}

func (s *SegmentStorage) GetSegmentById(ctx context.Context, id int) (model.Segment, error) {
	query := `SELECT id, slug FROM segments WHERE id=$1`

	var segment model.Segment
	err := s.client.QueryRow(ctx, query, id).Scan(&segment.Id, &segment.Slug)
	if err != nil {
		return segment, fmt.Errorf("unable to query segment: %w", err)
	}

	return segment, nil
}

func (s *SegmentStorage) GetSegmentBySlug(ctx context.Context, slug string) (model.Segment, error) {
	query := `SELECT id, slug FROM segments WHERE slug=$1`

	var segment model.Segment
	err := s.client.QueryRow(ctx, query, slug).Scan(&segment.Id, &segment.Slug)
	if err != nil {
		return segment, fmt.Errorf("unable to query segment: %w", err)
	}

	return segment, nil
}

func (s *SegmentStorage) CreateSegment(ctx context.Context, slug string) (int, error) {
	var id int
	query := "INSERT INTO segments (slug) values ($1) RETURNING id" // todo table name
	row := s.client.QueryRow(ctx, query, slug)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SegmentStorage) DeleteSegmentBySlug(ctx context.Context, slug string) (int, error) {
	var id int
	query := "DELETE from segments WHERE slug=$1 RETURNING id" // todo table name
	row := s.client.QueryRow(ctx, query, slug)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
