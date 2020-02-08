package dblib

import (
	"database/sql"
)

// Wrappers so we can add actions as methods to a transaction instead of writing them inline
type DB struct {
	*sql.DB
}
type Tx struct {
	*sql.Tx
}

// Open returns a DB reference for a data source.
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Ping checks if the datasource is valid
func (db *DB) Ping() error {
	err := db.DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

// Prepare prepares a query statement but does not execute
func (db *DB) Prepare(statement string) (*sql.Stmt, error) {
	queryStmt, err := db.DB.Prepare(statement)
	if err != nil {
		return nil, err
	}
	return queryStmt, nil
}

// Begin starts an returns a new transaction.
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
