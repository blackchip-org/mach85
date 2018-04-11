package mach85

import (
	"testing"
)

func newTestCPU() *CPU {
	mem := NewMemory(0x0300)
	c := NewCPU(mem)
	c.SP = 0xff
	c.PC = 0x1ff
	return c
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

func flagError(t *testing.T, want uint8, have uint8) {
	t.Errorf("\n       nv-bdizc\n want: %08b \n have: %08b \n", want, have)
}

// ----------------------------------------------------------------------------
// and
// ----------------------------------------------------------------------------
func TestAndImmediate(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x29) // and #$0f
	c.mem.Store(0x0201, 0x0f)
	c.A = 0xcd
	c.Run()
	want := uint8(0x0d)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
	want = flagB | flag5
	have = c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestAndZero(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x29) // and #$f0
	c.mem.Store(0x0201, 0xf0)
	c.A = 0x0f
	c.Run()
	want := flagZ | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestAndSigned(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa9) // and #$f0
	c.mem.Store(0x0201, 0xf0)
	c.A = 0xff
	c.Run()
	want := flagN | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestAndZeroPage(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0034, 0x0f) // .byte $0f
	c.mem.Store(0x0200, 0x25) // and $34
	c.mem.Store(0x0201, 0x34)
	c.A = 0xcd
	c.Run()
	want := uint8(0x0d)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAndZeroPageX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0034, 0x0f) // .byte $0f
	c.mem.Store(0x0200, 0x35) // and $30,X
	c.mem.Store(0x0201, 0x30)
	c.A = 0xcd
	c.X = 0x04
	c.Run()
	want := uint8(0x0d)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAndAbsolute(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x0f) // .byte $0f
	c.mem.Store(0x0200, 0x2d) // and $02ab
	c.mem.Store16(0x0201, 0x02ab)
	c.A = 0xcd
	c.Run()
	want := uint8(0x0d)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAndAbsoluteX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x0f) // .byte $0f
	c.mem.Store(0x0200, 0x3d) // and $02a0,X
	c.mem.Store16(0x0201, 0x02a0)
	c.A = 0xcd
	c.X = 0x0b
	c.Run()
	want := uint8(0x0d)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAndAbsoluteY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x0f) // .byte $0f
	c.mem.Store(0x0200, 0x39) // and $02a0,Y
	c.mem.Store16(0x0201, 0x02a0)
	c.A = 0xcd
	c.Y = 0x0b
	c.Run()
	want := uint8(0x0d)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAndIndirectX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store16(0x4a, 0x02ab) // .word $02ab
	c.mem.Store(0x02ab, 0x0f)   // .byte $0f
	c.mem.Store(0x0200, 0x21)   // and ($40,X)
	c.mem.Store(0x0201, 0x40)
	c.A = 0xcd
	c.X = 0x0a
	c.Run()
	want := uint8(0x0d)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAndIndirectY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store16(0x4a, 0x02a0) // .word $02a0
	c.mem.Store(0x02ab, 0x0f)   // .byte $0f
	c.mem.Store(0x0200, 0x31)   // and ($4a),Y
	c.mem.Store(0x0201, 0x4a)
	c.A = 0xcd
	c.Y = 0x0b
	c.Run()
	want := uint8(0x0d)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

// ----------------------------------------------------------------------------
// asl
// ----------------------------------------------------------------------------
func TestAslAccumulator(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x0a) // asl a
	c.A = 4
	c.Run()
	want := uint8(8)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAslSigned(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x0a) // asl a
	c.A = 1 << 6
	c.Run()
	want := flagN | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestAslCarry(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x0a) // asl a
	c.A = 1 << 7
	c.Run()
	want := flagC | flagZ | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestAslZeroPage(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x00ab, 4)    // .byte 4
	c.mem.Store(0x0200, 0x06) // asl $ab
	c.mem.Store(0x0201, 0xab)
	c.Run()
	want := uint8(8)
	have := c.mem.Load(0x00ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAslZeroPageX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x00ab, 4)    // .byte 4
	c.mem.Store(0x0200, 0x16) // asl $a0
	c.mem.Store(0x0201, 0xa0)
	c.X = 0x0b
	c.Run()
	want := uint8(8)
	have := c.mem.Load(0x00ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAslAbsolute(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 4)    // .byte 4
	c.mem.Store(0x0200, 0x0e) // asl $02ab
	c.mem.Store16(0x0201, 0x02ab)
	c.Run()
	want := uint8(8)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestAslAbsoluteX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 4)    // .byte 4
	c.mem.Store(0x0200, 0x1e) // asl $02a0,X
	c.mem.Store16(0x0201, 0x02a0)
	c.X = 0x0b
	c.Run()
	want := uint8(8)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

