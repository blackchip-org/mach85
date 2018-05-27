package mach85

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

const (
	CmdBreakpoint          = "b"
	CmdDisassemble         = "d"
	CmdGo                  = "g"
	CmdHalt                = "h"
	CmdLoad                = "l"
	CmdLoadProgram         = "lp"
	CmdMemory              = "m"
	CmdMemoryShifted       = "M"
	CmdNext                = "n"
	CmdScreenMemory        = "sm"
	CmdScreenMemoryShifted = "SM"
	CmdPokePeek            = "p"
	CmdPokePeekWord        = "pw"
	CmdStep                = "s"
	CmdQuit                = "q"
	CmdQuitLong            = "quit"
	CmdRegisters           = "r"
	CmdTrace               = "t"
	CmdZap                 = "z"
)

var (
	memPageLen  = 0x100
	dasmPageLen = 0x3f
)

type Monitor struct {
	Disassembler *Disassembler
	Prompt       string
	mach         *Mach85
	cpu          *CPU
	mem          *Memory
	in           io.ReadCloser
	out          *log.Logger
	lastCmd      string
	memPtr       uint16
	dasmPtr      uint16
}

func NewMonitor(mach *Mach85) *Monitor {
	mon := &Monitor{
		Prompt:       "mach85> ",
		mach:         mach,
		cpu:          mach.cpu,
		mem:          mach.Memory,
		in:           ioutil.NopCloser(os.Stdin),
		out:          log.New(os.Stdout, "", 0),
		Disassembler: NewDisassembler(mach.Memory),
	}
	return mon
}

const maxArgs = 0x100

func (m *Monitor) Run() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      m.Prompt,
		HistoryFile: filepath.Join(usr.HomeDir, ".mach85.history"),
		Stdin:       m.in,
	})
	if err != nil {
		return err
	}
	for {
		line, err := rl.Readline()
		if err != nil {
			return err
		}
		m.parse(line)
	}
}

func (m *Monitor) Go() {
	m.goCmd([]string{})
	m.Run()
}

func (m *Monitor) parse(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	fields := strings.Split(line, " ")

	if len(fields) == 0 {
		return
	}

	cmd := fields[0]
	args := fields[1:]
	var err error
	switch cmd {
	case CmdBreakpoint:
		err = m.breakpoint(args)
	case CmdDisassemble:
		err = m.disassemble(args)
	case CmdLoad:
		err = m.load(args)
	case CmdLoadProgram:
		err = m.loadProgram(args)
	case CmdGo:
		err = m.goCmd(args)
	case CmdHalt:
		err = m.halt(args)
	case CmdMemory:
		err = m.memory(args, PetsciiUnshiftedDecoder)
	case CmdMemoryShifted:
		err = m.memory(args, PetsciiShiftedDecoder)
	case CmdScreenMemory:
		err = m.memory(args, ScreenUnshiftedDecoder)
	case CmdScreenMemoryShifted:
		err = m.memory(args, ScreenShiftedDecoder)
	case CmdNext:
		err = m.next(args)
	case CmdStep:
		err = m.step(args)
	case CmdPokePeek:
		err = m.pokePeek(args)
	case CmdPokePeekWord:
		err = m.pokePeekWord(args)
	case CmdQuit, CmdQuitLong:
		os.Exit(0)
	case CmdRegisters:
		err = m.registers(args)
	case CmdTrace:
		err = m.trace(args)
	case CmdZap:
		err = m.zap(args)
	default:
		err = fmt.Errorf("unknown command: %v", cmd)
	}
	if err != nil {
		m.out.Println(err)
	} else {
		m.lastCmd = cmd
	}
}

func (m *Monitor) breakpoint(args []string) error {
	if err := checkLen(args, 1, 2); err != nil {
		return err
	}
	address, err := parseAddress(args[0])
	if err != nil {
		return err
	}
	if len(args) == 1 {
		if !m.mach.Breakpoints[address] {
			m.out.Println("breakpoint off")
		} else {
			m.out.Println("breakpoint on")
		}
		return nil
	}
	switch args[1] {
	case "on":
		m.mach.Breakpoints[address] = true
	case "off":
		delete(m.mach.Breakpoints, address)
	default:
		return fmt.Errorf("invalid: %v", args[1])
	}
	return nil
}

func (m *Monitor) disassemble(args []string) error {
	if err := checkLen(args, 0, 2); err != nil {
		return err
	}
	addrStart := m.cpu.PC
	if len(args) == 0 {
		if strings.HasPrefix(m.lastCmd, CmdDisassemble) {
			addrStart = m.dasmPtr
		}
	}
	if len(args) > 0 {
		addr, err := parseAddress(args[0])
		if err != nil {
			return err
		}
		addrStart = addr - 1
	}
	addrEnd := addrStart + uint16(dasmPageLen)
	if len(args) > 1 {
		addr, err := parseAddress(args[1])
		if err != nil {
			return err
		}
		addrEnd = addr - 1
	}
	for m.Disassembler.PC = addrStart; m.Disassembler.PC <= addrEnd; {
		m.out.Println(m.Disassembler.Next().String())
	}
	m.dasmPtr = m.Disassembler.PC
	return nil
}

func (m *Monitor) halt(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.mach.Stop()
	return nil
}

func (m *Monitor) load(args []string) error {
	if err := checkLen(args, 2, 2); err != nil {
		return err
	}
	addr, err := parseAddress(args[0])
	if err != nil {
		return err
	}
	filename := args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	m.mem.Import(addr, data)
	return nil
}

