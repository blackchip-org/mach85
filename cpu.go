package mach85

import (
	"fmt"
	"log"
)

const (
	Stack       = uint16(0x0100)
	ResetVector = uint16(0xfffc)
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

	Trace bool
	mem   *Memory
	dasm  *Disassembler
}

func NewCPU(mem *Memory) *CPU {
	return &CPU{
		mem:  mem,
		dasm: NewDisassembler(mem),
	}
}

func (c *CPU) SR() uint8 {
	boolBit := func(v bool) uint8 {
		if v {
			return 1
		}
		return 0
	}
	return boolBit(c.C) |
		boolBit(c.Z)<<1 |
		boolBit(c.I)<<2 |
		boolBit(c.D)<<3 |
		boolBit(c.B)<<4 |
		1<<5 |
		boolBit(c.V)<<6 |
		boolBit(c.N)<<7
}

func (c *CPU) SetSR(v uint8) {
	c.C = v&(1<<0) != 0
	c.Z = v&(1<<1) != 0
	c.I = v&(1<<2) != 0
	c.D = v&(1<<3) != 0
	c.B = v&(1<<4) != 0
	c.V = v&(1<<6) != 0
	c.N = v&(1<<7) != 0
}

func (c *CPU) String() string {
	boolChar := func(v bool) string {
		if v {
			return "*"
		}
		return "."
	}
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

func (c *CPU) push(value uint8) {
	c.mem.Store(Stack+uint16(c.SP), value)
	c.SP--
}

func (c *CPU) push16(value uint16) {
	c.push(uint8(value >> 8))
	c.push(uint8(value & 0xff))
}

func (c *CPU) pull() uint8 {
	c.SP++
	return c.mem.Load(Stack + uint16(c.SP))
}

func (c *CPU) pull16() uint16 {
	return uint16(c.pull()) | uint16(c.pull())<<8
}

func (c *CPU) Reset() {
	// Vector is actual start address so set the PC one byte behind
	c.PC = c.mem.Load16(ResetVector) - 1
}

func (c *CPU) Next() {
	if c.Trace {
		c.dasm.PC = c.PC
		op := c.dasm.Next()
		fmt.Println(op)
	}
	opcode := c.fetch()
	execute, ok := executors[opcode]
	if !ok {
		log.Printf("$%04x: illegal opcode: $%02x", c.PC, opcode)
	} else {
		execute(c)
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

type loader func() (uint8, storer)
type storer func(uint8)

func (c *CPU) loadAbsolute() (uint8, storer) {
	address := c.fetch16()
	value := c.mem.Load(address)
	return value, func(v uint8) { c.mem.Store(address, v) }
}

func (c *CPU) loadAbsoluteX() (uint8, storer) {
	address := c.fetch16() + uint16(c.X)
	value := c.mem.Load(address)
	return value, func(v uint8) { c.mem.Store(address, v) }
}

func (c *CPU) loadAbsoluteY() (uint8, storer) {
	address := c.fetch16() + uint16(c.Y)
	value := c.mem.Load(address)
	return value, func(v uint8) { c.mem.Store(address, v) }
}

func (c *CPU) loadAccumulator() (uint8, storer) {
	value := c.A
	return value, func(v uint8) { c.A = v }
}

func (c *CPU) loadImmediate() (uint8, storer) {
	value := c.fetch()
	return value, nil
}

func (c *CPU) loadIndirectX() (uint8, storer) {
	address := c.mem.Load16(uint16(c.fetch()) + uint16(c.X))
	value := c.mem.Load(address)
	return value, func(v uint8) { c.mem.Store(address, v) }
}

func (c *CPU) loadIndirectY() (uint8, storer) {
	address := c.mem.Load16(uint16(c.fetch())) + uint16(c.Y)
	value := c.mem.Load(address)
	return value, func(v uint8) { c.mem.Store(address, v) }
}

func (c *CPU) loadZeroPage() (uint8, storer) {
	address := uint16(c.fetch())
	value := c.mem.Load(address)
	return value, func(v uint8) { c.mem.Store(address, v) }
}

func (c *CPU) loadZeroPageX() (uint8, storer) {
	address := uint16(c.fetch() + c.X)
	value := c.mem.Load(address)
	return value, func(v uint8) { c.mem.Store(address, v) }
}

func (c *CPU) loadZeroPageY() (uint8, storer) {
	address := uint16(c.fetch() + c.Y)
	value := c.mem.Load(address)
	return value, func(v uint8) { c.mem.Store(address, v) }
}

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
