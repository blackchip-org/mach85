package mach85

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"
)

func newTestMonitor() (*Monitor, *bytes.Buffer) {
	var out bytes.Buffer
	mach := New()
	mach.StopOnBreak = true
	mach.QuitOnStop = true
	mach.cpu.PC = 0x0800 - 1
	mon := NewMonitor(mach)
	mon.Prompt = ""
	mon.out.SetOutput(&out)
	return mon, &out
}

func testMonitorRun(mon *Monitor) {
	go mon.Run()
	time.Sleep(time.Millisecond * 10)
	mon.mach.Start()
	mon.mach.Run()
}

func testMonitorInput(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func TestBreakpointOn(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.Memory.StoreN(0x0800, 0xea, 0xea, 0xea) // nop
	mon.in = testMonitorInput("b 0x0802 on \n g")
	testMonitorRun(mon)
	want := uint16(0x0801)
	have := mon.cpu.PC
	if want != have {
		t.Errorf("\n want: %04x \n have: %04x \n", want, have)
	}
}

func TestBreakpointOff(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mach.Memory.StoreN(0x0800, 0xea, 0xea, 0xea) // nop
	mon.in = testMonitorInput("b 0x0802 on \n b 0x0802 off \n g")
	testMonitorRun(mon)
	want := uint16(0x0804)
	have := mon.cpu.PC
	if want != have {
		fmt.Println(out.String())
		t.Errorf("\n want: %04x \n have: %04x \n", want, have)
	}
}

func TestBreakpointNotEnoughArguments(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("b")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "not enough arguments"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestBreakpointTooManyArguments(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("b 0x0800 on on")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "too many arguments"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestDisassembleFirstLine(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mach.Memory.StoreN(0x0800,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = testMonitorInput("d")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "$0800: a9 12     lda #$12"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestDisassembleLastLine(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mach.Memory.StoreN(0x0800+uint16(dasmPageLen-1),
		0xa9, 0x34, // lda #$34
	)
	mon.in = testMonitorInput("d")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$083e: a9 34     lda #$34"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestDisassemblePage(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("d 0800")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$083f: 00        brk"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestDisassembleNextPage(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("d 0800 \n d \n d")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$08bf: 00        brk"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestDisassembleRange(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("d 0800 0812")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$0812: 00        brk"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestDisassembleTooManyArguments(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("d 0800 0812 0812")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "too many arguments"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestGo(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.Memory.StoreN(0x0800,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = testMonitorInput("g")
	testMonitorRun(mon)
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoContinued(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.Memory.StoreN(0x0800,
		0xa9, 0x12, // lda #$12
		0x00,       // brk
		0xea,       // nop
		0xa9, 0x34, // lda #34
		0x00, // brk
	)
	mon.in = testMonitorInput("g")
	testMonitorRun(mon)
	mon.in = testMonitorInput("g")
	testMonitorRun(mon)
	want := uint8(0x34)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoAddress(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.Memory.StoreN(0x0900,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = testMonitorInput("g 0900")
	testMonitorRun(mon)
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoAddressOldHexSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.Memory.StoreN(0x0900,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = testMonitorInput("g $0900")
	testMonitorRun(mon)
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoAddressNewHexSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.Memory.StoreN(0x0900,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = testMonitorInput("g 0x0900")
	testMonitorRun(mon)
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoAddressDecimalSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.mach.Memory.StoreN(0x0900,
		0xa9, 0x12, // lda #$12
		0x00, // brk
	)
	mon.in = testMonitorInput("g +2304")
	testMonitorRun(mon)
	want := uint8(0x12)
	have := mon.mach.cpu.A
	if want != have {
		t.Errorf("\n want: %02x \n have: %02x \n", want, have)
	}
}

func TestGoInvalidAddress(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("g foo")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "invalid address: foo"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestGoTooManyArguments(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("g $0800 foo")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "too many arguments"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryFirstLine(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("m")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "$0800 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryLastLine(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("m")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$08f0 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryPage(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("m 0800")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$08f0 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryNextPage(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("m 0800 \n m \n m")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$0af0 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryRange(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("m 0800 081a")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$0810 00 00 00 00 00 00 00 00  00 00 00                ..........."
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryUnshifted(t *testing.T) {
	mon, out := newTestMonitor()
	for i := 0; i < 0x10; i++ {
		code := uint8(0x40 + i)
		mon.mem.Store(uint16(0x0800+i), code)
	}
	mon.in = testMonitorInput("m 0800 080f")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$0800 40 41 42 43 44 45 46 47  48 49 4a 4b 4c 4d 4e 4f @ABCDEFGHIJKLMNO"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryShifted(t *testing.T) {
	mon, out := newTestMonitor()
	for i := 0; i < 0x10; i++ {
		code := uint8(0x40 + i)
		mon.mem.Store(uint16(0x0800+i), code)
	}
	mon.in = testMonitorInput("M 0800 080f")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := "$0800 40 41 42 43 44 45 46 47  48 49 4a 4b 4c 4d 4e 4f @abcdefghijklmno"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryTooManyArguments(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("m 0800 0812 0812")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "too many arguments"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPoke(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.in = testMonitorInput("p 0900 ab")
	testMonitorRun(mon)
	want := uint8(0xab)
	have := mon.mem.Load(0x0900)
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPokeWord(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.in = testMonitorInput("pw 0900 abcd")
	testMonitorRun(mon)
	want := uint16(0xabcd)
	have := mon.mem.Load16(0x0900)
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPokeOldHexSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.in = testMonitorInput("p 0900 $ab")
	testMonitorRun(mon)
	want := uint8(0xab)
	have := mon.mem.Load(0x0900)
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPokeOldNewSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.in = testMonitorInput("p 0900 0xab")
	testMonitorRun(mon)
	want := uint8(0xab)
	have := mon.mem.Load(0x0900)
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPokeDecimalSigil(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.in = testMonitorInput("p 0900 +171")
	testMonitorRun(mon)
	want := uint8(0xab)
	have := mon.mem.Load(0x0900)
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPokeInvalid(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("p 0900 foo")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "invalid value: foo"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPokeOutOfRange(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("p 0900 1234")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "invalid value: 1234"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPokeN(t *testing.T) {
	mon, _ := newTestMonitor()
	mon.in = testMonitorInput("p 0900 ab cd ef 12 34")
	testMonitorRun(mon)
	want := uint8(0x34)
	have := mon.mem.Load(0x0904)
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPokeNotEnoughArguments(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("p")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "not enough arguments"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPeek(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mem.Store(0x0900, 0xab)
	mon.in = testMonitorInput("p 0900")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "$ab +171"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPeekWord(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mem.Store16(0x0900, 0xabcd)
	mon.in = testMonitorInput("pw 0900")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "$abcd +43981"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestTrace(t *testing.T) {
	mon, out := newTestMonitor()
	mon.mach.Memory.StoreN(0x0800,
		0xa9, 0x34, // lda #$34
	)
	mon.in = testMonitorInput("t on \n t \n g")
	testMonitorRun(mon)
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
	mon.mach.Memory.StoreN(0x0800,
		0xa9, 0x34, // lda #$34
	)
	mon.in = testMonitorInput("t on \n t off \n t \n g")
	testMonitorRun(mon)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := []string{
		"trace off",
	}
	have := lines
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestTraceTooManyArguments(t *testing.T) {
	mon, out := newTestMonitor()
	mon.in = testMonitorInput("t on on")
	testMonitorRun(mon)
	lines := strings.Split(out.String(), "\n")
	want := "too many arguments"
	have := lines[0]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}
