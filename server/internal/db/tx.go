package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type (
	PgTxManager struct {
		db      *pgx.Conn
		queries *Queries
	}
)

func NewTxManager(db *pgx.Conn, queries *Queries) *PgTxManager {
	return &PgTxManager{db: db, queries: queries}
}

func (m *PgTxManager) DoTx(ctx context.Context, cb func(qtx *Queries) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = cb(m.queries.WithTx(tx))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
