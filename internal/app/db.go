package app

import "database/sql"

type DB interface {
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Begin() (*sql.Tx, error)
	Query(query string, args ...any) (*sql.Rows, error)
}
