package pwaciii

import (
	"bytes"
	"testing"
)

func TestExtractData(t *testing.T) {
	apdu := []byte{0x10, 0x02, 0x41, 0x10, 0x10, 0x50, 0x10, 0x10, 0x51, 0x10, 0x10, 0x10, 0x10, 0x10, 0x03}
	expected := []byte{0x41, 0x10, 0x50, 0x10, 0x51, 0x10, 0x10}
	result := extractData(apdu)
	if !bytes.Equal(result, expected) {
		t.Errorf("Failed to extract data correctly. Got: ([% X]) [% X], Expected: [% X]", apdu, result, expected)
	}
}

func Test_Formatapdu(t *testing.T) {
	apdu := []byte{0x10, 0x02, 0x41, 0x50, 0x10, 0x10, 0x51, 0x50, 0x10, 0x03}
	data := []byte{0x41, 0x50, 0x10, 0x51}
	result := formatapdu(data)
	if !bytes.Equal(result, apdu) {
		t.Errorf("Failed to extract data correctly. Got: [% X], Expected: [% X]", result, apdu)
	}
}
