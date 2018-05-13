package mach85

func adc(c *CPU, load loader) {
	v1 := c.A
	v2, _ := load()
	if c.D {
		v1 = fromBCD(v1)
		v2 = fromBCD(v2)
	}
	utotal := uint16(v1) + uint16(v2)
	total := int16(int8(v1)) + int16(int8(v2))
	if c.C {
		utotal++
		total++
	}
	if c.D {
		c.C = utotal > 99
		c.A = toBCD(uint8(utotal))
	} else {
		c.C = utotal > 0xff
		c.V = total < -128 || total > 127
		c.A = uint8(utotal)
	}
	c.setFlagsNZ(c.A)
}

func and(c *CPU, load loader) {
	value, _ := load()
	c.A = c.A & value
	c.setFlagsNZ(c.A)
}

func asl(c *CPU, load loader) {
	value, store := load()
	c.C = (value & 0x80) != 0
	value = (value << 1)
	c.setFlagsNZ(value)
	store(value)
}

func bit(c *CPU, load loader) {
	value, _ := load()
	c.Z = (c.A & value) == 0
	c.N = (value & (1 << 7)) != 0
	c.V = (value & (1 << 6)) != 0
}

func branch(c *CPU, do bool) {
	displacement := int8(c.fetch())
	if do {
		if displacement >= 0 {
			c.PC += uint16(displacement)
		} else {
			c.PC -= uint16(displacement * -1)
		}
	}
}

func brk(c *CPU) {
	c.B = true
	c.fetch()
	if c.StopOnBreak {
		return
	}
	c.push16(c.PC + 1)
	c.push(c.SR())
	c.I = true
	c.PC = c.mem.Load16(AddrIrqVector) - 1
}

func cmp(c *CPU, register uint8, load loader) {
	value, _ := load()
	// C set as if subtraction. Clear if 'borrow', otherwise set
	result := int16(register) - int16(value)
	c.C = result >= 0
	c.setFlagsNZ(uint8(result))
}

func dec(c *CPU, load loader) {
	value, store := load()
	value = value - 1
	c.setFlagsNZ(value)
	store(value)
}

func dex(c *CPU) {
	c.X--
	c.setFlagsNZ(c.X)
}

func dey(c *CPU) {
	c.Y--
	c.setFlagsNZ(c.Y)
}

func eor(c *CPU, load loader) {
	value, _ := load()
	c.A = c.A ^ value
	c.setFlagsNZ(c.A)
}

func inc(c *CPU, load loader) {
	value, store := load()
	value = value + 1
	c.setFlagsNZ(value)
	store(value)
}

func inx(c *CPU) {
	c.X++
	c.setFlagsNZ(c.X)
}

func iny(c *CPU) {
	c.Y++
	c.setFlagsNZ(c.Y)
}

func jmp(c *CPU) {
	c.PC = c.fetch16() - 1
}

func jmpIndirect(c *CPU) {
	c.PC = c.mem.Load16(c.fetch16()) - 1
}

func jsr(c *CPU) {
	address := c.fetch16()
	c.push16(c.PC)
	c.PC = address - 1
}

func lda(c *CPU, load loader) {
	value, _ := load()
	c.A = value
	c.setFlagsNZ(value)
}

func ldx(c *CPU, load loader) {
	value, _ := load()
	c.X = value
	c.setFlagsNZ(value)
}

func ldy(c *CPU, load loader) {
	value, _ := load()
	c.Y = value
	c.setFlagsNZ(value)
}

func lsr(c *CPU, load loader) {
	value, store := load()
	c.C = (value & 0x01) != 0
	value = value >> 1
	c.setFlagsNZ(value)
	store(value)
}

func ora(c *CPU, load loader) {
	value, _ := load()
	c.A = c.A | value
	c.setFlagsNZ(c.A)
}

func php(c *CPU) {
	// https://wiki.nesdev.com/w/index.php/Status_flags
	c.push(c.SR() | flagB)
}

func pla(c *CPU) {
	c.A = c.pull()
	c.setFlagsNZ(c.A)
}

func rol(c *CPU, load loader) {
	value, store := load()
	rotate := uint8(0)
	if c.C {
		rotate = 1
	}
	c.C = value&0x80 != 0
	value = value<<1 + rotate
	c.setFlagsNZ(value)
	store(value)
}

func ror(c *CPU, load loader) {
	value, store := load()
	rotate := uint8(0)
	if c.C {
		rotate = 0x80
	}
	c.C = value&0x01 != 0
	value = value>>1 + rotate
	c.setFlagsNZ(value)
	store(value)
}

func rti(c *CPU) {
	c.SetSR(c.pull())
	c.PC = c.pull16() - 1
}

func sbc(c *CPU, load loader) {
	v1 := c.A
	v2, _ := load()
	if c.D {
		v1 = fromBCD(v1)
		v2 = fromBCD(v2)
	}
	utotal := v1 - v2
	total := int16(int8(v1)) - int16(int8(v2))
	borrow := false
	if v1 < v2 {
		borrow = true
	}
	// borrow if carry is clear
	if !c.C {
		if total == 0 {
			borrow = true
		}
		utotal--
		total--
	}
	c.C = !borrow
	if c.D {
		if total < 0 {
			total += 100
		}
		c.A = toBCD(uint8(total))
	} else {
		c.V = total < -128 || total > 127
		c.A = uint8(utotal)
	}
	c.setFlagsNZ(c.A)
}

func sta(c *CPU, store storer) {
	store(c.A)
}

func stx(c *CPU, store storer) {
	store(c.X)
}

func sty(c *CPU, store storer) {
	store(c.Y)
}

func transfer(c *CPU, from uint8, to *uint8) {
	*to = from
	c.setFlagsNZ(*to)
}
