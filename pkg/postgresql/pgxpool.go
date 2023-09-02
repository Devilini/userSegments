package psql

import (
	"context"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgConfig struct {
	username string
	password string
	host     string
	port     string
	database string
}

func NewPgConfig(username string, password string, host string, port string, database string) *pgConfig {
	return &pgConfig{
		username: username,
		password: password,
		host:     host,
		port:     port,
		database: database,
	}
}

func NewClient(
	ctx context.Context,
	maxAttempts int,
	maxDelay time.Duration,
	cfg *pgConfig,
) (pool *pgxpool.Pool, err error) {
	pgxCfg, parseConfigErr := pgxpool.ParseConfig(cfg.ConnStringFromCfg())
	if parseConfigErr != nil {
		logrus.Info("Unable to parse config: %v\n", parseConfigErr)
		return nil, parseConfigErr
	}

	pool, parseConfigErr = pgxpool.NewWithConfig(ctx, pgxCfg)
	if parseConfigErr != nil {
		logrus.Info("Failed to parse PostgreSQL configuration due to error: %v\n", parseConfigErr)
		return nil, parseConfigErr
	}

	err = DoWithAttempts(func() error {
		pingErr := pool.Ping(ctx)
		if pingErr != nil {
			logrus.Info("Failed to connect to postgres due to error %v... Going to do the next attempt\n", pingErr)
			return pingErr
		}

		return nil
	}, maxAttempts, maxDelay)
	if err != nil {
		logrus.Info("All attempts are exceeded. Unable to connect to PostgreSQL")
		return pool, err
	}

	return pool, nil
}

func DoWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--

			continue
		}

		return nil
	}

	return err
}

func (c *pgConfig) ConnStringFromCfg() string {
	url := strings.Builder{}
	url.WriteString("postgresql://")
	url.WriteString(c.username)
	url.WriteString(":")
	url.WriteString(c.password)
	url.WriteString("@")
	url.WriteString(c.database)
	url.WriteString(":")
	url.WriteString(c.port)
	url.WriteString("/")
	url.WriteString(c.database)

	return url.String()
}
