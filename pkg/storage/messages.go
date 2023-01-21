package storage

import (
	"database/sql"
	"fmt"

	"encoding/base64"

	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

type dataStatements struct {
	createUpstream   *sql.Stmt
	listUpstream     *sql.Stmt
	createDownstream *sql.Stmt
	deleteDownstream *sql.Stmt
	updateDownstream *sql.Stmt
	listDownstream   *sql.Stmt
}

// Close closes the resources opened by the DBDataStorage instance
func (d *dataStatements) Close() {
	d.createUpstream.Close()
	d.listUpstream.Close()
	d.createDownstream.Close()
	d.deleteDownstream.Close()
	d.updateDownstream.Close()
	d.listDownstream.Close()
}

func (d *dataStatements) prepare(db *sql.DB) error {
	var err error

	if d.createUpstream, err = db.Prepare(`
		INSERT INTO
			lora_upstream_messages (
				device_eui,
				data,
				time_stamp,
				gateway_eui,
				rssi,
				snr,
				frequency,
				data_rate,
				dev_addr)
		VALUES ($1,	$2,	$3, $4, $5, $6, $7, $8, $9)`); err != nil {
		return fmt.Errorf("unable to prepare insert statement: %v", err)
	}

	if d.listUpstream, err = db.Prepare(`
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
			lora_upstream_messages
		WHERE
			device_eui = $1
		ORDER BY
			time_stamp DESC
		LIMIT $2`); err != nil {
		return fmt.Errorf("unable to prepare list statement: %v", err)
	}

	if d.createDownstream, err = db.Prepare(`
		INSERT INTO lora_downstream_messages (
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
			$7)`); err != nil {
		return fmt.Errorf("unable to prepare downstream put statement: %v", err)
	}

	if d.deleteDownstream, err = db.Prepare(`
		DELETE FROM
			lora_downstream_messages
		WHERE
			device_eui = $1`); err != nil {
		return fmt.Errorf("unable to prepare downstream delete statement: %v", err)
	}

	sqlUpdateDownstream := `
		UPDATE lora_downstream_messages
			SET
				sent_time = $1,
				ack_time = $2
			WHERE
				device_eui = $3 
	`
	if d.updateDownstream, err = db.Prepare(sqlUpdateDownstream); err != nil {
		return fmt.Errorf("unable to prepare downstream update statement")
	}

	if d.listDownstream, err = db.Prepare(`
		SELECT
			data,
			port,
			ack,
			created_time,
			sent_time,
			ack_time
		FROM
			lora_downstream_messages
		WHERE
			device_eui = $1
		ORDER BY
			created_time
		LIMIT 100
	`); err != nil {
		return fmt.Errorf("unable to prepare downstream select statement")
	}
	return nil
}

// CreateUpstreamMessage stores a new data element in the backend. The element is associated with the specified DevAddr
func (s *Storage) CreateUpstreamMessage(deviceEUI protocol.EUI, data model.UpstreamMessage) error {
	return s.doSQLExec(s.dataStmt.createUpstream, func(st *sql.Stmt) (sql.Result, error) {
		b64str := base64.StdEncoding.EncodeToString(data.Data)
		return st.Exec(deviceEUI.ToInt64(),
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
func (s *Storage) readData(rows *sql.Rows) (model.UpstreamMessage, error) {
	ret := model.UpstreamMessage{}
	var err error
	var dataStr, gwEUI, devAddr string
	var devEUI int64
	if err = rows.Scan(&devEUI, &dataStr, &ret.Timestamp, &gwEUI, &ret.RSSI, &ret.SNR, &ret.Frequency, &ret.DataRate, &devAddr); err != nil {
		return ret, err
	}
	ret.DeviceEUI = protocol.EUIFromInt64(devEUI)
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

func (s *Storage) doQuery(stmt *sql.Stmt, eui protocol.EUI, limit int) ([]model.UpstreamMessage, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	rows, err := stmt.Query(eui.ToInt64(), limit)
	if err != nil {
		return nil, fmt.Errorf("unable to query device data for device with EUI %s: %v", eui, err)
	}
	var ret []model.UpstreamMessage
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

// ListUpstreamMessages retrieves all of the data stored for that DevAddr
func (s *Storage) ListUpstreamMessages(deviceEUI protocol.EUI, limit int) ([]model.UpstreamMessage, error) {
	return s.doQuery(s.dataStmt.listUpstream, deviceEUI, limit)
}

// CreateDownstreamMessage creates new downstream data for a device
func (s *Storage) CreateDownstreamMessage(deviceEUI protocol.EUI, message model.DownstreamMessage) error {
	return s.doSQLExec(s.dataStmt.createDownstream, func(st *sql.Stmt) (sql.Result, error) {
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

// DeleteDownstreamMessage deletes a downstream message
func (s *Storage) DeleteDownstreamMessage(deviceEUI protocol.EUI) error {
	return s.doSQLExec(s.dataStmt.deleteDownstream, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(deviceEUI.String())
	})
}

// GetNextDownstreamMessage returns a downstream message
func (s *Storage) GetNextDownstreamMessage(deviceEUI protocol.EUI) (model.DownstreamMessage, error) {
	ret := model.NewDownstreamMessage(deviceEUI, 0)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	rows, err := s.dataStmt.listDownstream.Query(deviceEUI.String())
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

// ListDownstreamMessages lists the scheduled downstream messages for a device
func (s *Storage) ListDownstreamMessages(deviceEUI protocol.EUI) ([]model.DownstreamMessage, error) {
	var ret []model.DownstreamMessage

	s.mutex.Lock()
	defer s.mutex.Unlock()

	rows, err := s.dataStmt.listDownstream.Query(deviceEUI.String())
	if err != nil {
		return ret, fmt.Errorf("unable to query for downstream message: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		dm := model.DownstreamMessage{
			DeviceEUI: deviceEUI,
		}
		if err := rows.Scan(&dm.Data, &dm.Port, &dm.Ack, &dm.CreatedTime, &dm.SentTime, &dm.AckTime); err != nil {
			return ret, fmt.Errorf("unable to read fields from downstream result: %v", err)
		}
		ret = append(ret, dm)
	}
	return ret, nil
}

// UpdateDownstreamMessage updates a downstream message
func (s *Storage) UpdateDownstreamMessage(deviceEUI protocol.EUI, sentTime int64, ackTime int64) error {
	return s.doSQLExec(s.dataStmt.updateDownstream, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(
			sentTime,
			ackTime,
			deviceEUI.String())
	})
}
