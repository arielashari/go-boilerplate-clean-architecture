package postgres

import (
	"context"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) entity.Transactor {
	return &transactor{pool: pool}
}

func (t *transactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(injectTx(ctx, tx)); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// context key for injecting tx
type txKey struct{}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txKey{}).(pgx.Tx)
	return tx
}
