package mach85

import (
	"bytes"
	"fmt"
)

type Memory struct {
	bytes []uint8
}

func NewMemory(size int) *Memory {
	m := &Memory{
		bytes: make([]uint8, size, size),
	}
	return m
}

func NewMemory64k() *Memory {
	return NewMemory(0x10000)
}

func (m *Memory) Load(address uint16) uint8 {
	return m.bytes[address]
}

func (m *Memory) Store(address uint16, value uint8) {
	m.bytes[address] = value
}

func (m *Memory) StoreN(address uint16, values ...uint8) {
	for i, value := range values {
		m.Store(address+uint16(i), value)
	}
}

func (m *Memory) Load16(address uint16) uint16 {
	lo := uint16(m.Load(address))
	hi := uint16(m.Load(address + 1))
	return hi<<8 + lo
}

func (m *Memory) Store16(address uint16, value uint16) {
	m.Store(address, uint8(value%0x100))
	m.Store(address+1, uint8(value>>8))
}

func (m *Memory) Dump(start uint16, end uint16) string {
	var buf bytes.Buffer
	var chars bytes.Buffer

	a0 := start / 0x10 * 0x10
	a1 := end / 0x10 * 0x10
	if a1 != end {
		a1 += 0x10
	}
	for addr := a0; addr < a1; addr++ {
		if addr%0x10 == 0 {
			buf.WriteString(fmt.Sprintf("$%04x", addr))
			chars.Reset()
		}
		if addr < start || addr > end {
			buf.WriteString("   ")
			chars.WriteString(" ")
		} else {
			value := m.Load(addr)
			buf.WriteString(fmt.Sprintf(" %02x", value))
			if value >= 0x20 && value < 0x80 {
				chars.WriteString(fmt.Sprintf("%c", value))
			} else {
				chars.WriteString(".")
			}
		}
		if addr%0x10 == 7 {
			buf.WriteString(" ")
		}
		if addr%0x10 == 0x0f {
			buf.WriteString(" " + chars.String())
			if addr < end-1 {
				buf.WriteString("\n")
			}
		}
	}
	return buf.String()
}

func (m *Memory) Import(address uint16, data []uint8) {
	for i, value := range data {
		m.Store(address+uint16(i), value)
	}
}
