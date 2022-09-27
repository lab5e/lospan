package storage

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver for testing, local instances and in-memory database
)

// Storage holds all of the storage objects
type Storage struct {
	db       *sql.DB
	appStmt  applicationStatements
	devStmt  deviceStatements
	dataStmt dataStatements
	gwStmt   gatewayStatements
	keyStmt  keyStatements
}

// Close closes all of the storage instances.
func (s *Storage) Close() {
	s.appStmt.Close()
	s.devStmt.Close()
	s.dataStmt.Close()
	s.gwStmt.Close()
	s.keyStmt.Close()
}

// CreateStorage creates a new storage
func CreateStorage(connectionString string, maxConn, idleConn int, maxConnLifetime time.Duration) (*Storage, error) {
	return newStorage(driverName, connectionString, maxConn, idleConn, maxConnLifetime)
}

const driverName = "sqlite3"

func newStorage(driver, connectionString string, maxConn, idleConn int, maxConnLifetime time.Duration) (*Storage, error) {
	db, err := sql.Open(driver, connectionString)
	if nil != err {
		log.Fatalf("Unable to connect to database: %s", err)
		return nil, err
	}
	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)
	db.SetConnMaxLifetime(maxConnLifetime)

	if err := createSchema(db); err != nil {
		return nil, err
	}
	ret := &Storage{db: db}
	if err := ret.appStmt.prepare(db); err != nil {
		return nil, err
	}
	if err := ret.devStmt.prepare(db); err != nil {
		return nil, err
	}
	if err := ret.dataStmt.prepare(db); err != nil {
		return nil, err
	}
	if err := ret.gwStmt.prepare(db); err != nil {
		return nil, err
	}
	if err := ret.keyStmt.prepare(db); err != nil {
		return nil, err
	}
	return ret, nil
}

// KeyGeneratorFunc is a function that generates identifiers
type KeyGeneratorFunc func(string) uint64

// dreateSchema crreates the schema for the database
func createSchema(db *sql.DB) error {
	commands := schemaCommandList()
	for _, v := range commands {
		if _, err := db.Exec(v); err != nil {
			return err
		}
	}
	return nil
}

// putFunc is a function used by the dbSQLExec wrappers
type stmtFunc func(stmt *sql.Stmt) (sql.Result, error)

func doSQLExec(db *sql.DB, statement *sql.Stmt, execFunc stmtFunc) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt := tx.Stmt(statement)
	defer stmt.Close()
	var result sql.Result
	if result, err = execFunc(stmt); err != nil {
		tx.Rollback()
		errMsg := err.Error()
		if strings.Index(errMsg, "duplicate key value violates") > 0 {
			return ErrAlreadyExists
		}
		if strings.Index(errMsg, "violates foreign key constraint") > 0 {
			return ErrDeleteConstraint
		}
		return err
	}
	tx.Commit()
	if count, _ := result.RowsAffected(); count == 0 {
		return ErrNotFound
	}
	return nil
}

// NewMemoryStorage creates a memory-backed storage based on SQLite3
func NewMemoryStorage() *Storage {
	s, err := CreateStorage(":memory:", 12, 10, time.Hour)
	if err != nil {
		panic(err.Error())
	}
	return s
}
