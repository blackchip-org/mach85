package mach85

import (
	"fmt"
)

// CPU is the MOS Technology 6502 processor.
type CPU struct {
	PC uint16 // Program counter
	A  uint8  // Accumulator
	X  uint8  // X register
	Y  uint8  // Y register
	SP uint8  // Stack pointer

	C bool // Carry flag
	Z bool // Zero flag
	I bool // Interrupt disable flag
	D bool // Decimal mode flag
	B bool // Break flag
	V bool // Overflow flag
	N bool // Signed flag

	mem *Memory
}

func NewCPU(mem *Memory) *CPU {
	return &CPU{
		mem: mem,
	}
}

func boolBit(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

func (c *CPU) SR() uint8 {
	return boolBit(c.C) |
		boolBit(c.Z)<<1 |
		boolBit(c.I)<<2 |
		boolBit(c.D)<<3 |
		boolBit(c.B)<<4 |
		1<<5 |
		boolBit(c.V)<<6 |
		boolBit(c.N)<<7
}

func boolChar(v bool) string {
	if v {
		return "*"
	}
	return "."
}

func (c *CPU) String() string {
	return fmt.Sprintf(""+
		" pc  sr ac xr yr sp  n v - b d i z c\n"+
		"%04x %02x %02x %02x %02x %02x  %s %s %s %s %s %s %s %s",
		c.PC, c.SR(), c.A, c.X, c.Y, c.SP,
		boolChar(c.N), boolChar(c.V), boolChar(true), boolChar(c.B),
		boolChar(c.D), boolChar(c.I), boolChar(c.Z), boolChar(c.C),
	)
}

func (c *CPU) fetch() uint8 {
	c.PC++
	return c.mem.Load(c.PC)
}

func (c *CPU) fetch16() uint16 {
	return uint16(c.fetch()) + (uint16(c.fetch()))<<8
}

func (c *CPU) Next() {
	opcode := c.fetch()
	exec, ok := opcodes[opcode]
	if ok {
		exec(c)
	}
}

func (c *CPU) Run() {
	for !c.B {
		c.Next()
	}
}

func (c *CPU) setFlagsNZ(value uint8) {
	c.Z = value == 0
	c.N = value&(1<<7) != 0
}

type loader func() (uint8, uint16)

func (c *CPU) loadAbsolute() (uint8, uint16) {
	address := c.fetch16()
	value := c.mem.Load(address)
	return value, address
}

func (c *CPU) loadAbsoluteX() (uint8, uint16) {
	address := c.fetch16() + uint16(c.X)
	value := c.mem.Load(address)
	return value, address
}

func (c *CPU) loadAbsoluteY() (uint8, uint16) {
	address := c.fetch16() + uint16(c.Y)
	value := c.mem.Load(address)
	return value, address
}

func (c *CPU) loadImmediate() (uint8, uint16) {
	address := uint16(0)
	value := c.fetch()
	return value, address
}

func (c *CPU) loadIndirectX() (uint8, uint16) {
	address := c.mem.Load16(uint16(c.fetch()) + uint16(c.X))
	value := c.mem.Load(address)
	return value, address
}

func (c *CPU) loadIndirectY() (uint8, uint16) {
	address := c.mem.Load16(uint16(c.fetch())) + uint16(c.Y)
	value := c.mem.Load(address)
	return value, address
}

func (c *CPU) loadZeroPage() (uint8, uint16) {
	address := uint16(c.fetch())
	value := c.mem.Load(address)
	return value, address
}

func (c *CPU) loadZeroPageX() (uint8, uint16) {
	address := uint16(c.fetch() + c.X)
	value := c.mem.Load(address)
	return value, address
}

func (c *CPU) loadZeroPageY() (uint8, uint16) {
	address := uint16(c.fetch() + c.Y)
	value := c.mem.Load(address)
	return value, address
}

type storer func(uint8)

func (c *CPU) storeAbsolute(value uint8) {
	address := c.fetch16()
	c.mem.Store(address, value)
}

func (c *CPU) storeAbsoluteX(value uint8) {
	address := c.fetch16() + uint16(c.X)
	c.mem.Store(address, value)
}

func (c *CPU) storeAbsoluteY(value uint8) {
	address := c.fetch16() + uint16(c.Y)
	c.mem.Store(address, value)
}

func (c *CPU) storeIndirectX(value uint8) {
	address := c.mem.Load16(uint16(c.fetch()) + uint16(c.X))
	c.mem.Store(address, value)
}

func (c *CPU) storeIndirectY(value uint8) {
	address := c.mem.Load16(uint16(c.fetch())) + uint16(c.Y)
	c.mem.Store(address, value)
}

func (c *CPU) storeZeroPage(value uint8) {
	address := uint16(c.fetch())
	c.mem.Store(address, value)
}

func (c *CPU) storeZeroPageX(value uint8) {
	address := uint16(c.fetch()) + uint16(c.X)
	c.mem.Store(address, value)
}

func (c *CPU) storeZeroPageY(value uint8) {
	address := uint16(c.fetch()) + uint16(c.Y)
	c.mem.Store(address, value)
}
