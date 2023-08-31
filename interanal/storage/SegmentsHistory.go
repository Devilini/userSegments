package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"userSegments/interanal/model"
)

type SegmentHistory interface {
	GetSegmentsHistory(ctx context.Context, date string) ([]model.SegmentHistoryReport, error)
}

type SegmentHistoryStorage struct {
	client *pgxpool.Pool
}

func NewSegmentHistoryStorage(client *pgxpool.Pool) SegmentHistoryStorage {
	return SegmentHistoryStorage{client: client}
}

func (s *SegmentHistoryStorage) GetSegmentsHistory(ctx context.Context, date string) ([]model.SegmentHistoryReport, error) {
	query := `SELECT segments_history.id, user_id, segments.slug as segment, operation, created_at FROM segments_history inner join segments on segments.id = segment_history.segment_id WHERE created_at >= $1`
	var segments []model.SegmentHistoryReport
	rows, err := s.client.Query(ctx, query, date)
	if err != nil {
		return segments, fmt.Errorf("unable to query user: %w", err)
	}

	for rows.Next() {
		segment := model.SegmentHistoryReport{}
		err := rows.Scan(&segment.Id, &segment.UserId, &segment.Segment, &segment.Operation, &segment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		segments = append(segments, segment)
	}
	log.Println(segments)

	return segments, nil
}
