package storage

import (
	"database/sql"
	"fmt"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

type deviceStatements struct {
	putStatement         *sql.Stmt
	devAddrStatement     *sql.Stmt
	euiStatement         *sql.Stmt
	nonceStatement       *sql.Stmt
	appEUIStatement      *sql.Stmt
	getNonceStatement    *sql.Stmt
	updateStateStatement *sql.Stmt
	deleteStatement      *sql.Stmt
	updateStatement      *sql.Stmt
}

func (d *deviceStatements) Close() {
	d.putStatement.Close()
	d.devAddrStatement.Close()
	d.euiStatement.Close()
	d.nonceStatement.Close()
	d.appEUIStatement.Close()
	d.getNonceStatement.Close()
	d.updateStateStatement.Close()
	d.deleteStatement.Close()
	d.updateStatement.Close()
}

func (d *deviceStatements) prepare(db *sql.DB) error {
	var err error

	sqlInsert := `
		INSERT INTO
			lora_devices (
				eui,
				dev_addr,
				app_key,
				apps_key,
				nwks_key,
				application_eui,
				state,
				fcnt_up,
				fcnt_dn,
				relaxed_counter,
				key_warning)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11)`
	if d.putStatement, err = db.Prepare(sqlInsert); err != nil {
		return fmt.Errorf("unable to prepare insert statement: %v", err)
	}

	sqlSelect := `
		SELECT
			eui,
			dev_addr,
			app_key,
			apps_key,
			nwks_key,
			application_eui,
			state,
			fcnt_up,
			fcnt_dn,
			relaxed_counter,
			key_warning
		FROM
			lora_devices
		WHERE
			dev_addr = $1`
	if d.devAddrStatement, err = db.Prepare(sqlSelect); err != nil {
		return fmt.Errorf("unable to prepare select statement: %v", err)
	}

	sqlList := `
		SELECT
			eui,
			dev_addr,
			app_key,
			apps_key,
			nwks_key,
			application_eui,
			state,
			fcnt_up,
			fcnt_dn,
			relaxed_counter,
			key_warning
		FROM
			lora_devices
		WHERE
			application_eui = $1`

	if d.appEUIStatement, err = db.Prepare(sqlList); err != nil {
		return fmt.Errorf("unable to prepare list statement: %v", err)
	}

	euiSelect := `
		SELECT
			eui,
			dev_addr,
			app_key,
			apps_key,
			nwks_key,
			application_eui,
			state,
			fcnt_up,
			fcnt_dn,
			relaxed_counter,
			key_warning
		FROM
			lora_devices
		WHERE
			eui = $1`

	if d.euiStatement, err = db.Prepare(euiSelect); err != nil {
		return fmt.Errorf("unable to prepare eui select statement: %v", err)
	}

	nonceInsert := `INSERT INTO lora_device_nonces (device_eui, nonce) VALUES ($1, $2)`
	if d.nonceStatement, err = db.Prepare(nonceInsert); err != nil {
		return fmt.Errorf("unable to prepare nonce insert statement: %v", err)
	}

	nonceSelect := `SELECT nonce FROM lora_device_nonces WHERE device_eui = $1`
	if d.getNonceStatement, err = db.Prepare(nonceSelect); err != nil {
		return fmt.Errorf("unable to prepare nonce select statement: %v", err)
	}

	updateState := `UPDATE lora_devices SET fcnt_dn = $1, fcnt_up = $2, key_warning = $3 WHERE eui = $4`
	if d.updateStateStatement, err = db.Prepare(updateState); err != nil {
		return fmt.Errorf("unable to prepare update state statement: %v", err)
	}

	delete := `DELETE FROM lora_devices WHERE eui = $1`
	if d.deleteStatement, err = db.Prepare(delete); err != nil {
		return fmt.Errorf("unable to prepare delete statement: %v", err)
	}

	update := `
		UPDATE
			lora_devices
		SET
			dev_addr = $1,
			app_key = $2,
			apps_key = $3,
			nwks_key = $4,
			state = $5,
			fcnt_up = $6,
			fcnt_dn = $7,
			relaxed_counter = $8,
			key_warning = $9
		WHERE eui = $10`
	if d.updateStatement, err = db.Prepare(update); err != nil {
		return fmt.Errorf("unable to prepare device update statement: %v", err)
	}
	return nil
}

// Read nonces for device.
func (s *Storage) retrieveNonces(device *model.Device) error {
	rows, err := s.devStmt.getNonceStatement.Query(device.DeviceEUI.ToInt64())
	if err != nil {
		return fmt.Errorf("unable to query nonces: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var nonce int
		if err := rows.Scan(&nonce); err != nil {
			lg.Warning("Unable to read DevNonce for device with EUI %s: %v", device.DeviceEUI, err)
			continue
		}
		device.DevNonceHistory = append(device.DevNonceHistory, uint16(nonce))
	}
	return nil
}

func (s *Storage) readDeviceSansNonce(row *sql.Rows) (model.Device, error) {
	ret := model.Device{}
	var devAddrStr, appKeyStr, appSkeyStr, nwkSkeyStr string
	var devEUI, appEUI int64
	var err error
	if err = row.Scan(
		&devEUI,
		&devAddrStr,
		&appKeyStr,
		&appSkeyStr,
		&nwkSkeyStr,
		&appEUI,
		&ret.State,
		&ret.FCntUp,
		&ret.FCntDn,
		&ret.RelaxedCounter,
		&ret.KeyWarning); err != nil {
		return ret, err
	}

	ret.DeviceEUI = protocol.EUIFromInt64(devEUI)

	if ret.DevAddr, err = protocol.DevAddrFromString(devAddrStr); err != nil {
		return ret, fmt.Errorf("invalid DevAddr for device with EUI %s (devaddr=%s)", ret.DeviceEUI, devAddrStr)
	}
	ret.AppEUI = protocol.EUIFromInt64(appEUI)

	if ret.AppKey, err = protocol.AESKeyFromString(appKeyStr); err != nil {
		return ret, fmt.Errorf("invalid AppKey: %v (key=%s)", err, appKeyStr)
	}
	if ret.AppSKey, err = protocol.AESKeyFromString(appSkeyStr); err != nil {
		return ret, fmt.Errorf("invalid AppSKey: %v (key=%s)", err, appSkeyStr)
	}
	if ret.NwkSKey, err = protocol.AESKeyFromString(nwkSkeyStr); err != nil {
		return ret, fmt.Errorf("invalid NwkSKey: %v (key=%s)", err, nwkSkeyStr)
	}

	return ret, nil
}

func (s *Storage) getDevice(rows *sql.Rows, err error) (model.Device, error) {
	emptyDevice := model.Device{}

	if err != nil {
		return emptyDevice, err
	}
	if !rows.Next() {
		rows.Close()
		return emptyDevice, ErrNotFound
	}
	device, err := s.readDeviceSansNonce(rows)
	if err != nil {
		defer rows.Close()
		return emptyDevice, err
	}
	rows.Close()
	return device, s.retrieveNonces(&device)

}

func (s *Storage) getDeviceList(rows *sql.Rows, err error) ([]model.Device, error) {
	if err != nil {
		return nil, fmt.Errorf("unable to query device list: %v", err)
	}
	var ret []model.Device

	for rows.Next() {
		device, err := s.readDeviceSansNonce(rows)
		if err != nil {
			rows.Close()
			return ret, err
		}
		ret = append(ret, device)
	}
	rows.Close()
	for i := range ret {
		if err := s.retrieveNonces(&ret[i]); err != nil {
			return ret, err
		}
	}
	return ret, nil
}

// GetDeviceByDevAddr returns the device with the matching device address
func (s *Storage) GetDeviceByDevAddr(devAddr protocol.DevAddr) ([]model.Device, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.getDeviceList(s.devStmt.devAddrStatement.Query(devAddr.String()))
}

// GetDeviceByEUI retrieves a device by its EUI
func (s *Storage) GetDeviceByEUI(devEUI protocol.EUI) (model.Device, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.getDevice(s.devStmt.euiStatement.Query(devEUI.ToInt64()))
}

// GetDevicesByApplicationEUI returns all devices for the given application
func (s *Storage) GetDevicesByApplicationEUI(appEUI protocol.EUI) ([]model.Device, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.getDeviceList(s.devStmt.appEUIStatement.Query(appEUI.ToInt64()))
}

// CreateDevice creates a device in the store
func (s *Storage) CreateDevice(device model.Device, appEUI protocol.EUI) error {
	return s.doSQLExec(s.devStmt.putStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(device.DeviceEUI.ToInt64(),
			device.DevAddr.String(),
			device.AppKey.String(),
			device.AppSKey.String(),
			device.NwkSKey.String(),
			device.AppEUI.ToInt64(),
			uint8(device.State),
			device.FCntUp,
			device.FCntDn,
			device.RelaxedCounter,
			device.KeyWarning)
	})
}

// AddDevNonce adds a nonce to the device
func (s *Storage) AddDevNonce(device model.Device, nonce uint16) error {
	return s.doSQLExec(s.devStmt.nonceStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(device.DeviceEUI.ToInt64(), nonce)
	})
}

// UpdateDeviceState updates the device state in the store
func (s *Storage) UpdateDeviceState(device model.Device) error {
	return s.doSQLExec(s.devStmt.updateStateStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(device.FCntDn, device.FCntUp, device.KeyWarning, device.DeviceEUI.ToInt64())
	})
}

// DeleteDevice removes a device from the store
func (s *Storage) DeleteDevice(eui protocol.EUI) error {
	return s.doSQLExec(s.devStmt.deleteStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(eui.ToInt64())
	})
}

// UpdateDevice updates the device
func (s *Storage) UpdateDevice(device model.Device) error {
	return s.doSQLExec(s.devStmt.updateStatement, func(st *sql.Stmt) (sql.Result, error) {
		return st.Exec(
			device.DevAddr.String(),
			device.AppKey.String(),
			device.AppSKey.String(),
			device.NwkSKey.String(),
			uint8(device.State),
			device.FCntUp,
			device.FCntDn,
			device.RelaxedCounter,
			device.KeyWarning,
			device.DeviceEUI.ToInt64())
	})
}