// ----------------------------------------------------------------------------
// bit
// ----------------------------------------------------------------------------
var bitTests = []struct {
	name          string
	a             uint8
	fetch         uint8
	expectedFlags uint8
}{
	{"zero", 0x00, 0x00, flagZ | flagB | flag5},
	{"non-zero", 0x01, 0x01, flagB | flag5},
	{"and-zero", 0x01, 0x02, flagZ | flagB | flag5},
	{"bit6", 0x00, 1 << 6, flagV | flagZ | flagB | flag5},
	{"bit7", 0x00, 1 << 7, flagN | flagZ | flagB | flag5},
}

func TestBitAbsolute(t *testing.T) {
	for _, test := range bitTests {
		t.Run(test.name, func(t *testing.T) {
			c := newTestCPU()
			c.mem.Store(0x02ab, test.fetch) // .byte fetch
			c.mem.Store(0x0200, 0x2c)       // bit $02ab
			c.mem.Store16(0x0201, 0x02ab)
			c.A = test.a
			c.Run()
			want := test.expectedFlags
			have := c.SR()
			if want != have {
				flagError(t, want, have)
			}
		})
	}
}

func TestBitZeroPage(t *testing.T) {
	for _, test := range bitTests {
		t.Run(test.name, func(t *testing.T) {
			c := newTestCPU()
			c.mem.Store(0xab, test.fetch) // .byte fetch
			c.mem.Store(0x0200, 0x2c)     // bit $ab
			c.mem.Store(0x0201, 0xab)
			c.A = test.a
			c.Run()
			want := test.expectedFlags
			have := c.SR()
			if want != have {
				flagError(t, want, have)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// branches
// ----------------------------------------------------------------------------
var branchTests = []struct {
	name      string
	op        uint8
	flags     uint8
	expectedA uint8
}{
	{"bpl yea", 0x10, 0, 0x02},
	{"bpl nay", 0x10, flagN, 0x01},
	{"bmi yea", 0x30, flagN, 0x02},
	{"bmi nay", 0x30, 0, 0x01},
	{"bvc yea", 0x50, 0, 0x02},
	{"bvc nay", 0x50, flagV, 0x01},
	{"bvs yea", 0x70, flagV, 0x02},
	{"bvs nay", 0x70, 0, 0x01},
	{"bcc yea", 0x90, 0, 0x02},
	{"bcc nay", 0x90, flagC, 0x01},
	{"bcs yea", 0xb0, flagC, 0x02},
	{"bcs nay", 0xb0, 0, 0x01},
	{"bne yea", 0xd0, 0, 0x02},
	{"bne nay", 0xd0, flagZ, 0x01},
	{"beq yea", 0xf0, flagZ, 0x02},
	{"beq nay", 0xf0, 0, 0x01},
}

func TestBranches(t *testing.T) {
	for _, test := range branchTests {
		t.Run(test.name, func(t *testing.T) {
			c := newTestCPU()
			c.mem.Store(0x0200, test.op) // branch to $0205
			c.mem.Store(0x0201, 0x03)    //
			c.mem.Store(0x0202, 0xa9)    // lda #01
			c.mem.Store(0x0203, 0x01)    //
			c.mem.Store(0x0204, 0x00)    // brk
			c.mem.Store(0x0205, 0xa9)    // lda #02
			c.mem.Store(0x0206, 0x02)
			c.SetSR(test.flags)
			c.Run()
			want := test.expectedA
			have := c.A
			if want != have {
				t.Errorf("\n want: %02x \n have: %02x \n", want, have)
			}
		})
	}
}

func TestBranchBackwards(t *testing.T) {
	c := newTestCPU()
	c.PC = 0x05
	c.mem.Store(0x06, 0xf0) // beq $02
	c.mem.Store(0x07, 0xfb) // -5
	c.Z = true
	c.Next()
	want := uint16(0x02)
	have := c.PC
	if want != have {
		t.Errorf("\n want: %04x \n have: %04x \n", want, have)
	}
}

// ----------------------------------------------------------------------------
// lda
// ----------------------------------------------------------------------------
func TestLdaImmediate(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa9) // lda #$12
	c.mem.Store(0x0201, 0x12)
	c.Run()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
	want = flagB | flag5
	have = c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdaZero(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa9) // lda #$00
	c.mem.Store(0x0201, 0x00)
	c.Run()
	want := flagZ | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdaSigned(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa9) // lda #$ff
	c.mem.Store(0x0201, 0xff)
	c.Run()
	want := flagN | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdaZeroPage(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0034, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xa5) // lda $34
	c.mem.Store(0x0201, 0x34)
	c.Run()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdaZeroPageX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0034, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xb5) // lda $30,X
	c.mem.Store(0x0201, 0x30)
	c.X = 0x4
	c.Run()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdaAbsolute(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xad) // lda $02ab
	c.mem.Store16(0x0201, 0x02ab)
	c.Run()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdaAbsoluteX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xbd) // lda $02a0,X
	c.mem.Store16(0x0201, 0x02a0)
	c.X = 0xb
	c.Run()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdaAbsoluteY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xb9) // lda $02a0,Y
	c.mem.Store16(0x0201, 0x02a0)
	c.Y = 0xb
	c.Run()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdaIndirectX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store16(0x4a, 0x02ab) // .word $02ab
	c.mem.Store(0x02ab, 0x12)   // .byte $12
	c.mem.Store(0x0200, 0xa1)   // lda ($40,X)
	c.mem.Store(0x0201, 0x40)
	c.X = 0xa
	c.Run()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdaIndirectY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store16(0x4a, 0x02a0) // .word $02a0
	c.mem.Store(0x02ab, 0x12)   // .byte $12
	c.mem.Store(0x0200, 0xb1)   // lda ($4a),Y
	c.mem.Store(0x0201, 0x4a)
	c.Y = 0xb
	c.Run()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

// ----------------------------------------------------------------------------
// ldx
// ----------------------------------------------------------------------------
func TestLdxImmediate(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa2) // ldx #$12
	c.mem.Store(0x0201, 0x12)
	c.Run()
	want := uint8(0x12)
	have := c.X
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
	want = flagB | flag5
	have = c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdxZero(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa2) // ldx #$00
	c.mem.Store(0x0201, 0x00)
	c.Run()
	want := flagZ | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdxSigned(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa2) // ldx #$ff
	c.mem.Store(0x0201, 0xff)
	c.Run()
	want := flagN | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdxZeroPage(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0034, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xa6) // ldx $34
	c.mem.Store(0x0201, 0x34)
	c.Run()
	want := uint8(0x12)
	have := c.X
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdxZeroPageY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0034, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xb6) // ldx $30,Y
	c.mem.Store(0x0201, 0x30)
	c.Y = 0x4
	c.Run()
	want := uint8(0x12)
	have := c.X
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdxAbsolute(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xae) // ldx $02ab
	c.mem.Store16(0x0201, 0x02ab)
	c.Run()
	want := uint8(0x12)
	have := c.X
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdxAbsoluteY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xbe) // ldx $02a0,Y
	c.mem.Store16(0x0201, 0x02a0)
	c.Y = 0xb
	c.Run()
	want := uint8(0x12)
	have := c.X
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

