package mach85

var opcodes = map[uint8]func(c *CPU){
	0x00: func(c *CPU) { brk(c) },
	0x01: func(c *CPU) { ora(c, c.loadIndirectX) },
	0x05: func(c *CPU) { ora(c, c.loadZeroPage) },
	0x06: func(c *CPU) { asl(c, c.loadZeroPage) },
	0x09: func(c *CPU) { ora(c, c.loadImmediate) },
	0x0a: func(c *CPU) { asl(c, c.loadAccumulator) },
	0x0d: func(c *CPU) { ora(c, c.loadAbsolute) },
	0x0e: func(c *CPU) { asl(c, c.loadAbsolute) },

	0x10: func(c *CPU) { branch(c, !c.N) }, // bpl
	0x11: func(c *CPU) { ora(c, c.loadIndirectY) },
	0x15: func(c *CPU) { ora(c, c.loadZeroPageX) },
	0x16: func(c *CPU) { asl(c, c.loadZeroPageX) },
	0x18: func(c *CPU) { c.C = false }, // clc
	0x19: func(c *CPU) { ora(c, c.loadAbsoluteY) },
	0x1d: func(c *CPU) { ora(c, c.loadAbsoluteX) },
	0x1e: func(c *CPU) { asl(c, c.loadAbsoluteX) },

	0x20: func(c *CPU) { jsr(c) },
	0x21: func(c *CPU) { and(c, c.loadIndirectX) },
	0x24: func(c *CPU) { bit(c, c.loadZeroPage) },
	0x25: func(c *CPU) { and(c, c.loadZeroPage) },
	0x29: func(c *CPU) { and(c, c.loadImmediate) },
	0x2c: func(c *CPU) { bit(c, c.loadAbsolute) },
	0x2d: func(c *CPU) { and(c, c.loadAbsolute) },

	0x30: func(c *CPU) { branch(c, c.N) }, // bmi
	0x31: func(c *CPU) { and(c, c.loadIndirectY) },
	0x35: func(c *CPU) { and(c, c.loadZeroPageX) },
	0x38: func(c *CPU) { c.C = true }, // sec
	0x39: func(c *CPU) { and(c, c.loadAbsoluteY) },
	0x3d: func(c *CPU) { and(c, c.loadAbsoluteX) },

	0x41: func(c *CPU) { eor(c, c.loadIndirectX) },
	0x45: func(c *CPU) { eor(c, c.loadZeroPage) },
	0x46: func(c *CPU) { lsr(c, c.loadZeroPage) },
	0x49: func(c *CPU) { eor(c, c.loadImmediate) },
	0x4a: func(c *CPU) { lsr(c, c.loadAccumulator) },
	0x4c: func(c *CPU) { jmp(c) },
	0x4d: func(c *CPU) { eor(c, c.loadAbsolute) },
	0x4e: func(c *CPU) { lsr(c, c.loadAbsolute) },

	0x50: func(c *CPU) { branch(c, !c.V) }, // bvc
	0x51: func(c *CPU) { eor(c, c.loadIndirectY) },
	0x55: func(c *CPU) { eor(c, c.loadZeroPageX) },
	0x56: func(c *CPU) { lsr(c, c.loadZeroPageX) },
	0x58: func(c *CPU) { c.I = false }, // cli
	0x59: func(c *CPU) { eor(c, c.loadAbsoluteY) },
	0x5d: func(c *CPU) { eor(c, c.loadAbsoluteX) },
	0x5e: func(c *CPU) { lsr(c, c.loadAbsoluteX) },

	0x6c: func(c *CPU) { jmpIndirect(c) },

	0x70: func(c *CPU) { branch(c, c.V) }, // bvs
	0x78: func(c *CPU) { c.I = true },     // sei

	0x81: func(c *CPU) { sta(c, c.storeIndirectX) },
	0x84: func(c *CPU) { sty(c, c.storeZeroPage) },
	0x85: func(c *CPU) { sta(c, c.storeZeroPage) },
	0x86: func(c *CPU) { stx(c, c.storeZeroPage) },
	0x88: func(c *CPU) { dey(c) },
	0x8a: func(c *CPU) { transfer(c, c.X, &c.A) },
	0x8c: func(c *CPU) { sty(c, c.storeAbsolute) },
	0x8d: func(c *CPU) { sta(c, c.storeAbsolute) },
	0x8e: func(c *CPU) { stx(c, c.storeAbsolute) },

	0x90: func(c *CPU) { branch(c, !c.C) }, // bcc
	0x91: func(c *CPU) { sta(c, c.storeIndirectY) },
	0x94: func(c *CPU) { sty(c, c.storeZeroPageX) },
	0x95: func(c *CPU) { sta(c, c.storeZeroPageX) },
	0x96: func(c *CPU) { stx(c, c.storeZeroPageY) },
	0x98: func(c *CPU) { transfer(c, c.Y, &c.A) },
	0x99: func(c *CPU) { sta(c, c.storeAbsoluteY) },
	0x9d: func(c *CPU) { sta(c, c.storeAbsoluteX) },

	0xa0: func(c *CPU) { ldy(c, c.loadImmediate) },
	0xa1: func(c *CPU) { lda(c, c.loadIndirectX) },
	0xa2: func(c *CPU) { ldx(c, c.loadImmediate) },
	0xa4: func(c *CPU) { ldy(c, c.loadZeroPage) },
	0xa5: func(c *CPU) { lda(c, c.loadZeroPage) },
	0xa6: func(c *CPU) { ldx(c, c.loadZeroPage) },
	0xa8: func(c *CPU) { transfer(c, c.A, &c.Y) },
	0xa9: func(c *CPU) { lda(c, c.loadImmediate) },
	0xaa: func(c *CPU) { transfer(c, c.A, &c.X) },
	0xac: func(c *CPU) { ldy(c, c.loadAbsolute) },
	0xad: func(c *CPU) { lda(c, c.loadAbsolute) },
	0xae: func(c *CPU) { ldx(c, c.loadAbsolute) },

	0xb0: func(c *CPU) { branch(c, c.C) }, // bcs
	0xb1: func(c *CPU) { lda(c, c.loadIndirectY) },
	0xb4: func(c *CPU) { ldy(c, c.loadZeroPageX) },
	0xb5: func(c *CPU) { lda(c, c.loadZeroPageX) },
	0xb6: func(c *CPU) { ldx(c, c.loadZeroPageY) },
	0xb8: func(c *CPU) { c.V = false }, // clv
	0xb9: func(c *CPU) { lda(c, c.loadAbsoluteY) },
	0xbd: func(c *CPU) { lda(c, c.loadAbsoluteX) },
	0xbc: func(c *CPU) { ldy(c, c.loadAbsoluteX) },
	0xbe: func(c *CPU) { ldx(c, c.loadAbsoluteY) },

	0xc0: func(c *CPU) { cmp(c, c.Y, c.loadImmediate) },
	0xc1: func(c *CPU) { cmp(c, c.A, c.loadIndirectX) },
	0xc4: func(c *CPU) { cmp(c, c.Y, c.loadZeroPage) },
	0xc5: func(c *CPU) { cmp(c, c.A, c.loadZeroPage) },
	0xc6: func(c *CPU) { dec(c, c.loadZeroPage) },
	0xc8: func(c *CPU) { iny(c) },
	0xc9: func(c *CPU) { cmp(c, c.A, c.loadImmediate) },
	0xca: func(c *CPU) { dex(c) },
	0xcc: func(c *CPU) { cmp(c, c.Y, c.loadAbsolute) },
	0xcd: func(c *CPU) { cmp(c, c.A, c.loadAbsolute) },
	0xce: func(c *CPU) { dec(c, c.loadAbsolute) },

	0xd0: func(c *CPU) { branch(c, !c.Z) }, // bne
	0xd1: func(c *CPU) { cmp(c, c.A, c.loadIndirectY) },
	0xd5: func(c *CPU) { cmp(c, c.A, c.loadZeroPageX) },
	0xd6: func(c *CPU) { dec(c, c.loadZeroPageX) },
	0xd8: func(c *CPU) { c.D = false }, // cld
	0xd9: func(c *CPU) { cmp(c, c.A, c.loadAbsoluteY) },
	0xdd: func(c *CPU) { cmp(c, c.A, c.loadAbsoluteX) },
	0xde: func(c *CPU) { dec(c, c.loadAbsoluteX) },

	0xe0: func(c *CPU) { cmp(c, c.X, c.loadImmediate) },
	0xe4: func(c *CPU) { cmp(c, c.X, c.loadZeroPage) },
	0xe6: func(c *CPU) { inc(c, c.loadZeroPage) },
	0xe8: func(c *CPU) { inx(c) },
	0xea: func(c *CPU) {}, // nop
	0xec: func(c *CPU) { cmp(c, c.X, c.loadAbsolute) },
	0xee: func(c *CPU) { inc(c, c.loadAbsolute) },

	0xf0: func(c *CPU) { branch(c, c.Z) }, // beq
	0xf6: func(c *CPU) { inc(c, c.loadZeroPageX) },
	0xf8: func(c *CPU) { c.D = true }, // sed
	0xfe: func(c *CPU) { inc(c, c.loadAbsoluteX) },
}
