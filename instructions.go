package mach85

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

func sta(c *CPU, store storer) {
	store(c.A)
}

func stx(c *CPU, store storer) {
	store(c.X)
}

func sty(c *CPU, store storer) {
	store(c.Y)
}
