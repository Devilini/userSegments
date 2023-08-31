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
	query := fmt.Sprintf("SELECT id, slug FROM %s WHERE id=$1", segmentsTable)
	var segment model.Segment
	err := s.client.QueryRow(ctx, query, id).Scan(&segment.Id, &segment.Slug)
	if err != nil {
		return segment, fmt.Errorf("unable to query segment: %w", err)
	}

	return segment, nil
}

func (s *SegmentStorage) GetSegmentBySlug(ctx context.Context, slug string) (model.Segment, error) {
	query := fmt.Sprintf("SELECT id, slug FROM %s WHERE slug=$1", segmentsTable)
	var segment model.Segment
	err := s.client.QueryRow(ctx, query, slug).Scan(&segment.Id, &segment.Slug)
	if err != nil {
		return segment, fmt.Errorf("unable to query segment: %w", err)
	}

	return segment, nil
}

func (s *SegmentStorage) CreateSegment(ctx context.Context, slug string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (slug) values ($1) RETURNING id", segmentsTable)
	row := s.client.QueryRow(ctx, query, slug)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SegmentStorage) DeleteSegmentBySlug(ctx context.Context, slug string) (int, error) {
	var id int
	query := fmt.Sprintf("DELETE from %s WHERE slug=$1 RETURNING id", segmentsTable)
	row := s.client.QueryRow(ctx, query, slug)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
