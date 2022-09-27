package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ExploratoryEngineering/logging"
)

type keyStatements struct {
	selectStatement *sql.Stmt
	updateStatement *sql.Stmt
	insertStatement *sql.Stmt
}

func (k *keyStatements) Close() {
	k.selectStatement.Close()
	k.updateStatement.Close()
	k.insertStatement.Close()
}

func (k *keyStatements) prepare(db *sql.DB) error {
	sqlSelect := `SELECT counter FROM lora_sequence WHERE identifier = $1`
	sqlUpdate := `UPDATE lora_sequence SET counter = $1 WHERE identifier = $2`
	sqlInsert := `INSERT INTO lora_sequence (identifier, counter) VALUES ($1, $2)`
	var err error
	if k.selectStatement, err = db.Prepare(sqlSelect); err != nil {
		return fmt.Errorf("unable to prepare select statement: %v", err)
	}
	if k.insertStatement, err = db.Prepare(sqlInsert); err != nil {
		return fmt.Errorf("unable to prepare insert statement: %v", err)
	}
	if k.updateStatement, err = db.Prepare(sqlUpdate); err != nil {
		return fmt.Errorf("unable to prepare update statement: %v", err)
	}
	return nil
}

// AllocateKeys allocates a new set of keys from the backend store
func (s *Storage) AllocateKeys(identifier string, interval uint64, initial uint64) (chan uint64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	row := tx.Stmt(s.keyStmt.selectStatement).QueryRow(identifier)

	var start uint64
	var counter int64
	err = row.Scan(&counter)

	switch err {
	case sql.ErrNoRows:
		// Not found - insert a new one with interval prepopulated
		_, err = tx.Stmt(s.keyStmt.insertStatement).Exec(identifier, int64(initial+interval))
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), "lora_sequence_pk") {
				// Retry since the key already exists
				return s.AllocateKeys(identifier, interval, initial)
			}
			return nil, err
		}
		counter = int64(initial)
	default:
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		_, err = tx.Stmt(s.keyStmt.updateStatement).Exec(counter+int64(interval), identifier)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		logging.Error("Unable to commit sequence with identifier %s (interval: %d, initial: %d): %v",
			identifier, interval, initial, err)
	}
	start = uint64(counter)

	ret := make(chan uint64)
	go func() {
		for i := start; i < (start + interval); i++ {
			ret <- i
		}
		close(ret)
	}()
	return ret, nil
}