// ----------------------------------------------------------------------------
// ldy
// ----------------------------------------------------------------------------
func TestLdyImmediate(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa0) // ldy #$12
	c.mem.Store(0x0201, 0x12)
	c.Run()
	want := uint8(0x12)
	have := c.Y
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
	want = flagB | flag5
	have = c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdyZero(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa0) // ldy #$00
	c.mem.Store(0x0201, 0x00)
	c.Run()
	want := flagZ | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdySigned(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xa0) // ldy #$ff
	c.mem.Store(0x0201, 0xff)
	c.Run()
	want := flagN | flagB | flag5
	have := c.SR()
	if want != have {
		flagError(t, want, have)
	}
}

func TestLdyZeroPage(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0034, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xa4) // ldy $34
	c.mem.Store(0x0201, 0x34)
	c.Run()
	want := uint8(0x12)
	have := c.Y
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdyZeroPageX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0034, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xb4) // ldy $30,X
	c.mem.Store(0x0201, 0x30)
	c.X = 0x4
	c.Run()
	want := uint8(0x12)
	have := c.Y
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdyAbsolute(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xac) // ldy $02ab
	c.mem.Store16(0x0201, 0x02ab)
	c.Run()
	want := uint8(0x12)
	have := c.Y
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestLdyAbsoluteX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x02ab, 0x12) // .byte $12
	c.mem.Store(0x0200, 0xbc) // ldy $02a0,X
	c.mem.Store16(0x0201, 0x02a0)
	c.X = 0xb
	c.Run()
	want := uint8(0x12)
	have := c.Y
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

