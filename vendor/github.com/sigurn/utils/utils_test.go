package utils

import (
	"testing"
)

func TestReverseByte(t *testing.T) {
	rb := ReverseByte(0x55)
	if rb != 0xAA {
		t.Errorf("Reverse byte failed. expected: 0xAA, real: 0x%X", rb)
	}
}

func TestReverseUint16(t *testing.T) {
	rval := ReverseUint16(0x5555)
	if rval != 0xAAAA {
		t.Errorf("Reverse uint16 failed. expected: 0xAAAA, real: 0x%X", rval)
	}
}
