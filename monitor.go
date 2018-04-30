package mach85

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

const (
	CmdBreakpoint  = "b"
	CmdDisassemble = "d"
	CmdGo          = "g"
	CmdMemory      = "m"
	CmdPoke        = "p"
	CmdQuit        = "q"
	CmdReset       = "reset"
	CmdRegisters   = "r"
	CmdTrace       = "t"
)

var (
	memPageLen  = 0x100
	dasmPageLen = 0x3f
)

type Monitor struct {
	Disassembler *Disassembler
	mach         *Mach85
	cpu          *CPU
	mem          *Memory
	in           io.Reader
	out          *log.Logger
	interactive  bool
	quit         bool
	lastCmd      string
	memPtr       uint16
	dasmPtr      uint16
}

func NewMonitor(mach *Mach85) *Monitor {
	mon := &Monitor{
		mach:         mach,
		cpu:          mach.cpu,
		mem:          mach.mem,
		in:           os.Stdin,
		out:          log.New(os.Stdout, "", 0),
		Disassembler: NewDisassembler(mach.mem),
		interactive:  true,
	}
	return mon
}

const maxArgs = 0x100

func (m *Monitor) Run() {
	s := bufio.NewScanner(m.in)
	s.Split(bufio.ScanLines)
	for {
		if m.interactive {
			fmt.Print("mach85> ")
		}
		if !s.Scan() {
			return
		}
		m.parse(s.Text())
		if m.quit {
			return
		}
	}
}

func (m *Monitor) parse(line string) {
	line = strings.ToLower(line)
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
	switch {
	case strings.HasPrefix(cmd, CmdBreakpoint):
		err = m.breakpoint(args)
	case strings.HasPrefix(cmd, CmdDisassemble):
		err = m.disassemble(args)
	case strings.HasPrefix(cmd, CmdGo):
		err = m.goCmd(args)
	case strings.HasPrefix(cmd, CmdMemory):
		err = m.memory(args)
	case strings.HasPrefix(cmd, CmdPoke):
		err = m.poke(args)
	case strings.HasPrefix(cmd, CmdQuit):
		m.quit = true
		return
	case cmd == CmdReset:
		err = m.reset(args)
	case strings.HasPrefix(cmd, CmdRegisters):
		err = m.registers(args)
	case strings.HasPrefix(cmd, CmdTrace):
		err = m.trace(args)
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
		if !m.cpu.Breakpoints[address] {
			m.out.Println("breakpoint off")
		} else {
			m.out.Println("breakpoint on")
		}
		return nil
	}
	switch args[1] {
	case "on":
		m.cpu.Breakpoints[address] = true
	case "off":
		delete(m.cpu.Breakpoints, address)
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
		addrStart = addr
	}
	addrEnd := addrStart + uint16(dasmPageLen)
	if len(args) > 1 {
		addr, err := parseAddress(args[1])
		if err != nil {
			return err
		}
		addrEnd = addr
	}
	for m.Disassembler.PC = addrStart; m.Disassembler.PC < addrEnd; {
		m.out.Println(m.Disassembler.Next().String())
	}
	m.dasmPtr = m.Disassembler.PC + 1
	return nil
}

func (m *Monitor) memory(args []string) error {
	if err := checkLen(args, 0, 2); err != nil {
		return err
	}
	addrStart := m.cpu.PC + 1
	if len(args) == 0 {
		if strings.HasPrefix(m.lastCmd, CmdMemory) {
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
	m.out.Println(m.mem.Dump(addrStart, addrEnd))
	m.memPtr = addrEnd
	return nil
}

func (m *Monitor) poke(args []string) error {
	if err := checkLen(args, 2, maxArgs); err != nil {
		return err
	}
	address, err := parseAddress(args[0])
	if err != nil {
		return err
	}
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

func (m *Monitor) registers(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.out.Println(m.cpu.String())
	return nil
}

func (m *Monitor) reset(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	go m.signalHandler()
	m.mach.Reset()
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
	go m.signalHandler()
	m.mach.Run()
	return nil
}

func (m *Monitor) trace(args []string) error {
	if err := checkLen(args, 0, 1); err != nil {
		return err
	}
	if len(args) == 0 {
		if m.cpu.Trace == nil {
			m.out.Println("trace off")
		} else {
			m.out.Println("trace on")
		}
		return nil
	}
	switch args[0] {
	case "on":
		m.cpu.Trace = func(op Operation) {
			m.out.Println(op)
		}
	case "off":
		m.cpu.Trace = nil
	default:
		return fmt.Errorf("invalid: %v", args[0])
	}
	return nil
}

func (m *Monitor) signalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	signal.Reset(os.Interrupt)
	m.mach.Stop()
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
