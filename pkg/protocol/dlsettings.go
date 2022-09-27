package protocol

// DLSettings contains the downlink configuration for the end-device in a JoinAccept
// message [6.2.5]
type DLSettings struct {
	RX1DRoffset byte
	RX2DataRate byte
}

// Encode DLSettings type into buffer.
func (d *DLSettings) encode(buffer []byte, pos *int) error {
	if pos == nil {
		return ErrNilError
	}
	if len(buffer) <= *pos {
		return ErrBufferTruncated
	}
	buffer[*pos] = ((d.RX1DRoffset & 0x07) << 4) | (d.RX2DataRate & 0x0F)
	*pos++
	return nil
}

func (d *DLSettings) decode(buffer []byte, pos *int) error {
	if pos == nil {
		return ErrNilError
	}
	if len(buffer) <= *pos {
		return ErrBufferTruncated
	}
	d.RX1DRoffset = (buffer[*pos] & 0x70) >> 4
	d.RX2DataRate = buffer[*pos] & 0x0F
	*pos++
	return nil
}
