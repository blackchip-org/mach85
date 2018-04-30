package mach85

type Instruction int

// Instructions
const (
	Illegal Instruction = iota // illegal instruction
	Adc                        // add with carry
	And                        // bitwise And with accumulator
	Asl                        // arithmetic shift left
	Bit                        // test bits
	Bmi                        // branch on minus
	Bcc                        // branch on carry clear
	Bcs                        // branch on carry set
	Beq                        // branch on equal
	Bne                        // branch on not equal
	Bpl                        // branch on plus
	Brk                        // break
	Bvc                        // branch on overflow clear
	Bvs                        // branch on overflow set
	Clc                        // clear carry
	Cld                        // clear decimal mode
	Cli                        // clear interrupt
	Clv                        // clear overflow
	Cmp                        // compare accumulator
	Cpx                        // compare x register
	Cpy                        // compare y register
	Dec                        // decrement memory
	Dex                        // decrement x
	Dey                        // decrement y
	Eor                        // bitwise exclusive or
	Inc                        // increment memory
	Inx                        // increment x
	Iny                        // increment y
	Jmp                        // jump
	Jsr                        // jump to subroutine
	Lda                        // load accumulator
	Ldx                        // load x register
	Ldy                        // load y register
	Lsr                        // logical shift right
	Nop                        // no operation
	Ora                        // bitwise or with accumulator
	Pha                        // push accumulator
	Php                        // push processor status
	Pla                        // pull accumulator
	Plp                        // pull processor status
	Rol                        // rotate left
	Ror                        // rotate right
	Rti                        // return from interrupt
	Rts                        // return from subroutine
	Sbc                        // subrtact with carry
	Sec                        // set carry
	Sed                        // set decimal mode
	Sei                        // set interrupt
	Sta                        // store accumulator
	Stx                        // store x register
	Sty                        // store y register
	Tax                        // transfer a to x
	Tay                        // transfer a to y
	Tsx                        // transfer stack pointer to x
	Txs                        // transfer x to stack pointer
	Txa                        // transfer x to a
	Tya                        // transfer y to a
)

var instructionStrings = map[Instruction]string{
	Illegal: "???",
	Adc:     "adc",
	And:     "and",
	Asl:     "asl",
	Bit:     "bit",
	Bmi:     "bmi",
	Bcc:     "bcc",
	Bcs:     "bcs",
	Beq:     "beq",
	Bne:     "bne",
	Bpl:     "bpl",
	Brk:     "brk",
	Bvc:     "bvc",
	Bvs:     "bvs",
	Clc:     "clc",
	Cld:     "cld",
	Cli:     "cli",
	Clv:     "clv",
	Cmp:     "cmp",
	Cpx:     "cpx",
	Cpy:     "cpy",
	Dec:     "dec",
	Dex:     "dex",
	Dey:     "dey",
	Eor:     "eor",
	Inc:     "inc",
	Inx:     "inx",
	Iny:     "iny",
	Jmp:     "jmp",
	Jsr:     "jsr",
	Lda:     "lda",
	Ldx:     "ldx",
	Ldy:     "ldy",
	Lsr:     "lsr",
	Nop:     "nop",
	Ora:     "ora",
	Pha:     "pha",
	Php:     "php",
	Pla:     "pla",
	Plp:     "plp",
	Rol:     "rol",
	Ror:     "ror",
	Rti:     "rti",
	Rts:     "rts",
	Sbc:     "sbc",
	Sec:     "sec",
	Sed:     "sed",
	Sei:     "sei",
	Sta:     "sta",
	Stx:     "stx",
	Sty:     "sty",
	Tax:     "tax",
	Tay:     "tay",
	Tsx:     "tsx",
	Txs:     "txs",
	Txa:     "txa",
	Tya:     "tya",
}

func (i Instruction) String() string {
	return instructionStrings[i]
}

type op struct {
	inst Instruction
	mode Mode
}

