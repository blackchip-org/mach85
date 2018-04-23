package mach85

import (
	"bytes"
	"strings"
	"testing"
)

func newTestMonitor() (*Monitor, *bytes.Buffer) {
	var out bytes.Buffer
	mach := New()
	mach.cpu.PC = 0x0800 - 1
	mon := NewMonitor(mach)
	mon.out = &out
	return mon, &out
}

func TestRun(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.mem.StoreN(0x0800,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = strings.NewReader("r")
	mon.Run()
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestRunContinued(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.mem.StoreN(0x0800,
		0xa9, 0x12, // lda #$12
		0x00,       // brk
		0xea,       // nop
		0xa9, 0x34, // lda #34
		0x00, // brk
	)
	mon.in = strings.NewReader("r\nr")
	mon.Run()
	want := uint8(0x34)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestDisassembleFirstLine(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mach.mem.StoreN(0x0800,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = strings.NewReader("d")
	mon.Run()
	lines := strings.Split(out.String(), "\n")
	want := "$0800: a9 12     lda #$12"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestDisassembleLastLine(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mach.mem.StoreN(0x0800+uint16(mon.PageLen-1),
		0xa9, 0x34, // lda #$34
	)
	mon.in = strings.NewReader("d")
	mon.Run()
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$080f: a9 34     lda #$34"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}
