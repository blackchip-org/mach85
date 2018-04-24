package mach85

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func newTestMonitor() (*Monitor, *bytes.Buffer) {
	var out bytes.Buffer
	mach := New()
	mach.cpu.PC = 0x0800 - 1
	mon := NewMonitor(mach)
	mon.interactive = false
	mon.out.SetOutput(&out)
	return mon, &out
}

func TestGo(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.mem.StoreN(0x0800,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = strings.NewReader("g")
	mon.Run()
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoContinued(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.mem.StoreN(0x0800,
		0xa9, 0x12, // lda #$12
		0x00,       // brk
		0xea,       // nop
		0xa9, 0x34, // lda #34
		0x00, // brk
	)
	mon.in = strings.NewReader("g \n g")
	mon.Run()
	want := uint8(0x34)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoAddress(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.mem.StoreN(0x0900,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = strings.NewReader("g 0900")
	mon.Run()
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoAddressOldHexSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.mem.StoreN(0x0900,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = strings.NewReader("g $0900")
	mon.Run()
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoAddressNewHexSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.mem.StoreN(0x0900,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = strings.NewReader("g 0x0900")
	mon.Run()
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoAddressDecimalSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.mem.StoreN(0x0900,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = strings.NewReader("g +2304")
	mon.Run()
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoInvalidAddress(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = strings.NewReader("g foo")
	mon.Run()
	lines := strings.Split(out.String(), "\n")
	want := "invalid address: foo"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestGoTooManyArguments(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = strings.NewReader("g $0800 foo")
	mon.Run()
	lines := strings.Split(out.String(), "\n")
	want := "too many arguments"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
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
	mon.mach.mem.StoreN(0x0800+uint16(dasmPageLen-2),
		0xa9, 0x34, // lda #$34
	)
	mon.in = strings.NewReader("d")
	mon.Run()
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$083e: a9 34     lda #$34"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestTrace(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mach.mem.StoreN(0x0800,
		0xa9, 0x34, // lda #$34
	)
	mon.in = strings.NewReader("t on \n t \n g")
	mon.Run()
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := []string{
		"trace on",
		"$0800: a9 34     lda #$34",
		"$0802: 00        brk",
	}
	have := lines
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestTraceDisabled(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mach.mem.StoreN(0x0800,
		0xa9, 0x34, // lda #$34
	)
	mon.in = strings.NewReader("t on \n t off \n t \n g")
	mon.Run()
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := []string{
		"trace off",
	}
	have := lines
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}
