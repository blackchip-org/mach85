package mach85

import (
	"bytes"
	"fmt"
)

type MemoryChunk interface {
	Load(address uint16) uint8
	Store(address uint16, value uint8)
}

type Memory struct {
	Base MemoryChunk
}

func NewMemory(base MemoryChunk) *Memory {
	return &Memory{Base: base}
}

func (m *Memory) Load(address uint16) uint8 {
	return m.Base.Load(address)
}

func (m *Memory) Store(address uint16, value uint8) {
	m.Base.Store(address, value)
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

func (m *Memory) Import(address uint16, data []uint8) {
	for i, value := range data {
		m.Store(address+uint16(i), value)
	}
}

func (m *Memory) Dump(start uint16, end uint16, decode Decoder) string {
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
			ch, printable := decode(value)
			if printable {
				chars.WriteString(fmt.Sprintf("%c", ch))
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

type RAM struct {
	bytes []uint8
}

func NewRAM(size int) *RAM {
	r := &RAM{
		bytes: make([]uint8, size, size),
	}
	return r
}

func (r *RAM) Load(address uint16) uint8 {
	return r.bytes[address]
}

func (r *RAM) Store(address uint16, value uint8) {
	r.bytes[address] = value
}

type ROM struct {
	bytes []uint8
}

func NewROM(bytes []uint8) *ROM {
	return &ROM{
		bytes: bytes,
	}
}

func (r *ROM) Load(address uint16) uint8 {
	return r.bytes[address]
}

func (r *ROM) Store(address uint16, value uint8) {
}

type NullMemory struct {
	Value uint8
}

func (m NullMemory) Load(address uint16) uint8 {
	return m.Value
}

func (m NullMemory) Store(address uint16, value uint8) {
}
