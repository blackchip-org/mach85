package mach85

import "testing"

func TestFromBCD(t *testing.T) {
	want := uint8(42)
	have := fromBCD(0x42)
	if want != have {
		t.Errorf("\n want: %v \n have: %v\n", want, have)
	}
}

func TestToBCD(t *testing.T) {
	want := uint8(0x42)
	have := toBCD(42)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}

func TestToBCDOverflow(t *testing.T) {
	want := uint8(0x12)
	have := toBCD(112)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}
func TestSigned8Plus1(t *testing.T) {
	want := int8(1)
	have := signed8(uint8(1))
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}

func TestSigned8Minus1(t *testing.T) {
	want := int8(-1)
	have := signed8(uint8(0xff))
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}
