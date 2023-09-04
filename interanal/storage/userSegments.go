package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
	"userSegments/interanal/apperror"
	"userSegments/interanal/model"
)

type UserSegments interface {
	GetUserSegments(ctx context.Context, id int) ([]model.Segment, error)
	ChangeUserSegments(ctx context.Context, userSegments []model.UserSegments, segmentsToDel []int, userId int) (int, error)
	DeleteSlugForUsers(ctx context.Context, segmentId int) ([]int, error)
	DeleteAllExpired(ctx context.Context) ([]model.UserSegments, error)
	CreateByPercent(ctx context.Context, percent int, segmentId int) error
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
		res, err := tx.Exec(
			ctx,
			fmt.Sprintf("INSERT INTO %s(user_id, segment_id, expired_date) VALUES($1, $2, $3) ON CONFLICT ON CONSTRAINT user_segments_pkey DO NOTHING", userSegmentsTable),
			user.UserId, user.SegmentId, (*time.Time)(user.ExpiredAt),
		)
		if err != nil {
			return 0, err
		}

		if res.RowsAffected() == 0 {
			continue
		}
		_, err = tx.Exec(
			ctx,
			fmt.Sprintf("INSERT INTO %s(user_id, segment_id, operation) VALUES($1, $2, $3)", segmentsHistoryTable),
			user.UserId, user.SegmentId, OperationTypeAdd,
		)
		if err != nil {
			return 0, err
		}
	}

	res, err := tx.Exec(
		ctx,
		fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND segment_id = any($2)", userSegmentsTable),
		userId,
		segmentsToDel,
	)
	if err != nil {
		return 0, err
	}

	if res.RowsAffected() > 0 {
		for _, seg := range segmentsToDel {
			_, err = tx.Exec(
				ctx,
				fmt.Sprintf("INSERT INTO %s(user_id, segment_id, operation) VALUES($1, $2, $3)", segmentsHistoryTable),
				userId, seg, OperationTypeDelete,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (s *UserSegmentsStorage) GetUserSegments(ctx context.Context, userId int) ([]model.Segment, error) {
	query := fmt.Sprintf(
		"SELECT s.id, s.slug FROM %s s inner join %s us on us.segment_id = s.id WHERE us.user_id=$1 AND (us.expired_date is null or us.expired_date > $2)",
		segmentsTable,
		userSegmentsTable,
	)
	var segments []model.Segment
	rows, err := s.client.Query(ctx, query, userId, time.Now().Format(time.DateTime))
	if err != nil {
		if err == pgx.ErrNoRows {
			return segments, apperror.NotFoundError("not found user segments")
		}
		return segments, err
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

func (s *UserSegmentsStorage) DeleteSlugForUsers(ctx context.Context, segmentId int) ([]int, error) {
	var userIds []int
	query := fmt.Sprintf("DELETE from %s WHERE segment_id=$1 RETURNING user_id", userSegmentsTable)
	rows, err := s.client.Query(ctx, query, segmentId)
	if err != nil {
		return userIds, fmt.Errorf("unable to scan row: %w", err)
	}

	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			return userIds, fmt.Errorf("unable to scan row: %w", err)
		}
		userIds = append(userIds, userId)
	}

	return userIds, err
}

func (s *UserSegmentsStorage) DeleteAllExpired(ctx context.Context) ([]model.UserSegments, error) {
	var deletedSegments []model.UserSegments
	query := fmt.Sprintf("DELETE from %s WHERE expired_date <= $1 RETURNING segment_id, user_id", userSegmentsTable)
	rows, err := s.client.Query(ctx, query, time.Now().Format(time.DateTime))
	if err != nil {
		return deletedSegments, fmt.Errorf("unable to scan row: %w", err)
	}

	for rows.Next() {
		segment := model.UserSegments{}
		err := rows.Scan(&segment.UserId, &segment.SegmentId)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		deletedSegments = append(deletedSegments, segment)
	}

	return deletedSegments, err
}

func (s *UserSegmentsStorage) CreateByPercent(ctx context.Context, percent int, segmentId int) error {
	rows, err := s.client.Query(
		ctx,
		fmt.Sprintf("SELECT id FROM %s order by RANDOM() LIMIT (SELECT count(*)::float FROM %s) / 100 * $1", usersTable, usersTable),
		percent,
	)
	if err != nil {
		return err
	}

	var users []int
	for rows.Next() {
		var user int
		err := rows.Scan(&user)
		if err != nil {
			return fmt.Errorf("unable to scan row: %w", err)
		}
		users = append(users, user)
	}

	tx, err := s.client.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("unable to begin transaction because %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			logrus.Error(err)
		}
	}(tx, ctx)

	for _, userId := range users {
		_, err = tx.Exec(
			ctx,
			fmt.Sprintf("INSERT INTO %s(user_id, segment_id) VALUES($1, $2) ON CONFLICT ON CONSTRAINT user_segments_pkey DO NOTHING", userSegmentsTable),
			userId, segmentId,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(
			ctx,
			fmt.Sprintf("INSERT INTO %s(user_id, segment_id, operation) VALUES($1, $2, 'add')", segmentsHistoryTable),
			userId, segmentId,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
