package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"userSegments/interanal/domain/user/model"
)

type User interface {
	//CreateProduct(ctx context.Context, name string) (int, error)
	GetById(ctx context.Context, id int) (model.User, error)
	Create(ctx context.Context, name string) (int, error)
	//GetAllProducts(ctx context.Context) ([]entity.Product, error)
}

type UserStorage struct {
	client PostgresClient
}

type PostgresClient interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	Stat() *pgxpool.Stat
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewUserStorage(client PostgresClient) UserStorage {
	return UserStorage{client: client}
}

//func (repo *UserStorage) All(ctx context.Context) ([]model.User, error) {
//	all, err := repo.findBy(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	resp := make([]model.User, len(all))
//	for i, e := range all {
//		resp[i] = e.ToDomain()
//	}
//
//	return resp, nil
//}
//
//func (r *UserStorage) GetById(ctx context.Context, id int) (entity.User, error) {
//	sql, args, _ := r.Builder.
//		Select("id, username, password, created_at").
//		From("users").
//		Where("id = ?", id).
//		ToSql()
//
//	var user entity.User
//	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
//		&user.Id,
//		&user.Username,
//		&user.Password,
//		&user.CreatedAt,
//	)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return entity.User{}, repoerrs.ErrNotFound
//		}
//		return entity.User{}, fmt.Errorf("UserRepo.GetUserById - r.Pool.QueryRow: %v", err)
//	}
//
//	return user, nil
//}

func (s *UserStorage) GetById(ctx context.Context, id int) (model.User, error) {
	query := `SELECT id, name FROM users WHERE id=$1`
	//
	//rows, err := s.client.Query(ctx, query, id)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to query users: %w", err)
	//}
	//defer rows.Close()

	//return pgx.CollectRows(rows, pgx.RowToStructByName(model.User))

	var user model.User
	err := s.client.QueryRow(ctx, query, id).Scan(&user.Id, &user.Name)
	if err != nil {
		return user, fmt.Errorf("unable to query user: %w", err)
	}

	return user, nil
}

func (s *UserStorage) Create(ctx context.Context, name string) (int, error) {
	var id int
	query := "INSERT INTO users (name) values ($1) RETURNING id" // todo table name
	row := s.client.QueryRow(ctx, query, name)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
