package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps pgx pool for application repositories.
type DB struct {
    Pool *pgxpool.Pool
}

// NewDB creates a new pgx connection pool with sane defaults.
func NewDB(ctx context.Context, dsn string) (*DB, error) {
    cfg, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, err
    }
    // 커넥션 풀 설정 기본값 조정
    cfg.MaxConns = 10
    cfg.MinConns = 1
    cfg.MaxConnLifetime = time.Hour
    cfg.MaxConnIdleTime = 30 * time.Minute
    cfg.HealthCheckPeriod = 30 * time.Second

    pool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil {
        return nil, err
    }
    // 연결 확인
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, err
    }
    return &DB{Pool: pool}, nil
}

// Close closes the underlying pool.
func (db *DB) Close() {
    if db != nil && db.Pool != nil {
        db.Pool.Close()
    }
}


