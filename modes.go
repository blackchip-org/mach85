package mach85

type Mode int

const (
	Absolute Mode = iota
	AbsoluteX
	AbsoluteY
	Accumulator
	Immediate
	Implied
	Indirect
	IndirectX
	IndirectY
	Relative
	ZeroPage
	ZeroPageX
	ZeroPageY
)

var operandLengths = map[Mode]int{
	Absolute:    2,
	AbsoluteX:   2,
	AbsoluteY:   2,
	Accumulator: 0,
	Immediate:   1,
	Implied:     0,
	Indirect:    2,
	IndirectX:   1,
	IndirectY:   1,
	Relative:    1,
	ZeroPage:    1,
	ZeroPageX:   1,
	ZeroPageY:   1,
}

var operandFormats = map[Mode]string{
	Absolute:    "$%04x",
	AbsoluteX:   "$%04x,x",
	AbsoluteY:   "$%04x,y",
	Accumulator: "a",
	Immediate:   "#$%02x",
	Indirect:    "($%04x)",
	IndirectX:   "($%02x,x)",
	IndirectY:   "($%02x),y",
	Relative:    "$%04x",
	ZeroPage:    "$%02x",
	ZeroPageX:   "$%02x,x",
	ZeroPageY:   "$%02x,y",
}
