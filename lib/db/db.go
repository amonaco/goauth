package db

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

// Connect initializes the PostgreSQL connection pool.
func Connect(ctx context.Context, connString string) error {
    cfg, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return err
    }

    cfg.MaxConns = 10
    cfg.MinConns = 2
    cfg.MaxConnLifetime = 30 * time.Minute
    cfg.HealthCheckPeriod = 1 * time.Minute

    dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil {
        return err
    }

    if err := dbPool.Ping(ctx); err != nil {
        dbPool.Close()
        return err
    }

    pool = dbPool
    return nil
}

// Close closes the database pool.
func Close() {
    if pool != nil {
        pool.Close()
        pool = nil
    }
}

// DB returns the global pool.
func DB() *pgxpool.Pool {
    if pool == nil {
        panic("database pool is not initialized")
    }
    return pool
}

// QueryRow executes a query and returns a row.
func QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
    return DB().QueryRow(ctx, query, args...)
}

// Exec executes a query without returning rows.
func Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
    return DB().Exec(ctx, query, args...)
}

func init() {
    // noop to keep package importable for the DB layer
}

func example() {
    _ = fmt.Sprintf("placeholder")
}
