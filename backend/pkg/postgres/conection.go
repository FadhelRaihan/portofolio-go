package postgres

import (
    "context"
    "fmt"
    "net/url"

    "github.com/jackc/pgx/v5/pgxpool"
    "backend/internal/config"
)

func NewPool(env config.Env) (*pgxpool.Pool, error) {
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        env.DBUser,
        url.QueryEscape(env.DBPass),
        env.DBHost,
        env.DBPort,
        env.DBName,
    )

    pool, err := pgxpool.New(context.Background(), dsn)
    if err != nil {
        return nil, err
    }
    if err := pool.Ping(context.Background()); err != nil {
        pool.Close()
        return nil, err
    }
    return pool, nil
}
