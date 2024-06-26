package protocol

import (
	"encoding/base64"
	"testing"
)

func TestJoinRequestPayloadDecode(t *testing.T) {
	rawData := "AAgHBgUEAwIBvrrvvr66775+xPtJBQI="
	buffer, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		t.Fatal("Couldn't decode string: ", err)
	}

	pos := 1

	joinRequest := JoinRequestPayload{}
	if err := joinRequest.decode(buffer, &pos); err != nil {
		t.Fatal("Couldn't decode buffer")
	}

	appEUI, _ := EUIFromString("01-02-03-04-05-06-07-08")
	devEUI, _ := EUIFromString("BE-EF-BA-BE-BE-EF-BA-BE")

	if joinRequest.AppEUI != appEUI {
		t.Fatalf("Expected %v for app EUI but got %v", appEUI.String(), joinRequest.AppEUI.String())
	}
	if joinRequest.DevEUI != devEUI {
		t.Fatalf("Expected %v for dev EUI but got %v", devEUI.String(), joinRequest.DevEUI.String())
	}
}

func TestEncodeDecodeJoinRequest(t *testing.T) {
	jr1 := JoinRequestPayload{
		AppEUI:   EUIFromInt64(0x0102030405060708),
		DevEUI:   EUIFromInt64(0x0807060504030201),
		DevNonce: 0x1234,
	}

	buf := make([]byte, 50)
	pos := 0
	if err := jr1.encode(buf, &pos); err != nil {
		t.Fatalf("Couldn't encode: %v", err)
	}
	pos = 0
	jr2 := JoinRequestPayload{}
	if err := jr2.decode(buf, &pos); err != nil {
		t.Fatalf("Couldn't decode: %v", err)
	}

	if jr1 != jr2 {
		t.Fatalf("Encoded and decoded aren't equal: %+v != %+v", jr1, jr2)
	}
}
func TestJoinRequestPayloadDecodeInvalidBuffer(t *testing.T) {

	var buffer [1]byte
	pos := 0
	joinRequest := JoinRequestPayload{}
	if err := joinRequest.decode(buffer[:], &pos); err == nil {
		t.Fatal("Expected error when using invalid buffer but got none")
	}
}

func TestBasicEncodeDecode(t *testing.T) {
	basicEncoderTests(t, &JoinRequestPayload{})
	basicDecoderTests(t, &JoinRequestPayload{})
}
