package storage

import (
	"database/sql"
	"fmt"

	"github.com/ExploratoryEngineering/logging"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

type applicationStatements struct {
	putStatement       *sql.Stmt // Prepared statement for Put
	getStatement       *sql.Stmt // Prepared statement for GetByEUI
	listStatement      *sql.Stmt // Prepared statement for GetByNetworkEUI
	deleteStatement    *sql.Stmt // Prepared statement for Delete
	systemGetStatement *sql.Stmt // Prepared statement for system get
}

// Close releases all of the resources used by the application storage.
func (a *applicationStatements) Close() {
	a.putStatement.Close()
	a.getStatement.Close()
	a.listStatement.Close()
	a.deleteStatement.Close()
	a.systemGetStatement.Close()
}

func (a *applicationStatements) prepare(db *sql.DB) error {
	var err error
	sqlInsert := `
		INSERT INTO
			lora_application (eui)
		VALUES ($1)`
	if a.putStatement, err = db.Prepare(sqlInsert); err != nil {
		return fmt.Errorf("unable to prepare insert statement: %v", err)
	}

	sqlSelect := `
		SELECT
			a.eui
		FROM
			lora_application a
		WHERE
			a.eui = $1`
	if a.getStatement, err = db.Prepare(sqlSelect); err != nil {
		return fmt.Errorf("app:unable to prepare select statement: %v", err)
	}

	sqlList := `
		SELECT
			a.eui
		FROM
			lora_application a`

	if a.listStatement, err = db.Prepare(sqlList); err != nil {
		return fmt.Errorf("app:unable to prepare list statement: %v", err)
	}

	sqlDelete := `
		DELETE
		FROM lora_application 
		WHERE eui = $1`
	if a.deleteStatement, err = db.Prepare(sqlDelete); err != nil {
		return fmt.Errorf("app:unable to prepare delete statement: %v", err)
	}

	sqlSystemGet := `
		SELECT
			a.eui
		FROM
			lora_application a
		WHERE
			a.eui = $1`
	if a.systemGetStatement, err = db.Prepare(sqlSystemGet); err != nil {
		return fmt.Errorf("app:unable to prepare system select statement: %v", err)
	}

	return nil
}

func (s *Storage) readApplication(rows *sql.Rows) (model.Application, error) {
	var appEUI int64
	var err error
	ret := model.NewApplication()
	if err = rows.Scan(&appEUI); err != nil {
		return ret, err
	}

	ret.AppEUI = protocol.EUIFromInt64(appEUI)
	return ret, nil
}

// GetByEUI retrieves the application with the specified application EUI.
func (s *Storage) GetApplicationByEUI(eui protocol.EUI) (model.Application, error) {
	rows, err := s.appStmt.systemGetStatement.Query(eui.ToInt64())
	ret := model.NewApplication()
	if err != nil {
		return ret, err
	}
	defer rows.Close()
	if !rows.Next() {
		return ret, ErrNotFound
	}
	app, err := s.readApplication(rows)
	return app, err
}

// ListApplications returns all applications with the given network EUI
func (s *Storage) ListApplications() (chan model.Application, error) {
	rows, err := s.appStmt.listStatement.Query()
	if err != nil {
		return nil, fmt.Errorf("unable to query application list: %v", err)
	}
	outputChan := make(chan model.Application)
	go func() {
		defer rows.Close()
		defer close(outputChan)
		for rows.Next() {
			app, err := s.readApplication(rows)
			if err != nil {
				logging.Warning("Unable to read application in list, skipping: %v", err)
				continue
			}
			outputChan <- app
		}
	}()
	return outputChan, nil
}

// Put stores an Application instance in the storage backend
func (s *Storage) CreateApplication(application model.Application) error {
	return doSQLExec(s.db, s.appStmt.putStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(application.AppEUI.ToInt64())
	})
}

func (s *Storage) DeleteApplication(eui protocol.EUI) error {
	return doSQLExec(s.db, s.appStmt.deleteStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(eui.ToInt64())
	})
}