// ----------------------------------------------------------------------------
// sta
// ----------------------------------------------------------------------------

func TestStaZeroPage(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x85) // sta $34
	c.mem.Store(0x0201, 0x34)
	c.A = 0x12
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x34)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStaZeroPageX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x95) // sta $30,X
	c.mem.Store(0x0201, 0x30)
	c.A = 0x12
	c.X = 0x04
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x34)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStaAbsolute(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x8d) // sta $02ab
	c.mem.Store16(0x0201, 0x02ab)
	c.A = 0x12
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStaAbsoluteX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x9d) // sta $02a0,X
	c.mem.Store16(0x0201, 0x02a0)
	c.A = 0x12
	c.X = 0x0b
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStaAbsoluteY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x99) // sta $02a0,Y
	c.mem.Store16(0x0201, 0x02a0)
	c.A = 0x12
	c.Y = 0x0b
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStaIndirectX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store16(0x4a, 0x02ab) // .word $02ab
	c.mem.Store(0x0200, 0x81)   // sta ($40,X)
	c.mem.Store(0x0201, 0x40)
	c.A = 0x12
	c.X = 0x0a
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStaIndirectY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store16(0x4a, 0x02a0) // .word $02a0
	c.mem.Store(0x0200, 0x91)   // sta ($4a),Y
	c.mem.Store(0x0201, 0x4a)
	c.A = 0x12
	c.Y = 0x0b
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

// ----------------------------------------------------------------------------
// stx
// ----------------------------------------------------------------------------

func TestStxZeroPage(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x86) // stx $34
	c.mem.Store(0x0201, 0x34)
	c.X = 0x12
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x34)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStxZeroPageY(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x96) // stx $30,Y
	c.mem.Store(0x0201, 0x30)
	c.X = 0x12
	c.Y = 0x04
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x34)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStxAbsolute(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x8e) // stx $02ab
	c.mem.Store16(0x0201, 0x02ab)
	c.X = 0x12
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

// ----------------------------------------------------------------------------
// sty
// ----------------------------------------------------------------------------

func TestStyZeroPage(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x84) // sty $34
	c.mem.Store(0x0201, 0x34)
	c.Y = 0x12
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x34)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStyZeroPageX(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x94) // sty $30,X
	c.mem.Store(0x0201, 0x30)
	c.Y = 0x12
	c.X = 0x04
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x34)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestStyAbsolute(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0x8c) // sty $02ab
	c.mem.Store16(0x0201, 0x02ab)
	c.Y = 0x12
	c.Run()
	want := uint8(0x12)
	have := c.mem.Load(0x02ab)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}
