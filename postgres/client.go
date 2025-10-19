package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DB *sqlx.DB

func NewClient(_ context.Context, dbURL string) (DB, error) {
	return sqlx.Open("postgres", dbURL)
}
