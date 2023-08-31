package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"userSegments/interanal/model"
)

type UserSegments interface {
	GetUserSegments(ctx context.Context, id int) ([]model.Segment, error)
	GetSegmentsBySlugs(ctx context.Context, slugs []string, userId int) ([]model.Segment, error)
	ChangeUserSegments(ctx context.Context, userSegments []model.UserSegments, segmentsToDel []int, userId int) (int, error)
	DeleteSlugForUsers(ctx context.Context, segmentId int) error
}

type UserSegmentsStorage struct {
	client *pgxpool.Pool
}

func NewUserSegmentsStorage(client *pgxpool.Pool) UserSegmentsStorage {
	return UserSegmentsStorage{client: client}
}

func (s *UserSegmentsStorage) ChangeUserSegments(
	ctx context.Context,
	userSegments []model.UserSegments,
	segmentsToDel []int,
	userId int) (int, error) {
	tx, err := s.client.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return 0, fmt.Errorf("unable to begin transaction because %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			logrus.Error(err)
		}
	}(tx, ctx)

	for _, user := range userSegments {
		_, err = tx.Exec(
			ctx,
			fmt.Sprintf("INSERT INTO %s(user_id, segment_id) VALUES($1, $2)", userSegmentsTable),
			user.UserId, user.SegmentId,
		)
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(
			ctx,
			fmt.Sprintf("INSERT INTO %s(user_id, segment_id, operation) VALUES($1, $2, 'add')", segmentsHistoryTable),
			user.UserId, user.SegmentId,
		)
		if err != nil {
			return 0, err
		}
	}

	_, err = tx.Exec(
		ctx,
		fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND segment_id = any($2)", userSegmentsTable),
		userId,
		segmentsToDel,
	)
	if err != nil {
		return 0, err
	}

	for _, seg := range segmentsToDel {
		_, err = tx.Exec(
			ctx,
			fmt.Sprintf("INSERT INTO %s(user_id, segment_id, operation) VALUES($1, $2, 'delete')", segmentsHistoryTable),
			userId, seg,
		)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (s *UserSegmentsStorage) GetUserSegments(ctx context.Context, userId int) ([]model.Segment, error) {
	query := fmt.Sprintf("SELECT id, slug FROM %s inner join %s on user_segments.segment_id = segments.id WHERE user_id=$1", segmentsTable, userSegmentsTable)
	var segments []model.Segment
	rows, err := s.client.Query(ctx, query, userId)
	if err != nil {
		return segments, fmt.Errorf("unable to query: %w", err)
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

func (s *UserSegmentsStorage) GetSegmentsBySlugs(ctx context.Context, slugs []string, userId int) ([]model.Segment, error) {
	query := fmt.Sprintf("SELECT id, slug FROM segments inner join %s on user_segments.segment_id = segments.id WHERE user_id = $1 AND segments.slug=any($2)", userSegmentsTable)
	var segments []model.Segment
	rows, err := s.client.Query(ctx, query, userId, slugs)
	if err != nil {
		return segments, fmt.Errorf("unable to query: %w", err)
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

func (s *UserSegmentsStorage) DeleteSlugForUsers(ctx context.Context, segmentId int) error {
	query := fmt.Sprintf("DELETE from %s WHERE segment_id=$1 RETURNING segment_id", userSegmentsTable)
	_, err := s.client.Exec(ctx, query, segmentId)

	return err
}