var opcodes = map[uint8]op{
	0x00: op{Brk, Implied},
	0x01: op{Ora, IndirectX},
	0x05: op{Ora, ZeroPage},
	0x06: op{Asl, ZeroPage},
	0x08: op{Php, Implied},
	0x09: op{Ora, Immediate},
	0x0a: op{Asl, Accumulator},
	0x0d: op{Ora, Absolute},
	0x0e: op{Asl, Absolute},

	0x10: op{Bpl, Relative},
	0x11: op{Ora, IndirectY},
	0x15: op{Ora, ZeroPageX},
	0x16: op{Asl, ZeroPageX},
	0x18: op{Clc, Implied},
	0x19: op{Ora, AbsoluteY},
	0x1d: op{Ora, AbsoluteX},
	0x1e: op{Asl, AbsoluteX},

	0x20: op{Jsr, Absolute},
	0x21: op{And, IndirectX},
	0x24: op{Bit, ZeroPage},
	0x25: op{And, ZeroPage},
	0x26: op{Rol, ZeroPage},
	0x28: op{Plp, Implied},
	0x29: op{And, Immediate},
	0x2a: op{Rol, Accumulator},
	0x2c: op{Bit, Absolute},
	0x2d: op{And, Absolute},
	0x2e: op{Rol, Absolute},

	0x30: op{Bmi, Relative},
	0x31: op{And, IndirectY},
	0x35: op{And, ZeroPageX},
	0x36: op{Rol, ZeroPageX},
	0x38: op{Sec, Implied},
	0x39: op{And, AbsoluteY},
	0x3d: op{And, AbsoluteX},
	0x3e: op{Rol, AbsoluteX},

	0x40: op{Rti, Implied},
	0x41: op{Eor, IndirectX},
	0x45: op{Eor, ZeroPage},
	0x46: op{Lsr, ZeroPage},
	0x48: op{Pha, Implied},
	0x49: op{Eor, Immediate},
	0x4a: op{Lsr, Accumulator},
	0x4c: op{Jmp, Absolute},
	0x4d: op{Eor, Absolute},
	0x4e: op{Lsr, Absolute},

	0x50: op{Bvc, Relative},
	0x51: op{Eor, IndirectY},
	0x55: op{Eor, ZeroPageX},
	0x56: op{Lsr, ZeroPageX},
	0x58: op{Cli, Implied},
	0x59: op{Eor, AbsoluteY},
	0x5d: op{Eor, AbsoluteX},
	0x5e: op{Lsr, AbsoluteX},

	0x60: op{Rts, Implied},
	0x61: op{Adc, IndirectX},
	0x65: op{Adc, ZeroPage},
	0x66: op{Ror, ZeroPage},
	0x68: op{Pla, Implied},
	0x69: op{Adc, Immediate},
	0x6a: op{Ror, Accumulator},
	0x6c: op{Jmp, Indirect},
	0x6d: op{Adc, Absolute},
	0x6e: op{Ror, Absolute},

	0x70: op{Bvs, Relative},
	0x71: op{Adc, IndirectY},
	0x75: op{Adc, ZeroPageX},
	0x76: op{Ror, ZeroPageX},
	0x78: op{Sei, Implied},
	0x79: op{Adc, AbsoluteY},
	0x7d: op{Adc, AbsoluteX},
	0x7e: op{Ror, AbsoluteX},

	0x81: op{Sta, IndirectX},
	0x84: op{Sty, ZeroPage},
	0x85: op{Sta, ZeroPage},
	0x86: op{Stx, ZeroPage},
	0x88: op{Dey, Implied},
	0x8a: op{Txa, Implied},
	0x8c: op{Sty, Absolute},
	0x8d: op{Sta, Absolute},
	0x8e: op{Stx, Absolute},

	0x90: op{Bcc, Relative},
	0x91: op{Sta, IndirectY},
	0x94: op{Sty, ZeroPageX},
	0x95: op{Sta, ZeroPageX},
	0x96: op{Stx, ZeroPageY},
	0x98: op{Tya, Implied},
	0x99: op{Sta, AbsoluteY},
	0x9a: op{Txs, Implied},
	0x9d: op{Sta, AbsoluteX},

	0xa0: op{Ldy, Immediate},
	0xa1: op{Lda, IndirectX},
	0xa2: op{Ldx, Immediate},
	0xa4: op{Ldy, ZeroPage},
	0xa5: op{Lda, ZeroPage},
	0xa6: op{Ldx, ZeroPage},
	0xa8: op{Tay, Implied},
	0xa9: op{Lda, Immediate},
	0xaa: op{Tax, Implied},
	0xac: op{Ldy, Absolute},
	0xad: op{Lda, Absolute},
	0xae: op{Ldx, Absolute},

	0xb0: op{Bcs, Relative},
	0xb1: op{Lda, IndirectY},
	0xb4: op{Ldy, ZeroPageX},
	0xb5: op{Lda, ZeroPageX},
	0xb6: op{Ldx, ZeroPageY},
	0xb8: op{Clv, Implied},
	0xb9: op{Lda, AbsoluteY},
	0xba: op{Tsx, Implied},
	0xbd: op{Lda, AbsoluteX},
	0xbc: op{Ldy, AbsoluteX},
	0xbe: op{Ldx, AbsoluteY},

	0xc0: op{Cpy, Immediate},
	0xc1: op{Cmp, IndirectX},
	0xc4: op{Cpy, ZeroPage},
	0xc5: op{Cmp, ZeroPage},
	0xc6: op{Dec, ZeroPage},
	0xc8: op{Iny, Implied},
	0xc9: op{Cmp, Immediate},
	0xca: op{Dex, Implied},
	0xcc: op{Cpy, Absolute},
	0xcd: op{Cmp, Absolute},
	0xce: op{Dec, Absolute},

	0xd0: op{Bne, Relative},
	0xd1: op{Cmp, IndirectY},
	0xd5: op{Cmp, ZeroPageX},
	0xd6: op{Dec, ZeroPageX},
	0xd8: op{Cld, Implied},
	0xd9: op{Cmp, AbsoluteY},
	0xdd: op{Cmp, AbsoluteX},
	0xde: op{Dec, AbsoluteX},

	0xe0: op{Cpx, Immediate},
	0xe1: op{Sbc, IndirectX},
	0xe4: op{Cpx, ZeroPage},
	0xe5: op{Sbc, ZeroPage},
	0xe6: op{Inc, ZeroPage},
	0xe8: op{Inx, Implied},
	0xe9: op{Sbc, Immediate},
	0xea: op{Nop, Implied},
	0xec: op{Cpx, Absolute},
	0xed: op{Sbc, Absolute},
	0xee: op{Inc, Absolute},

	0xf0: op{Beq, Relative},
	0xf1: op{Sbc, IndirectY},
	0xf5: op{Sbc, ZeroPageX},
	0xf6: op{Inc, ZeroPageX},
	0xf8: op{Sed, Implied},
	0xf9: op{Sbc, AbsoluteY},
	0xfd: op{Sbc, AbsoluteX},
	0xfe: op{Inc, AbsoluteX},
}

