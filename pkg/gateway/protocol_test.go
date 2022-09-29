package gateway

import (
	"testing"

	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/stretchr/testify/require"
)

func TestBinaryMarshal(t *testing.T) {
	assert := require.New(t)

	pkt := GwPacket{}
	// A PUSH_DATA sentence with EUI AABBCCDD and the string 'abcdef'
	buffer := []byte{0, 0x11, 0x22, 0, 0x11, 0xAA, 0xBB, 0xBB, 0xCC, 0xCC, 0xDD, 0xDD, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46}

	err := pkt.UnmarshalBinary(buffer)
	assert.NoError(err, "Error unmarshalling bytes for PUSH_DATA")

	eui, err := protocol.EUIFromString("11-AA-BB-BB-CC-CC-DD-DD")
	assert.NoError(err)

	assert.Equal(pkt.GatewayEUI, eui)
	assert.Equal(uint16(0x1122), pkt.Token)
	assert.Equal("ABCDEF", pkt.JSONString)

	// PUSH_ACK
	buffer = []byte{0, 0x11, 0x22, 1}
	assert.Nil(pkt.UnmarshalBinary(buffer), "PUSH_ACK")

	// PULL_DATA
	buffer = []byte{0, 0x11, 0x22, 2, 0xAA, 0xAA, 0xBB, 0xBB, 0xCC, 0xCC, 0xDD, 0xDD}
	assert.Nil(pkt.UnmarshalBinary(buffer), "PULL_DATA")

	// PULL_RESP
	buffer = []byte{0, 0x11, 0x22, 3, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46}
	assert.Nil(pkt.UnmarshalBinary(buffer), "PULL_RESP")

	// PULL_ACK
	buffer = []byte{0, 0x11, 0x22, 4}
	if pkt.UnmarshalBinary(buffer) != nil {
		t.Fatal("Couldn't unmarshal PULL_ACK")
	}

	// TX_ACK
	buffer = []byte{0, 0x11, 0x22, 5, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46}
	assert.Nil(pkt.UnmarshalBinary(buffer), "TX_ACK")

	// Unknown type
	buffer = []byte{0, 0x11, 0x22, 99}
	assert.NotNil(pkt.UnmarshalBinary(buffer), "Uknown type")

	assert.NotNil(pkt.UnmarshalBinary([]byte{0}), "too small buffer")

	assert.NotNil(pkt.UnmarshalBinary([]byte{1}), "tiny buffer")

	buffer = []byte{0, 0x11, 0x22, 0}
	assert.NotNil(pkt.UnmarshalBinary(buffer), "Too small PUSH_DATA")

	buffer = []byte{0, 0x11, 0x22, 2}
	assert.NotNil(pkt.UnmarshalBinary(buffer), "Shoulnd't be able to unmarshal small PULL_DATA buffer")

}

func TestBinaryUnmarsha(t *testing.T) {
	assert := require.New(t)

	pkt := GwPacket{
		ProtocolVersion: 0,
		Token:           0x1234,
		Identifier:      PushData,
		GatewayEUI:      protocol.EUIFromInt64(0x11AABBBBCCCCDDDD),
		JSONString:      "ABCDEF",
	}

	buf, err := pkt.MarshalBinary()
	assert.NoError(err)
	assert.NotNil(buf)

	pkt = GwPacket{
		Identifier: PullAck,
	}
	_, err = pkt.MarshalBinary()
	assert.NoError(err, "Got error marshaling PULL_ACK")

	pkt = GwPacket{
		Identifier: PushAck,
	}
	_, err = pkt.MarshalBinary()
	assert.NoError(err, "Got error marshaling PULL_ACK")

	pkt = GwPacket{
		Identifier: PullData,
	}
	_, err = pkt.MarshalBinary()
	assert.NoError(err, "Got error marshaling PULL_DATA")

	pkt = GwPacket{
		Identifier: PullResp,
	}
	_, err = pkt.MarshalBinary()
	assert.NoError(err, "Got error marshaling PULL_RESP")

	pkt = GwPacket{
		Identifier: PullResp,
	}
	_, err = pkt.MarshalBinary()
	assert.NoError(err, "Got error marshaling PULL_RESP")

	pkt = GwPacket{
		Identifier: TxAck,
	}
	_, err = pkt.MarshalBinary()
	assert.NoError(err, "Got error marshaling PULL_DATA")

	pkt = GwPacket{
		Identifier: UnknownType,
	}
	_, err = pkt.MarshalBinary()
	assert.Error(err, "Expected error marshaling unknown type")

}

func TestBinaryMarshalUnmarshal(t *testing.T) {
	pkt := GwPacket{
		ProtocolVersion: 1,
		Token:           0xAABB,
		Identifier:      PullResp,
		JSONString: `{"txpk":{"imme":false,"tmst":254014692,"freq":868.5,"rfch":0,"modu":"LORA","datr":"SF12BW125","size":17,"data":"oGL34y/vUiWG
+OYcwPZAKgA"}}`,
	}

	bytes, err := pkt.MarshalBinary()
	if err != nil {
		t.Fatalf("Couldn't marshal packet (%v): %v ", pkt, err)
	}

	pk2 := GwPacket{}
	if err := pk2.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Got error unmarshaling bytes (source=%v, bytes=%v):  %v", pkt, bytes, err)
	}

	if pkt != pk2 {
		t.Fatalf("Not the same packet (original: %v != copy: %v)", pkt, pk2)
	}
}
