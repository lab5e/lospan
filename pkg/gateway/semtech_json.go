package gateway

// Rxpk is a (JSON) struct used by the Semtech packet forwarder. It is sent from the gateway to the server.
type Rxpk struct {
	Time                string  `json:"time"` // Time stamp (unix-) for the gateway
	Timestamp           uint32  `json:"tmst"`
	Frequency           float32 `json:"freq"`
	ConcentratorChannel uint8   `json:"chan"`
	ConcentratorRFChain uint8   `json:"rfch"`
	ModulationID        string  `json:"modu"`
	DataRateID          string  `json:"datr"`
	CodingRateID        string  `json:"codr"`
	RSSI                int32   `json:"rssi"`
	LoraSNRRatio        float32 `json:"lsnr"`
	PayloadSize         uint32  `json:"size"`
	RFPackets           string  `json:"data"`
}

// Txpk is a (JSON) struct used by the Semtech packet forwarder. It is sent from the server to the gateway
type Txpk struct {
	Immediate             bool    `json:"imme"`           // (one of)Send packet immediately (will ignore tmst & time)
	Timestamp             uint32  `json:"tmst,omitempty"` // (one of)Send packet on a certain timestamp value (will ignore time)
	Time                  string  `json:"time,omitempty"` // (one of)Send packet at a certain time (GPS synchronization required)
	Frequency             float32 `json:"freq"`           // (mandatory)TX central frequency in MHz (unsigned float, Hz precision)
	RFChain               uint8   `json:"rfch"`           // (mandatory)Concentrator "RF chain" used for TX (unsigned integer)
	TxPower               uint32  `json:"powe,omitempty"` // TX output power in dBm (unsigned integer, dBm precision)
	Modulation            string  `json:"modu"`           // (mandatory)Modulation identifier "LORA" or "FSK"
	LoRaDataRate          string  `json:"datr"`           // (mandatory)LoRa datarate identifier (eg. SF12BW500)
	EccCoding             string  `json:"codr,omitempty"` // LoRa ECC coding rate identifier
	FskFrequencyDeviation uint32  `json:"fdev,omitempty"` // FSK frequency deviation (unsigned integer, in Hz)
	LoraInvPol            bool    `json:"ipol,omitempty"` // Lora modulation polarization inversion
	RfPreamble            uint32  `json:"prea,omitempty"` // RF preamble size (unsigned integer)
	PayloadSize           int     `json:"size"`           // (mandatory)RF packet payload size in bytes (unsigned integer)
	Data                  string  `json:"data"`           // (mandatory)Base64 encoded RF packet payload, padding optional
	NoCRC                 bool    `json:"ncrc,omitempty"` // If true, disable the CRC of the physical layer (optional)

}

// RXData contains device payload in "Data" and also (possibly) gateway status in "Stat". Both contain JSON
type RXData struct {
	Data []Rxpk `json:"rxpk"`
	//	Stat statv2 `json:"stat"` // this might be an array or a single value, depending on configuration.
}

// TXData is the struct used when transmitting data to the gateway
type TXData struct {
	Data Txpk `json:"txpk"`
}
