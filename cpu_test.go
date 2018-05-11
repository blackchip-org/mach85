package mach85

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func newTestCPU() *CPU {
	mem := NewMemory(NewRAM(0x10000))
	c := New6510(mem)
	c.SP = 0xff
	c.PC = 0x1ff
	return c
}

func testRunCPU(cpu *CPU) {
	for !cpu.B {
		cpu.Next()
	}
}

func TestIllegalOpcode(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() { log.SetOutput(os.Stderr) }()

	c := newTestCPU()
	// http://visual6502.org/wiki/index.php?title=6502_all_256_Opcodes
	c.mem.Store(0x0200, 0x02) // *KIL
	testRunCPU(c)
	if !strings.Contains(buf.String(), "illegal") {
		t.Errorf("illegal instruction not logged")
	}
}

var cpuStringTests = []struct {
	setup func(c *CPU)
	want  string
}{
	{func(c *CPU) { c.PC = 0x1234 },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"1234 20 00 00 00 00  . . * . . . . ."},
	{func(c *CPU) { c.A = 0x56 },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 20 56 00 00 00  . . * . . . . ."},
	{func(c *CPU) { c.X = 0x78 },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 20 00 78 00 00  . . * . . . . ."},
	{func(c *CPU) { c.Y = 0x9a },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 20 00 00 9a 00  . . * . . . . ."},
	{func(c *CPU) { c.SP = 0xbc },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 20 00 00 00 bc  . . * . . . . ."},
	{func(c *CPU) { c.C = true },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 21 00 00 00 00  . . * . . . . *"},
	{func(c *CPU) { c.Z = true },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 22 00 00 00 00  . . * . . . * ."},
	{func(c *CPU) { c.I = true },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 24 00 00 00 00  . . * . . * . ."},
	{func(c *CPU) { c.D = true },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 28 00 00 00 00  . . * . * . . ."},
	{func(c *CPU) { c.B = true },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 30 00 00 00 00  . . * * . . . ."},
	{func(c *CPU) { c.V = true },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 60 00 00 00 00  . * * . . . . ."},
	{func(c *CPU) { c.N = true },
		"" +
			" pc  sr ac xr yr sp  n v - b d i z c\n" +
			"0000 a0 00 00 00 00  * . * . . . . ."},
}

func TestCPUString(t *testing.T) {
	for _, test := range cpuStringTests {
		c := New6510(NewMemory(NewRAM(0x10000)))
		test.setup(c)
		have := c.String()
		if test.want != have {
			t.Errorf("\n want: \n%v \n have: \n%v\n", test.want, have)
		}
	}
}
func TestPush(t *testing.T) {
	c := New6510(NewMemory(NewRAM(0x10000)))
	c.SP = 0xff
	c.push(0x12)
	c.push(0x34)
	c.push(0x56)
	want := uint8(0x56)
	have := c.mem.Load(AddrStack + 0x100 - 3)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}

func TestPush16(t *testing.T) {
	c := New6510(NewMemory(NewRAM(0x10000)))
	c.SP = 0xff
	c.push(0x12)
	c.push16(0x3456)
	want := uint8(0x56)
	have := c.mem.Load(AddrStack + 0x100 - 3)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}

func TestPushOverflow(t *testing.T) {
	c := New6510(NewMemory(NewRAM(0x10000)))
	c.SP = 0x01
	c.push(0x12)
	c.push(0x34)
	c.push(0x56)
	want := uint8(0x56)
	have := c.mem.Load(AddrStack + 0x100 - 1)
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}

func TestPull(t *testing.T) {
	c := New6510(NewMemory(NewRAM(0x10000)))
	c.mem.Store(AddrStack+0xff, 0x12)
	c.SP = 0xfe
	want := uint8(0x12)
	have := c.pull()
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}

func TestPull16(t *testing.T) {
	c := New6510(NewMemory(NewRAM(0x10000)))
	c.mem.Store(AddrStack+0xfe, 0x34)
	c.mem.Store(AddrStack+0xff, 0x12)
	c.SP = 0xfd
	want := uint16(0x1234)
	have := c.pull16()
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}

func TestIRQ(t *testing.T) {
	c := newTestCPU()
	c.mem.Store(0x0200, 0xea)         // nop
	c.mem.StoreN(AddrISR, 0xa9, 0x12) // lda #12
	c.IRQ()
	c.Next()
	c.Next()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}

func TestIRQIgnore(t *testing.T) {
	c := newTestCPU()
	c.mem.StoreN(0x0200, 0xa9, 0x12)
	c.I = true
	c.IRQ()
	c.Next()
	want := uint8(0x12)
	have := c.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x\n", want, have)
	}
}
