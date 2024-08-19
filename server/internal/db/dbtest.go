package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"

	_ "github.com/mattn/go-sqlite3"
)

// type sqlitePgAdapter struct {
// 	db *sql.DB
// 	tx *sql.Tx
// }
//
// // Conn returns the underlying *Conn that on which this transaction is executing.
// func (s *sqlitePgAdapter) Conn() *pgx.Conn {
// 	return &pgx.Conn{}
// }
//
// type rowsAdapter struct {
// 	rows *sql.Rows
// }
//
// func (ra rowsAdapter) Close()                        { ra.rows.Close() }
// func (ra rowsAdapter) Err() error                    { return ra.rows.Err() }
// func (ra rowsAdapter) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }
// func (ra rowsAdapter) FieldDescriptions() []pgconn.FieldDescription {
// 	return []pgconn.FieldDescription{}
// }
// func (ra rowsAdapter) Next() bool             { return ra.rows.Next() }
// func (ra rowsAdapter) Scan(dest ...any) error { return ra.rows.Scan(dest...) }
// func (ra rowsAdapter) Values() ([]any, error) { return []any{}, nil }
// func (ra rowsAdapter) RawValues() [][]byte    { return [][]byte{} }
// func (ra rowsAdapter) Conn() *pgx.Conn        { return nil }
//
// type rowAdapter struct {
// 	rows *sql.Row
// }
//
// func (ra rowAdapter) Scan(dest ...any) error { return ra.rows.Scan(dest...) }
//
// func (s *sqlitePgAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
// 	_, err := s.db.ExecContext(ctx, sql, args...)
// 	if err != nil {
// 		return pgconn.CommandTag{}, err
// 	}
// 	return pgconn.CommandTag{}, nil
// }
//
// func (s *sqlitePgAdapter) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
// 	res, err := s.db.QueryContext(ctx, sql, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ra := rowsAdapter{rows: res}
//
// 	return ra, nil
// }
//
// func (s *sqlitePgAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
// 	res := s.db.QueryRow(sql, args...)
// 	ra := rowAdapter{rows: res}
//
// 	return ra
// }
//
// func (s *sqlitePgAdapter) Begin(ctx context.Context) (pgx.Tx, error) {
// 	tx, err := s.db.BeginTx(ctx, nil)
// 	s.tx = tx
// 	return s, err
// }
//
// func (s *sqlitePgAdapter) Rollback(ctx context.Context) error {
// 	if s.tx == nil {
// 		return fmt.Errorf("no transaction to rollback")
// 	}
// 	return s.tx.Rollback()
// }
//
// func (s *sqlitePgAdapter) Commit(ctx context.Context) error {
// 	if s.tx == nil {
// 		return fmt.Errorf("no transaction to commit")
// 	}
// 	return s.tx.Commit()
// }

func NewTestDB(t *testing.T) (*Queries, *PgTxManager) {
	t.Helper()

	ctx := context.Background()

	pgdb, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("opening ent client: %v", err)
	}

	// clean database
	if _, err := pgdb.Exec(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		t.Fatal("failed to clean schema", err)
	}

	schema, err := os.ReadFile("../../sql/schema.sql")
	if err != nil {
		t.Fatal("failed to open schema file", err)
	}

	// create tables
	if _, err := pgdb.Exec(ctx, string(schema)); err != nil {
		t.Fatal("failed to create schema", err)
	}

	db := New(pgdb)

	return db, NewTxManager(pgdb, db)
}
