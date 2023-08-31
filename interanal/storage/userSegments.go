package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"userSegments/interanal/model"
)

type UserSegments interface {
	GetUserSegments(ctx context.Context, id int) ([]model.Segment, error)
	GetSegmentsBySlugs(ctx context.Context, slugs []string, userId int) ([]model.Segment, error)
	AddUserToSegment(ctx context.Context, userSegments []model.UserSegments, segmentsToDel []int, userId int) (int, error)
	DeleteSlugForUsers(ctx context.Context, segmentId int) error
}

type UserSegmentsStorage struct {
	client *pgxpool.Pool
}

func NewUserSegmentsStorage(client *pgxpool.Pool) UserSegmentsStorage {
	return UserSegmentsStorage{client: client}
}

func (s *UserSegmentsStorage) AddUserToSegment(
	ctx context.Context,
	userSegments []model.UserSegments,
	segmentsToDel []int,
	userId int) (int, error) {
	tx, err := s.client.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		log.Fatalf("Unable to begin transaction because %s", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println(err)
		}
	}(tx, ctx)

	for _, user := range userSegments {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO user_segments(user_id, segment_id) VALUES($1, $2)`,
			user.UserId, user.SegmentId,
		)
		if err != nil {
			log.Fatalf("err %s", err)
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO segments_history(user_id, segment_id, operation) VALUES($1, $2, $3)`,
			user.UserId, user.SegmentId, "add",
		)
		if err != nil {
			log.Fatalf("err %s", err)
		}
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM user_segments WHERE user_id = $1 AND segment_id = any($2)`,
		userId,
		segmentsToDel,
	)
	if err != nil {
		log.Fatalf("Error deleting %s", err)
	}

	for _, seg := range segmentsToDel {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO segments_history(user_id, segment_id, operation) VALUES($1, $2, $3)`,
			userId, seg, "delete",
		)
		if err != nil {
			log.Fatalf("err %s", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return 1, nil
}

func (s *UserSegmentsStorage) GetUserSegments(ctx context.Context, userId int) ([]model.Segment, error) {
	query := `SELECT id, slug FROM segments inner join user_segments on user_segments.segment_id = segments.id WHERE user_id=$1`
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

func (s *UserSegmentsStorage) GetSegmentsBySlugs(ctx context.Context, slugs []string, userId int) ([]model.Segment, error) {
	query := `SELECT id, slug FROM segments inner join user_segments on user_segments.segment_id = segments.id WHERE user_id = $1 AND segments.slug=any($2)`
	var segments []model.Segment
	rows, err := s.client.Query(ctx, query, userId, slugs)
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

func (s *UserSegmentsStorage) DeleteSlugForUsers(ctx context.Context, segmentId int) error {
	query := "DELETE from user_segments WHERE segment_id=$1 RETURNING segment_id"
	_, err := s.client.Exec(ctx, query, segmentId)
	if err != nil {
		return err
	}

	return nil
}
