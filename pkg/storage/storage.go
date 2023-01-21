package storage

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	_ "modernc.org/sqlite" // use sqlite driver
)

// Storage holds all of the storage objects
type Storage struct {
	db       *sql.DB
	mutex    *sync.Mutex
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
func CreateStorage(connectionString string) (*Storage, error) {
	return newStorage(driverName, connectionString)
}

const driverName = "sqlite3"

func newStorage(driver, connectionString string) (*Storage, error) {
	db, err := sql.Open(driver, connectionString)
	if nil != err {
		log.Fatalf("Unable to connect to database: %s", err)
		return nil, err
	}

	if err := createSchema(db); err != nil {
		return nil, err
	}
	ret := &Storage{
		db:    db,
		mutex: &sync.Mutex{},
	}
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

// NewMemoryStorage creates a memory-backed storage based on SQLite3
func NewMemoryStorage() *Storage {
	s, err := CreateStorage(":memory:")
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (s *Storage) doSQLExec(statement *sql.Stmt, execFunc stmtFunc) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var result sql.Result
	var err error
	if result, err = execFunc(statement); err != nil {
		errMsg := err.Error()
		if strings.Index(errMsg, "duplicate key value violates") > 0 {
			return ErrAlreadyExists
		}
		if strings.Index(errMsg, "violates foreign key constraint") > 0 {
			return ErrDeleteConstraint
		}
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return ErrNotFound
	}
	return nil
}
