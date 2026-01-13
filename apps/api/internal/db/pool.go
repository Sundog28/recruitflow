package db

import (
"context"
"fmt"
"time"

"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(dbURL string) (*pgxpool.Pool, error) {
cfg, err := pgxpool.ParseConfig(dbURL)
if err != nil {
return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
}

cfg.MaxConns = 10
cfg.MinConns = 2
cfg.MaxConnLifetime = 30 * time.Minute
cfg.MaxConnIdleTime = 5 * time.Minute

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

pool, err := pgxpool.NewWithConfig(ctx, cfg)
if err != nil {
return nil, fmt.Errorf("connect db: %w", err)
}

if err := pool.Ping(ctx); err != nil {
pool.Close()
return nil, fmt.Errorf("ping db: %w", err)
}

return pool, nil
}
