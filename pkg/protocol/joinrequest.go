package protocol

import (
	"encoding/binary"
)

// JoinRequestPayload is the payload sent by the device in a JoinRequest
// message [6.2.4]. The message is not encrypted.
type JoinRequestPayload struct {
	AppEUI   EUI
	DevEUI   EUI
	DevNonce uint16
}

// Decode JoinRequest payload from a byte buffer.
func (j *JoinRequestPayload) decode(buffer []byte, pos *int) error {
	if buffer == nil || pos == nil {
		return ErrNilError
	}
	if len(buffer) < (*pos + 18) {
		return ErrBufferTruncated
	}
	j.AppEUI = EUIFromInt64(int64(binary.LittleEndian.Uint64(buffer[*pos:])))
	*pos += 8
	j.DevEUI = EUIFromInt64(int64(binary.LittleEndian.Uint64(buffer[*pos:])))
	*pos += 8

	// These should be big endian. Because keys.
	j.DevNonce = binary.BigEndian.Uint16(buffer[*pos:])
	*pos += 2
	return nil
}

// Encode the JoinRequest into a buffer
func (j *JoinRequestPayload) encode(buffer []byte, pos *int) error {
	if buffer == nil || pos == nil {
		return ErrNilError
	}
	if len(buffer) < (*pos + 18) {
		return ErrBufferTruncated
	}

	binary.LittleEndian.PutUint64(buffer[*pos:], uint64(j.AppEUI.ToInt64()))
	*pos += 8
	binary.LittleEndian.PutUint64(buffer[*pos:], uint64(j.DevEUI.ToInt64()))
	*pos += 8

	// These should be big endian. Because keys.
	binary.BigEndian.PutUint16(buffer[*pos:], j.DevNonce)
	*pos += 2
	return nil
}
