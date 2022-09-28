package storage

import (
	"database/sql"
	"fmt"

	"encoding/base64"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

type dataStatements struct {
	putStatement     *sql.Stmt
	listStatement    *sql.Stmt
	appDataList      *sql.Stmt
	putDownstream    *sql.Stmt
	deleteDownstream *sql.Stmt
	updateDownstream *sql.Stmt
	getDownstream    *sql.Stmt
}

// Close closes the resources opened by the DBDataStorage instance
func (d *dataStatements) Close() {
	d.putStatement.Close()
	d.listStatement.Close()
	d.appDataList.Close()
	d.putDownstream.Close()
	d.deleteDownstream.Close()
	d.updateDownstream.Close()
	d.getDownstream.Close()
}

func (d *dataStatements) prepare(db *sql.DB) error {
	var err error

	sqlInsert := `
		INSERT INTO
			lora_device_data (
				device_eui,
				application_eui,
				data,
				time_stamp,
				gateway_eui,
				rssi,
				snr,
				frequency,
				data_rate,
				dev_addr)
		VALUES ($1,	$2,	$3, $4, $5, $6, $7, $8, $9, $10)`
	if d.putStatement, err = db.Prepare(sqlInsert); err != nil {
		return fmt.Errorf("unable to prepare insert statement: %v", err)
	}

	sqlSelect := `
		SELECT
			device_eui,
			application_eui,
			data,
			time_stamp,
			gateway_eui,
			rssi,
			snr,
			frequency,
			data_rate,
			dev_addr
		FROM
			lora_device_data
		WHERE
			device_eui = $1
		ORDER BY
			time_stamp DESC
		LIMIT $2`
	if d.listStatement, err = db.Prepare(sqlSelect); err != nil {
		return fmt.Errorf("unable to prepare list statement: %v", err)
	}

	sqlDataList := `
		SELECT 
			device_eui, 
			application_eui, 
			data, 
			time_stamp, 
			gateway_eui, 
			rssi, 
			snr, 
			frequency, 
			data_rate, 
			dev_addr
		FROM 
			lora_device_data 
		WHERE 
			application_eui = $1
		ORDER BY 
			time_stamp DESC
		LIMIT $2`
	if d.appDataList, err = db.Prepare(sqlDataList); err != nil {
		return fmt.Errorf("unable to prepare app list statement: %v", err)
	}

	sqlPutDownstream := `
		INSERT INTO lora_downstream_message (
			device_eui,
			data,
			port,
			ack,
			created_time,
			sent_time,
			ack_time)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7)
	`
	if d.putDownstream, err = db.Prepare(sqlPutDownstream); err != nil {
		return fmt.Errorf("unable to prepare downstream put statement: %v", err)
	}

	sqlDeleteDownsteram := `
		DELETE FROM
			lora_downstream_message
		WHERE
			device_eui = $1
	`
	if d.deleteDownstream, err = db.Prepare(sqlDeleteDownsteram); err != nil {
		return fmt.Errorf("unable to prepare downstream delete statement: %v", err)
	}

	sqlUpdateDownstream := `
		UPDATE lora_downstream_message
			SET
				sent_time = $1,
				ack_time = $2
			WHERE
				device_eui = $3
	`
	if d.updateDownstream, err = db.Prepare(sqlUpdateDownstream); err != nil {
		return fmt.Errorf("unable to prepare downstream update statement")
	}

	sqlGetDownstream := `
		SELECT
			data,
			port,
			ack,
			created_time,
			sent_time,
			ack_time
		FROM
			lora_downstream_message
		WHERE
			device_eui = $1
	`
	if d.getDownstream, err = db.Prepare(sqlGetDownstream); err != nil {
		return fmt.Errorf("unable to prepare downstream select statement")
	}
	return nil
}

// CreateUpstreamData stores a new data element in the backend. The element is associated with the specified DevAddr
func (s *Storage) CreateUpstreamData(deviceEUI protocol.EUI, applicationEUI protocol.EUI, data model.DeviceData) error {
	return s.doSQLExec(s.dataStmt.putStatement, func(st *sql.Stmt) (sql.Result, error) {
		b64str := base64.StdEncoding.EncodeToString(data.Data)
		return st.Exec(deviceEUI.ToInt64(),
			applicationEUI.ToInt64(),
			b64str,
			data.Timestamp,
			data.GatewayEUI.String(),
			data.RSSI,
			data.SNR,
			data.Frequency,
			data.DataRate,
			data.DevAddr.String())
	})
}

// Decode a single row into a DeviceData instance.
func (s *Storage) readData(rows *sql.Rows) (model.DeviceData, error) {
	ret := model.DeviceData{}
	var err error
	var dataStr, gwEUI, devAddr string
	var devEUI, appEUI int64
	if err = rows.Scan(&devEUI, &appEUI, &dataStr, &ret.Timestamp, &gwEUI, &ret.RSSI, &ret.SNR, &ret.Frequency, &ret.DataRate, &devAddr); err != nil {
		return ret, err
	}
	ret.DeviceEUI = protocol.EUIFromInt64(devEUI)
	ret.AppEUI = protocol.EUIFromInt64(appEUI)
	if ret.Data, err = base64.StdEncoding.DecodeString(dataStr); err != nil {
		return ret, err
	}
	if ret.GatewayEUI, err = protocol.EUIFromString(gwEUI); err != nil {
		return ret, err
	}
	if ret.DevAddr, err = protocol.DevAddrFromString(devAddr); err != nil {
		return ret, err
	}
	return ret, nil
}

func (s *Storage) doQuery(stmt *sql.Stmt, eui protocol.EUI, limit int) ([]model.DeviceData, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	rows, err := stmt.Query(eui.ToInt64(), limit)
	if err != nil {
		return nil, fmt.Errorf("unable to query device data for device with EUI %s: %v", eui, err)
	}
	var ret []model.DeviceData
	defer rows.Close()
	for rows.Next() {
		data, err := s.readData(rows)
		if err != nil {
			lg.Warning("Unable to decode data for device with EUI %s: %v", eui, err)
			return ret, err
		}
		ret = append(ret, data)
	}

	return ret, nil
}

// GetUpstreamDataByDeviceEUI retrieves all of the data stored for that DevAddr
func (s *Storage) GetUpstreamDataByDeviceEUI(deviceEUI protocol.EUI, limit int) ([]model.DeviceData, error) {
	return s.doQuery(s.dataStmt.listStatement, deviceEUI, limit)
}

// GetDownstreamDataByApplicationEUI returns
func (s *Storage) GetDownstreamDataByApplicationEUI(applicationEUI protocol.EUI, limit int) ([]model.DeviceData, error) {
	return s.doQuery(s.dataStmt.appDataList, applicationEUI, limit)
}

// CreateDownstreamData creates new downstream data for a device
func (s *Storage) CreateDownstreamData(deviceEUI protocol.EUI, message model.DownstreamMessage) error {
	return s.doSQLExec(s.dataStmt.putDownstream, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(
			deviceEUI.String(),
			message.Data,
			message.Port,
			message.Ack,
			message.CreatedTime,
			message.SentTime,
			message.AckTime)
	})
}

