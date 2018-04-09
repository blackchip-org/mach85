package mach85

func and(c *CPU, load loader) {
	value, _ := load()
	c.A = c.A & value
	c.setFlagsNZ(c.A)
}

func brk(c *CPU) {
	c.B = true
	c.fetch()
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
