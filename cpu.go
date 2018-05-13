package mach85

import (
	"fmt"
	"log"
)

// CPU is the MOS Technology 6502 series processor.
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

	StopOnBreak bool

	mem   *Memory
	irq   chan bool
	inISR bool
}

const (
	flagC = uint8(1 << 0)
	flagZ = uint8(1 << 1)
	flagI = uint8(1 << 2)
	flagD = uint8(1 << 3)
	flagB = uint8(1 << 4)
	flag5 = uint8(1 << 5)
	flagV = uint8(1 << 6)
	flagN = uint8(1 << 7)
)

func New6510(mem *Memory) *CPU {
	return &CPU{
		mem: mem,
		irq: make(chan bool, 10),
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
	c.mem.Store(AddrStack+uint16(c.SP), value)
	c.SP--
}

func (c *CPU) push16(value uint16) {
	c.push(uint8(value >> 8))
	c.push(uint8(value & 0xff))
}

func (c *CPU) pull() uint8 {
	c.SP++
	return c.mem.Load(AddrStack + uint16(c.SP))
}

func (c *CPU) pull16() uint16 {
	return uint16(c.pull()) | uint16(c.pull())<<8
}

func (c *CPU) Reset() {
	// Vector is actual start address so set the PC one byte behind
	c.PC = c.mem.Load16(AddrResetVector) - 1
}

func (c *CPU) Next() {
	opcode := c.fetch()
	execute, ok := executors[opcode]
	if !ok {
		log.Printf("$%04x: illegal opcode: $%02x", c.PC, opcode)
	} else {
		execute(c)
	}
	if opcode == 0x40 { // rti
		c.inISR = false
	}
	select {
	case <-c.irq:
		if !c.I {
			// http://www.6502.org/tutorials/6502opcodes.html#RTI
			// Note that unlike RTS, the return address on the stack is the
			// actual address rather than the address-1.
			c.push16(c.PC + 1)
			c.push(c.SR())
			c.PC = c.mem.Load16(AddrIrqVector) - 1
			c.inISR = true
		}
	default:
	}
}

func (c *CPU) IRQ() {
	c.irq <- true
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
	address := c.fetch()
	value := c.mem.Load(uint16(address))
	return value, func(v uint8) { c.mem.Store(uint16(address), v) }
}

func (c *CPU) loadZeroPageX() (uint8, storer) {
	address := c.fetch() + c.X
	value := c.mem.Load(uint16(address))
	return value, func(v uint8) { c.mem.Store(uint16(address), v) }
}

func (c *CPU) loadZeroPageY() (uint8, storer) {
	address := c.fetch() + c.Y
	value := c.mem.Load(uint16(address))
	return value, func(v uint8) { c.mem.Store(uint16(address), v) }
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
