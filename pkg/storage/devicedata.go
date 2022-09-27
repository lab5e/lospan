package storage

import (
	"database/sql"
	"fmt"

	"encoding/base64"

	"github.com/ExploratoryEngineering/logging"
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
				data,
				time_stamp,
				gateway_eui,
				rssi,
				snr,
				frequency,
				data_rate,
				dev_addr)
		VALUES ($1,	$2,	$3, $4, $5, $6, $7, $8, $9)`
	if d.putStatement, err = db.Prepare(sqlInsert); err != nil {
		return fmt.Errorf("unable to prepare insert statement: %v", err)
	}

	sqlSelect := `
		SELECT
			device_eui,
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
		SELECT d.device_eui, d.data, d.time_stamp, gateway_eui, rssi, snr, frequency, data_rate, d.dev_addr
		FROM lora_device_data d
			INNER JOIN lora_device dev ON d.device_eui = dev.eui
			INNER JOIN lora_application app ON dev.application_eui = app.eui
		WHERE app.eui = $1
			ORDER BY d.time_stamp DESC
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

// Put stores a new data element in the backend. The element is associated with the specified DevAddr
func (s *Storage) CreateUpstreamData(deviceEUI protocol.EUI, data model.DeviceData) error {
	return doSQLExec(s.db, s.dataStmt.putStatement, func(st *sql.Stmt) (sql.Result, error) {
		b64str := base64.StdEncoding.EncodeToString(data.Data)
		return st.Exec(deviceEUI.String(),
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
	var devEUI, dataStr, gwEUI, devAddr string
	if err = rows.Scan(&devEUI, &dataStr, &ret.Timestamp, &gwEUI, &ret.RSSI, &ret.SNR, &ret.Frequency, &ret.DataRate, &devAddr); err != nil {
		return ret, err
	}
	if ret.DeviceEUI, err = protocol.EUIFromString(devEUI); err != nil {
		return ret, err
	}
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

func (s *Storage) doQuery(stmt *sql.Stmt, eui string, limit int) (chan model.DeviceData, error) {
	rows, err := stmt.Query(eui, limit)
	if err != nil {
		return nil, fmt.Errorf("unable to query device data for device with EUI %s: %v", eui, err)
	}
	outputChan := make(chan model.DeviceData)
	go func() {
		defer rows.Close()
		defer close(outputChan)
		for rows.Next() {
			ret, err := s.readData(rows)
			if err != nil {
				logging.Warning("Unable to decode data for device with EUI %s: %v", eui, err)
				continue
			}
			outputChan <- ret
		}
	}()
	return outputChan, nil
}

// GetUpstreamDataByDeviceEUI retrieves all of the data stored for that DevAddr
func (s *Storage) GetUpstreamDataByDeviceEUI(deviceEUI protocol.EUI, limit int) (chan model.DeviceData, error) {
	return s.doQuery(s.dataStmt.listStatement, deviceEUI.String(), limit)
}

func (s *Storage) GetDownstreamByApplicationEUI(applicationEUI protocol.EUI, limit int) (chan model.DeviceData, error) {
	return s.doQuery(s.dataStmt.appDataList, applicationEUI.String(), limit)
}

// CreateDownstreamData creates new downstream data for a device
func (s *Storage) CreateDownstreamData(deviceEUI protocol.EUI, message model.DownstreamMessage) error {
	return doSQLExec(s.db, s.dataStmt.putDownstream, func(st *sql.Stmt) (sql.Result, error) {
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
	return doSQLExec(s.db, s.dataStmt.deleteDownstream, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(deviceEUI.String())
	})
}

// GetDownstreamData returns a downstream message
func (s *Storage) GetDownstreamData(deviceEUI protocol.EUI) (model.DownstreamMessage, error) {
	ret := model.NewDownstreamMessage(deviceEUI, 0)

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
	return doSQLExec(s.db, s.dataStmt.updateDownstream, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(
			sentTime,
			ackTime,
			deviceEUI.String())
	})
}
