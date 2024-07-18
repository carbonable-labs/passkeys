package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/mattn/go-sqlite3"
)

type PgCloserFn = func() error

type sqlitePgAdapter struct {
	db *sql.DB
}

type rowsAdapter struct {
	rows *sql.Rows
}

func (ra rowsAdapter) Close()                        { ra.rows.Close() }
func (ra rowsAdapter) Err() error                    { return ra.rows.Err() }
func (ra rowsAdapter) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }
func (ra rowsAdapter) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{}
}
func (ra rowsAdapter) Next() bool             { return ra.rows.Next() }
func (ra rowsAdapter) Scan(dest ...any) error { return ra.rows.Scan(dest...) }
func (ra rowsAdapter) Values() ([]any, error) { return []any{}, nil }
func (ra rowsAdapter) RawValues() [][]byte    { return [][]byte{} }
func (ra rowsAdapter) Conn() *pgx.Conn        { return nil }

type rowAdapter struct {
	rows *sql.Row
}

func (ra rowAdapter) Scan(dest ...any) error { return ra.rows.Scan(dest...) }

func (s *sqlitePgAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	_, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	return pgconn.CommandTag{}, nil
}

func (s *sqlitePgAdapter) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	res, err := s.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	ra := rowsAdapter{rows: res}

	return ra, nil
}

func (s *sqlitePgAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	res := s.db.QueryRow(sql, args...)
	ra := rowAdapter{rows: res}

	return ra
}

func NewTestDB(t *testing.T) *Queries {
	t.Helper()

	ctx := context.Background()

	drv, err := sql.Open(
		"sqlite3",
		"file:db?mode=memory&_fk=1&_journal_mode=WAL",
	)
	if err != nil {
		t.Fatalf("opening ent client: %v", err)
	}
	pgAdapter := &sqlitePgAdapter{db: drv}

	// pgdb, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	t.Fatal("failed to acquire postgres connection", err)
	// }

	schema, err := os.ReadFile("../../sql/schema.sql")
	if err != nil {
		t.Fatal("failed to open schema file", err)
	}

	// create tables
	if _, err := pgAdapter.Exec(ctx, string(schema)); err != nil {
		t.Fatal("failed to create schema", err)
	}

	db := New(pgAdapter)

	return db
}
