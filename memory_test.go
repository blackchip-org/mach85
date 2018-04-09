package mach85

import "testing"

func TestLoad16(t *testing.T) {
	m := NewMemory(0x10000)
	m.Store(0x00, 0xcd)
	m.Store(0x01, 0xab)
	want := uint16(0xabcd)
	have := m.Load16(0)
	if want != have {
		t.Errorf("\n want: %x \n have: %x \n", want, have)
	}
}

func TestStore16(t *testing.T) {
	m := NewMemory(0x10000)
	m.Store16(0x00, 0xabcd)
	want := uint8(0xcd)
	have := m.Load(0)
	if want != have {
		t.Errorf("\n want: %x \n have: %x \n", want, have)
	}
	want = uint8(0xab)
	have = m.Load(1)
	if want != have {
		t.Errorf("\n want: %x \n have: %x \n", want, have)
	}
}
