package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"userSegments/interanal/apperror"
	"userSegments/interanal/model"
)

const (
	OperationTypeAdd    string = "add"
	OperationTypeDelete string = "delete"
)

type SegmentHistory interface {
	GetSegmentsHistory(ctx context.Context, dateFrom string, dateTo string) ([]model.SegmentHistoryReport, error)
	BulkInsert(ctx context.Context, segmentHistory []model.SegmentHistory) (int64, error)
}

type SegmentHistoryStorage struct {
	client *pgxpool.Pool
}

func NewSegmentHistoryStorage(client *pgxpool.Pool) SegmentHistoryStorage {
	return SegmentHistoryStorage{client: client}
}

func (s *SegmentHistoryStorage) GetSegmentsHistory(ctx context.Context, dateFrom string, dateTo string) ([]model.SegmentHistoryReport, error) {
	query := fmt.Sprintf("SELECT sh.id, sh.user_id, s.slug, sh.operation, sh.created_at "+
		"FROM %s sh inner join %s s on s.id = sh.segment_id WHERE sh.created_at BETWEEN $1 AND $2 order by sh.user_id",
		segmentsHistoryTable,
		segmentsTable,
	)
	var segments []model.SegmentHistoryReport
	rows, err := s.client.Query(ctx, query, dateFrom, dateTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return segments, apperror.NotFoundError("not found Segment History data")
		}
		return segments, err
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

func (s *SegmentHistoryStorage) BulkInsert(ctx context.Context, segmentHistory []model.SegmentHistory) (int64, error) {
	var entries [][]any
	columns := []string{"user_id", "segment_id", "operation"}
	tableName := segmentsHistoryTable
	for _, item := range segmentHistory {
		entries = append(entries, []any{item.UserId, item.SegmentId, item.Operation})
	}
	countRows, err := s.client.CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		columns,
		pgx.CopyFromRows(entries),
	)
	if err != nil {
		return 0, fmt.Errorf("error copying into %s table: %w", tableName, err)
	}

	return countRows, nil
}