var executors = map[uint8]func(c *CPU){
	0x00: func(c *CPU) { brk(c) },
	0x01: func(c *CPU) { ora(c, c.loadIndirectX) },
	0x05: func(c *CPU) { ora(c, c.loadZeroPage) },
	0x06: func(c *CPU) { asl(c, c.loadZeroPage) },
	0x08: func(c *CPU) { c.push(c.SR()) }, // php
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
	0x26: func(c *CPU) { rol(c, c.loadZeroPage) },
	0x28: func(c *CPU) { c.SetSR(c.pull()) }, // plp
	0x29: func(c *CPU) { and(c, c.loadImmediate) },
	0x2a: func(c *CPU) { rol(c, c.loadAccumulator) },
	0x2c: func(c *CPU) { bit(c, c.loadAbsolute) },
	0x2d: func(c *CPU) { and(c, c.loadAbsolute) },
	0x2e: func(c *CPU) { rol(c, c.loadAbsolute) },

	0x30: func(c *CPU) { branch(c, c.N) }, // bmi
	0x31: func(c *CPU) { and(c, c.loadIndirectY) },
	0x35: func(c *CPU) { and(c, c.loadZeroPageX) },
	0x36: func(c *CPU) { rol(c, c.loadZeroPageX) },
	0x38: func(c *CPU) { c.C = true }, // sec
	0x39: func(c *CPU) { and(c, c.loadAbsoluteY) },
	0x3d: func(c *CPU) { and(c, c.loadAbsoluteX) },
	0x3e: func(c *CPU) { rol(c, c.loadAbsoluteX) },

	0x40: func(c *CPU) { rti(c) },
	0x41: func(c *CPU) { eor(c, c.loadIndirectX) },
	0x45: func(c *CPU) { eor(c, c.loadZeroPage) },
	0x46: func(c *CPU) { lsr(c, c.loadZeroPage) },
	0x48: func(c *CPU) { c.push(c.A) }, // pha
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

	0x60: func(c *CPU) { c.PC = c.pull16() }, // rts
	0x61: func(c *CPU) { adc(c, c.loadIndirectX) },
	0x65: func(c *CPU) { adc(c, c.loadZeroPage) },
	0x66: func(c *CPU) { ror(c, c.loadZeroPage) },
	0x68: func(c *CPU) { pla(c) },
	0x69: func(c *CPU) { adc(c, c.loadImmediate) },
	0x6a: func(c *CPU) { ror(c, c.loadAccumulator) },
	0x6c: func(c *CPU) { jmpIndirect(c) },
	0x6d: func(c *CPU) { adc(c, c.loadAbsolute) },
	0x6e: func(c *CPU) { ror(c, c.loadAbsolute) },

	0x70: func(c *CPU) { branch(c, c.V) }, // bvs
	0x71: func(c *CPU) { adc(c, c.loadIndirectY) },
	0x75: func(c *CPU) { adc(c, c.loadZeroPageX) },
	0x76: func(c *CPU) { ror(c, c.loadZeroPageX) },
	0x78: func(c *CPU) { c.I = true }, // sei
	0x79: func(c *CPU) { adc(c, c.loadAbsoluteY) },
	0x7d: func(c *CPU) { adc(c, c.loadAbsoluteX) },
	0x7e: func(c *CPU) { ror(c, c.loadAbsoluteX) },

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
	0x9a: func(c *CPU) { c.SP = c.X }, // txs
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
	0xba: func(c *CPU) { transfer(c, c.SP, &c.X) }, // tsx
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
	0xe1: func(c *CPU) { sbc(c, c.loadIndirectX) },
	0xe4: func(c *CPU) { cmp(c, c.X, c.loadZeroPage) },
	0xe5: func(c *CPU) { sbc(c, c.loadZeroPage) },
	0xe6: func(c *CPU) { inc(c, c.loadZeroPage) },
	0xe8: func(c *CPU) { inx(c) },
	0xe9: func(c *CPU) { sbc(c, c.loadImmediate) },
	0xea: func(c *CPU) {}, // nop
	0xec: func(c *CPU) { cmp(c, c.X, c.loadAbsolute) },
	0xed: func(c *CPU) { sbc(c, c.loadAbsolute) },
	0xee: func(c *CPU) { inc(c, c.loadAbsolute) },

	0xf0: func(c *CPU) { branch(c, c.Z) }, // beq
	0xf1: func(c *CPU) { sbc(c, c.loadIndirectY) },
	0xf5: func(c *CPU) { sbc(c, c.loadZeroPageX) },
	0xf6: func(c *CPU) { inc(c, c.loadZeroPageX) },
	0xf8: func(c *CPU) { c.D = true }, // sed
	0xf9: func(c *CPU) { sbc(c, c.loadAbsoluteY) },
	0xfd: func(c *CPU) { sbc(c, c.loadAbsoluteX) },
	0xfe: func(c *CPU) { inc(c, c.loadAbsoluteX) },
}
