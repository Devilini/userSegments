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
	AddUserToSegment(ctx context.Context, userSegments []model.UserSegments, segmentsToDel []string, userId int) (int, error)
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
	segmentsToDel []string,
	userId int) (int, error) {
	// Connecting to database...
	//rows := [][]interface{}{
	//	{"John", "Smith", int32(36)},
	//	{"Jane", "Doe", int32(29)},
	//}
	//copyCount, err := s.client.CopyFrom(
	//	ctx,
	//	pgx.Identifier{"people"},
	//	[]string{"first_name", "last_name", "age"},
	//	pgx.CopyFromRows(rows),
	//)

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
	}

	// TODO batch insert
	//query := `INSERT INTO user_segments(user_id, segment_id) VALUES($1, $2)`
	//batch := &pgx.Batch{}
	//for _, user := range userSegments {
	//	batch.Queue(query, user.UserId, user.SegmentId)
	//}
	//
	//results := s.client.SendBatch(ctx, batch)
	////_, err = results.Exec()
	////if err != nil {
	////	return 0, err
	////}
	//defer func(results pgx.BatchResults) {
	//	err := results.Close()
	//	if err != nil {
	//
	//	}
	//}(results)

	//rows := [][]interface{}{
	//	{int32(5), int32(5)},
	//}
	//_, err = tx.CopyFrom(
	//	ctx,
	//	pgx.Identifier{"user_segments"},
	//	[]string{"user_id", "segment_id"},
	//	pgx.CopyFromRows(rows),
	//)

	//query := `INSERT INTO user_segments(user_id, segment_id) VALUES($1, $2)`
	//batch := &pgx.Batch{}
	//for _, user := range userSegments {
	//	//args := pgx.NamedArgs{
	//	//	"userName":  user.Name,
	//	//	"userEmail": user.Email,
	//	//}
	//	batch.Queue(query, user.UserId, user.SegmentId)
	//}
	//
	//tx.SendBatch(ctx, batch)
	//_, err = results.Exec()
	//if err != nil {
	//	return 0, err
	//}

	//defer results.Close()

	_, err = tx.Exec(
		ctx,
		`DELETE FROM user_segments WHERE user_id = $1 AND segment_id in(Select id from segments WHERE slug = any($2))`,
		userId,
		segmentsToDel,
	)
	if err != nil {
		log.Fatalf("Error deleting %s", err)
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
