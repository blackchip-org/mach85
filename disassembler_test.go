package mach85

import "testing"

func TestOperationString3(t *testing.T) {
	operation := Operation{
		Address:     0x1234,
		Instruction: Lda,
		Mode:        Absolute,
		Operand:     uint16(0x5678),
		Bytes:       []uint8{0xad, 0x78, 0x56},
	}
	want := "$1234: ad 78 56 lda $5678"
	have := operation.String()
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestOperationString2(t *testing.T) {
	operation := Operation{
		Address:     0x1234,
		Instruction: Lda,
		Mode:        ZeroPage,
		Operand:     uint16(0x56),
		Bytes:       []uint8{0xa5, 0x56},
	}
	want := "$1234: a5 56    lda $56"
	have := operation.String()
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestOperationString1(t *testing.T) {
	operation := Operation{
		Address:     0x1234,
		Instruction: Rts,
		Mode:        Implied,
		Operand:     uint16(0),
		Bytes:       []uint8{0x60},
	}
	want := "$1234: 60       rts"
	have := operation.String()
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}
