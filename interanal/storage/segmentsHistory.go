package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"userSegments/interanal/model"
)

type SegmentHistory interface {
	GetSegmentsHistory(ctx context.Context, dateFrom string, dateTo string) ([]model.SegmentHistoryReport, error)
}

type SegmentHistoryStorage struct {
	client *pgxpool.Pool
}

func NewSegmentHistoryStorage(client *pgxpool.Pool) SegmentHistoryStorage {
	return SegmentHistoryStorage{client: client}
}

func (s *SegmentHistoryStorage) GetSegmentsHistory(ctx context.Context, dateFrom string, dateTo string) ([]model.SegmentHistoryReport, error) {
	query := fmt.Sprintf("SELECT segments_history.id, user_id, segments.slug as segment, operation, created_at "+
		"FROM %s inner join %s on segments.id = segments_history.segment_id WHERE created_at BETWEEN $1 AND $2",
		segmentsHistoryTable,
		segmentsTable,
	)
	var segments []model.SegmentHistoryReport
	rows, err := s.client.Query(ctx, query, dateFrom, dateTo)
	if err != nil {
		return segments, fmt.Errorf("unable to query: %w", err)
	}

	for rows.Next() {
		segment := model.SegmentHistoryReport{}
		err := rows.Scan(&segment.Id, &segment.UserId, &segment.Segment, &segment.Operation, &segment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		segments = append(segments, segment)
	}

	return segments, nil
}
