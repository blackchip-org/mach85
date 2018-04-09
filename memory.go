package mach85

type Memory struct {
	bytes []uint8
}

func NewMemory(size int) *Memory {
	m := &Memory{
		bytes: make([]uint8, size, size),
	}
	return m
}

func (m *Memory) Load(address uint16) uint8 {
	return m.bytes[address]
}

func (m *Memory) Store(address uint16, value uint8) {
	m.bytes[address] = value
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
