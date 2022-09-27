package protocol

import "testing"

func TestEuiFromString(t *testing.T) {
	euiStr := "01-02-03-04-05-06-07-08"

	eui, err := EUIFromString(euiStr)
	if err != nil {
		t.Error("Couldn't create EUI from string")
	}

	if eui.String() != euiStr {
		t.Error("Did not get the expected EUI string")
	}
}

func TestEuiFromInvalidString(t *testing.T) {
	_, err := EUIFromString("")
	if err == nil {
		t.Error("Expected error on empty string")
	}

	_, err = EUIFromString("foof")
	if err == nil {
		t.Error("Expected error on invalid string")
	}

	_, err = EUIFromString("01-02-03-04")
	if err == nil {
		t.Error("Expected error on too short string")
	}

	_, err = EUIFromString("01-02-03-04-05-06-07-08-01-02-03-04-05-06-07-08")
	if err == nil {
		t.Error("Expected error on too long string")
	}
}

func TestToFromUint64(t *testing.T) {
	eui1, _ := EUIFromString("01-02-03-04-05-06-07-08")
	eui2 := EUIFromInt64(0x0102030405060708)
	if eui1 != eui2 {
		t.Fatal("Not what I'd expect")
	}
	eui3 := EUIFromInt64(eui2.ToInt64())
	if eui2 != eui3 {
		t.Fatalf("Not what I'd expect (%+v != %+v)", eui2, eui3)
	}
}
