package protocol

import "testing"

func TestPingSlotInfoReq(t *testing.T) {
	m := MACPingSlotInfoReq{macBase{PingSlotInfoReq, false}, 0x05, 0x04}
	macCommandStandardTests(&m, PingSlotInfoReq, t)

	buffer := make([]byte, 3)
	pos := 0
	err := m.encode(buffer, &pos)
	if err != nil {
		t.Error("Could not encode PingSlotInfoReq: ", err)
	}

	p := MACPingSlotInfoReq{macBase{PingSlotInfoReq, false}, 0, 0}
	dpos := 0
	err = p.decode(buffer, &dpos)
	if err != nil {
		t.Error("Could not decode PingSlotInfoReq: ", err)
	}

	if p != m {
		t.Errorf("PingSlotInfoReq encoded and decoded are different: %v != %v", p, m)
	}

	if dpos != pos {
		t.Errorf("PingSlotInfoReq decodes different number of bytes (%d != %d)", dpos, pos)
	}
}

func TestPingSlotInfoAns(t *testing.T) {
	m := MACPingSlotInfoAns{macBase{PingSlotInfoAns, false}}
	macCommandStandardTests(&m, PingSlotInfoAns, t)

	buffer := make([]byte, 2)
	pos := 0
	err := m.encode(buffer, &pos)
	if err != nil {
		t.Error("Could not encode PingSlotInfoAns: ", err)
	}

	p := MACPingSlotInfoAns{macBase{PingSlotInfoAns, false}}
	dpos := 0
	err = p.decode(buffer, &dpos)
	if err != nil {
		t.Error("Could not decode PingSlotInfoAns: ", err)
	}
	if dpos != pos {
		t.Errorf("PingSlotInfoAns decodes different number of bytes (%d != %d)", dpos, pos)
	}
}

func TestBeaconFreqReq(t *testing.T) {
	m := MACBeaconFreqReq{macBase{BeaconFreqReq, false}, 0x0A0B0C}
	macCommandStandardTests(&m, BeaconFreqReq, t)

	buffer := make([]byte, 5)
	pos := 0
	err := m.encode(buffer, &pos)
	if err != nil {
		t.Error("Could not encode BeaconFreqReq: ", err)
	}

	p := MACBeaconFreqReq{macBase{BeaconFreqReq, false}, 0}
	dpos := 0
	err = p.decode(buffer, &dpos)

	if err != nil {
		t.Error("Could not encode BeaconFreqReq: ", err)
	}

	if p != m {
		t.Errorf("BeaconFreqReq decoded incorrectly: %v != %v", p, m)
	}

	if dpos != pos {
		t.Errorf("BeaconFreqReq decodes different number of bytes (%d != %d)", dpos, pos)
	}
}

func TestBeaconFreqAns(t *testing.T) {
	m := MACBeaconFreqAns{macBase{BeaconFreqAns, false}}
	macCommandStandardTests(&m, BeaconFreqAns, t)

	buffer := make([]byte, 2)
	pos := 0
	err := m.encode(buffer, &pos)
	if err != nil {
		t.Error("Could not encode BeaconFreqAns: ", err)
	}

	p := MACBeaconFreqAns{macBase{BeaconFreqAns, false}}
	dpos := 0
	err = p.decode(buffer, &dpos)
	if err != nil {
		t.Error("Could not decode BeaconFreqAns: ", err)
	}

	if dpos != pos {
		t.Errorf("BeaconFreqAns decodes different number of bytes (%d != %d)", dpos, pos)
	}
}

func TestPingSlotChannelReq(t *testing.T) {
	m := MACPingSlotChannelReq{macBase{PingSlotChannelReq, false}, 0x010203, 0x4, 0x5}
	macCommandStandardTests(&m, PingSlotChannelReq, t)

	buffer := make([]byte, 12)
	pos := 0
	err := m.encode(buffer, &pos)
	if err != nil {
		t.Error("Could not encode PingSlotChannelReq: ", err)
	}

	p := MACPingSlotChannelReq{macBase{PingSlotChannelReq, false}, 0, 0, 0}
	dpos := 0
	err = p.decode(buffer, &dpos)
	if err != nil {
		t.Error("Could not decode PingSlotChannelReq: ", err)
	}

	if p.Frequency != m.Frequency || p.MaxDR != m.MaxDR || p.MinDR != m.MinDR {
		t.Errorf("PingSlotChannelReq decoded incorrectly: %v != %v", p, m)
	}

	if dpos != pos {
		t.Errorf("PingSlotChannelReq decodes different number of bytes (%d != %d)", dpos, pos)
	}
}

func TestPingSlotFreqAns(t *testing.T) {
	m := MACPingSlotFreqAns{macBase{PingSlotFreqAns, false}, true, false}
	macCommandStandardTests(&m, PingSlotFreqAns, t)

	buffer := make([]byte, 3)
	pos := 0
	err := m.encode(buffer, &pos)
	if err != nil {
		t.Error("Could not encode PingSlotFreqAns: ", err)
	}

	p := MACPingSlotFreqAns{macBase{PingSlotFreqAns, false}, false, false}
	dpos := 0
	err = p.decode(buffer, &dpos)
	if err != nil {
		t.Error("Could not decode PingSlotFreqAns: ", err)
	}

	if p != m {
		t.Errorf("PingSlotChannelReq decoded incorrectly: %v != %v", p, m)
	}

	if dpos != pos {
		t.Errorf("PingSlotChannelReq decodes different number of bytes (%d != %d)", dpos, pos)
	}
}

func TestBeaconTimingReq(t *testing.T) {
	m := MACBeaconTimingReq{macBase{BeaconTimingReq, false}}
	macCommandStandardTests(&m, BeaconTimingReq, t)

	buffer := make([]byte, 2)
	pos := 0
	err := m.encode(buffer, &pos)
	if err != nil {
		t.Error("Could not encode BeaconTimingReq: ", err)
	}

	p := MACBeaconTimingReq{macBase{BeaconTimingReq, false}}
	dpos := 0
	err = p.decode(buffer, &dpos)
	if err != nil {
		t.Error("Could not decode BeaconTimingReq: ", err)
	}
	if dpos != pos {
		t.Errorf("BeaconTimingReq decodes different number of bytes (%d != %d)", dpos, pos)
	}
}

func TestBeaconTimingAns(t *testing.T) {
	m := MACBeaconTimingAns{macBase{BeaconTimingAns, false}, 1, 2}
	macCommandStandardTests(&m, BeaconTimingAns, t)

	buffer := make([]byte, 5)
	pos := 0
	err := m.encode(buffer, &pos)
	if err != nil {
		t.Error("Could not encode BeaconTimingAns: ", err)
	}

	p := MACBeaconTimingAns{macBase{BeaconTimingAns, false}, 0, 0}
	dpos := 0
	err = p.decode(buffer, &dpos)
	if err != nil {
		t.Error("Could not decode BeaconTimingAns: ", err)
	}

	if p != m {
		t.Errorf("BeaconTimingAns decoded incorrectly: %v != %v", p, m)
	}

	if pos != dpos {
		t.Errorf("BeaconTimingAns decodes different number of bytes (%d != %d)", dpos, pos)
	}
}
