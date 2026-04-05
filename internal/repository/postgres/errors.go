package postgres

import (
	"errors"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	PgUniqueViolation     = "23505"
	PgForeignKeyViolation = "23503"
)

func ParseError(err error, operation string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return entity.ErrNotFound.WithOperation(operation)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PgUniqueViolation:
			return entity.ErrConflict.
				WithInternal(err).
				WithOperation(operation)
		case PgForeignKeyViolation:
			return entity.ErrNotFound.
				WithInternal(err).
				WithOperation(operation)
		}
	}

	return entity.NewAppError(entity.CodeInternal, "database operation failed").
		WithInternal(err).
		WithOperation(operation)
}
