package storage

import (
	"fmt"

	"database/sql"

	"net"

	"github.com/ExploratoryEngineering/logging"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

type gatewayStatements struct {
	putStatement    *sql.Stmt // Prepare statement for put operation
	deleteStatement *sql.Stmt // Prepare statement for delete operation
	listStatement   *sql.Stmt // Prepare statement for select
	getStatement    *sql.Stmt // Prepare statement for select
	getSysStatement *sql.Stmt // Prepare statement for system get (ie all gateways)
	updateStatement *sql.Stmt // Prepare statement for gatway update
}

func (g *gatewayStatements) Close() {
	g.putStatement.Close()
	g.deleteStatement.Close()
	g.listStatement.Close()
	g.getStatement.Close()
	g.getSysStatement.Close()
	g.updateStatement.Close()
}

func (g *gatewayStatements) prepare(db *sql.DB) error {
	var err error
	sqlSelect := `
		SELECT
			gateway_eui,
			latitude,
			longitude,
			altitude,
			ip,
			strict_ip
		FROM
			lora_gateway`

	if g.listStatement, err = db.Prepare(sqlSelect); err != nil {
		return fmt.Errorf("unable to prepare list statement: %v", err)
	}

	sqlInsert := `
		INSERT INTO lora_gateway (
			gateway_eui,
			latitude,
			longitude,
			altitude,
			ip,
			strict_ip)
		VALUES ($1, $2, $3, $4, $5, $6)`
	if g.putStatement, err = db.Prepare(sqlInsert); err != nil {
		return fmt.Errorf("unable to prepare insert statement: %v", err)
	}

	sqlDelete := `
		DELETE FROM
			lora_gateway 
		WHERE
			gateway_eui = $1`
	if g.deleteStatement, err = db.Prepare(sqlDelete); err != nil {
		return fmt.Errorf("unable to prepare delete statement: %v", err)
	}

	sqlSelectOne := `
		SELECT
			gw.gateway_eui,
			gw.latitude,
			gw.longitude,
			gw.altitude,
			gw.ip,
			gw.strict_ip
		FROM
			lora_gateway gw
		WHERE
			gw.gateway_eui = $1`
	if g.getStatement, err = db.Prepare(sqlSelectOne); err != nil {
		return fmt.Errorf("unable to prepare select statement: %v", err)
	}

	sysGetStatement := `
		SELECT
			gw.gateway_eui,
			gw.latitude,
			gw.longitude,
			gw.altitude,
			gw.ip,
			gw.strict_ip
		FROM
			lora_gateway gw
		WHERE
			gw.gateway_eui = $1`
	if g.getSysStatement, err = db.Prepare(sysGetStatement); err != nil {
		return fmt.Errorf("unable to prepare system get statement: %v", err)
	}

	updateStatement := `
		UPDATE
			lora_gateway 
		SET
			latitude = $1, longitude = $2, altitude = $3, ip = $4, strict_ip = $5
		WHERE
			gateway_eui = $6
	`
	if g.updateStatement, err = db.Prepare(updateStatement); err != nil {
		return fmt.Errorf("unable to prepare update statement: %v", err)
	}

	return nil
}

func (s *Storage) readGateway(rows *sql.Rows) (model.Gateway, error) {
	var euiStr, ipStr string
	var err error
	var json []uint8
	gw := model.NewGateway()
	if err := rows.Scan(&euiStr, &gw.Latitude, &gw.Longitude, &gw.Altitude, &ipStr, &gw.StrictIP, &json); err != nil {
		return gw, err
	}
	if gw.GatewayEUI, err = protocol.EUIFromString(euiStr); err != nil {
		return gw, err
	}
	gw.IP = net.ParseIP(ipStr)
	return gw, nil
}

func (s *Storage) getGwList(rows *sql.Rows, err error) (chan model.Gateway, error) {
	if err != nil {
		return nil, err
	}
	ret := make(chan model.Gateway)
	go func() {
		defer rows.Close()
		defer close(ret)
		for rows.Next() {
			gw, err := s.readGateway(rows)
			if err != nil {
				logging.Warning("Unable to read gateway list: %v", err)
				continue
			}
			ret <- gw
		}
	}()
	return ret, nil
}

// GetGatewayList returns a list of gateways
func (s *Storage) GetGatewayList() (chan model.Gateway, error) {
	return s.getGwList(s.gwStmt.listStatement.Query())
}

func (s *Storage) getGateway(rows *sql.Rows, err error) (model.Gateway, error) {
	if err != nil {
		return model.Gateway{}, err
	}

	defer rows.Close()

	if !rows.Next() {
		return model.Gateway{}, ErrNotFound
	}

	gw, err := s.readGateway(rows)
	if err != nil {
		return model.Gateway{}, err
	}
	return gw, nil
}

// GetGateway returns a gateway from the store
func (s *Storage) GetGateway(eui protocol.EUI) (model.Gateway, error) {
	return s.getGateway(s.gwStmt.getSysStatement.Query(eui.String()))
}

// CreateGateway creates a new gateway in the store
func (s *Storage) CreateGateway(gateway model.Gateway) error {
	return doSQLExec(s.db, s.gwStmt.putStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(
			gateway.GatewayEUI.String(),
			gateway.Latitude,
			gateway.Longitude,
			gateway.Altitude,
			gateway.IP.String(),
			gateway.StrictIP)
	})
}

// DeleteGateway removes a gateway from the store
func (s *Storage) DeleteGateway(eui protocol.EUI) error {
	return doSQLExec(s.db, s.gwStmt.deleteStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(eui.String())
	})
}

// UpdateGateway updates a gateway in the store
func (s *Storage) UpdateGateway(gateway model.Gateway) error {
	return doSQLExec(s.db, s.gwStmt.updateStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(gateway.Latitude, gateway.Longitude, gateway.Altitude,
			gateway.IP.String(), gateway.StrictIP, gateway.GatewayEUI.String())
	})
}