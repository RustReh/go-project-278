package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// Ошибки репозитория. Service маппит их в apperr.
var (
	ErrNotFound     = errors.New("link not found")
	ErrConflict     = errors.New("link already exists")
	ErrInvalidInput = errors.New("invalid input")
)

// ErrInternal — неожиданная ошибка БД; service отдаёт apperr.CodeInternal.
var ErrInternal = errors.New("repository internal error")

// MapError переводит ошибки sqlc/SQL/Postgres в ошибки репозитория.
func MapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}
	return fmt.Errorf("%w: %w", ErrInternal, err)
}
