package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/kokoc-hack/internal/config"
)

type Postgres struct {
	pool    *pgxpool.Pool
	timeout int
}

func NewPostgres(timeout int) (*Postgres, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", config.Cfg.PostgresUser, config.Cfg.PostgresPassword, config.Cfg.PostgresHost, config.Cfg.PostgresDb)

	for attempts := 0; attempts < 5; attempts++ {
		time.Sleep(1)
		pool, err := pgxpool.New(context.Background(), dsn)

		if err == nil {
			return &Postgres{pool: pool, timeout: timeout}, nil
		}
	}

	return nil, errors.New("can't connect to postgres")
}

func (pg *Postgres) GetSessionWithContext(c context.Context) (*pgxpool.Conn, error) {
	return pg.pool.Acquire(c)
}

func (pg *Postgres) Close() {
	pg.pool.Close()
}

func (pg *Postgres) Timeout() int {
	return pg.timeout
}

func (pg *Postgres) GetPool() *pgxpool.Pool {
	return pg.pool
}