func (m *Monitor) loadProgram(args []string) error {
	if err := checkLen(args, 1, 1); err != nil {
		return err
	}
	filename := args[0]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	addr := uint16(data[0]) + uint16(data[1])<<8
	end := addr + uint16(len(data)) - 2
	m.out.Printf("$%04x - $%04x\n", addr, end)
	m.mem.Import(addr, data[2:])
	return nil
}

func (m *Monitor) memory(args []string, decoder Decoder) error {
	if err := checkLen(args, 0, 2); err != nil {
		return err
	}
	addrStart := m.cpu.PC + 1
	if len(args) == 0 {
		if m.lastCmd == CmdMemory || m.lastCmd == CmdMemoryShifted {
			addrStart = m.memPtr
		}
	}
	if len(args) > 0 {
		addr, err := parseAddress(args[0])
		if err != nil {
			return err
		}
		addrStart = addr
	}
	addrEnd := addrStart + uint16(memPageLen)
	if len(args) > 1 {
		addr, err := parseAddress(args[1])
		if err != nil {
			return err
		}
		addrEnd = addr
	}
	m.out.Println(m.mem.Dump(addrStart, addrEnd, decoder))
	m.memPtr = addrEnd
	return nil
}

func (m *Monitor) next(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.Disassembler.PC = m.cpu.PC
	m.out.Println(m.Disassembler.Next())
	return nil
}

func (m *Monitor) pokePeek(args []string) error {
	if err := checkLen(args, 1, maxArgs); err != nil {
		return err
	}
	address, err := parseAddress(args[0])
	if err != nil {
		return err
	}
	// peek
	if len(args) == 1 {
		v := m.mem.Load(address)
		m.out.Printf("$%02x +%d\n", v, v)
		return nil
	}
	// poke
	values := []uint8{}
	for _, str := range args[1:] {
		v, err := parseValue(str)
		if err != nil {
			return err
		}
		values = append(values, v)
	}
	for offset, v := range values {
		m.mem.Store(address+uint16(offset), v)
	}
	return nil
}

func (m *Monitor) pokePeekWord(args []string) error {
	if err := checkLen(args, 1, maxArgs); err != nil {
		return err
	}
	address, err := parseAddress(args[0])
	if err != nil {
		return err
	}
	// peek
	if len(args) == 1 {
		v := m.mem.Load16(address)
		m.out.Printf("$%04x +%d\n", v, v)
		return nil
	}
	// poke
	values := []uint16{}
	for _, str := range args[1:] {
		v, err := parseValue16(str)
		if err != nil {
			return err
		}
		values = append(values, v)
	}
	for offset, v := range values {
		m.mem.Store16(address+uint16(offset), v)
	}
	return nil
}

func (m *Monitor) step(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.Disassembler.PC = m.cpu.PC
	m.out.Println(m.Disassembler.Next())
	m.cpu.Next()
	return nil
}

func (m *Monitor) registers(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	reason := ""
	if m.mach.Err != nil {
		reason = fmt.Sprintf(": %v", m.mach.Err)
	}
	m.out.Printf("[%v%v]\n", m.mach.Status, reason)

	m.out.Println(m.cpu.String())
	return nil
}

func (m *Monitor) goCmd(args []string) error {
	if err := checkLen(args, 0, 1); err != nil {
		return err
	}
	if len(args) > 0 {
		address, err := parseAddress(args[0])
		if err != nil {
			return err
		}
		m.cpu.PC = address - 1
	}
	go m.mach.Start()
	return nil
}

func (m *Monitor) trace(args []string) error {
	if err := checkLen(args, 0, 1); err != nil {
		return err
	}
	if len(args) == 0 {
		if m.mach.Trace == nil {
			m.out.Println("trace off")
		} else {
			m.out.Println("trace on")
		}
		return nil
	}
	switch args[0] {
	case "on":
		m.traceOn()
	case "off":
		m.traceOff()
	default:
		return fmt.Errorf("invalid: %v", args[0])
	}
	return nil
}

func (m *Monitor) zap(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	//go m.signalHandler()
	m.mach.Reset()
	return nil
}

func (m *Monitor) signalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	signal.Reset(os.Interrupt)
	m.out.Println()
	m.mach.Stop()
}

func (m *Monitor) traceOn() {
	m.mach.Trace = func(op Operation) {
		m.out.Println(op)
	}
}

func (m *Monitor) traceOff() {
	m.mach.Trace = nil
}

func checkLen(args []string, min int, max int) error {
	if len(args) < min {
		return errors.New("not enough arguments")
	}
	if len(args) > max {
		return errors.New("too many arguments")
	}
	return nil
}

func parseUint(str string, bitSize int) (uint64, error) {
	base := 16
	switch {
	case strings.HasPrefix(str, "$"):
		str = str[1:]
	case strings.HasPrefix(str, "0x"):
		str = str[2:]
	case strings.HasPrefix(str, "+"):
		str = str[1:]
		base = 10
	}
	return strconv.ParseUint(str, base, bitSize)
}

func parseAddress(str string) (uint16, error) {
	value, err := parseUint(str, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid address: %v", str)
	}
	return uint16(value), nil
}

func parseValue(str string) (uint8, error) {
	value, err := parseUint(str, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid value: %v", str)
	}
	return uint8(value), nil
}

func parseValue16(str string) (uint16, error) {
	value, err := parseUint(str, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid value: %v", str)
	}
	return uint16(value), nil
}