// DeleteDownstreamData deletes a downstream message
func (s *Storage) DeleteDownstreamData(deviceEUI protocol.EUI) error {
	return s.doSQLExec(s.dataStmt.deleteDownstream, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(deviceEUI.String())
	})
}

// GetDownstreamData returns a downstream message
func (s *Storage) GetDownstreamData(deviceEUI protocol.EUI) (model.DownstreamMessage, error) {
	ret := model.NewDownstreamMessage(deviceEUI, 0)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	rows, err := s.dataStmt.getDownstream.Query(deviceEUI.String())
	if err != nil {
		return ret, fmt.Errorf("unable to query for downstream message: %v", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return ret, ErrNotFound
	}
	if err := rows.Scan(&ret.Data, &ret.Port, &ret.Ack, &ret.CreatedTime, &ret.SentTime, &ret.AckTime); err != nil {
		return ret, fmt.Errorf("unable to read fields from downstream result: %v", err)
	}
	return ret, nil
}

// UpdateDownstreamData updates a downstream message
func (s *Storage) UpdateDownstreamData(deviceEUI protocol.EUI, sentTime int64, ackTime int64) error {
	return s.doSQLExec(s.dataStmt.updateDownstream, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(
			sentTime,
			ackTime,
			deviceEUI.String())
	})
}
