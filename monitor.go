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
	dasmPageLen = 0x40
)

type Monitor struct {
	mach        *Mach85
	in          io.Reader
	out         *log.Logger
	dasm        *Disassembler
	interactive bool
	quit        bool
	lastCmd     string
	memPtr      uint16
	dasmPtr     uint16
}

func NewMonitor(mach *Mach85) *Monitor {
	mon := &Monitor{
		mach:        mach,
		in:          os.Stdin,
		out:         log.New(os.Stdout, "", 0),
		dasm:        NewDisassembler(mach.mem),
		interactive: true,
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
	fields := strings.Split(line, " ")

	if len(fields) == 0 {
		return
	}

	cmd := fields[0]
	args := fields[1:]
	var err error
	switch {
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

func (m *Monitor) disassemble(args []string) error {
	if err := checkLen(args, 0, 2); err != nil {
		return err
	}
	addrStart := m.mach.cpu.PC
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
	for m.dasm.PC = addrStart; m.dasm.PC < addrEnd; {
		m.out.Println(m.dasm.Next().String())
	}
	m.dasmPtr = m.dasm.PC
	return nil
}

func (m *Monitor) memory(args []string) error {
	if err := checkLen(args, 0, 2); err != nil {
		return err
	}
	addrStart := m.mach.cpu.PC + 1
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
	m.out.Println(m.mach.mem.Dump(addrStart, addrEnd))
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
		m.mach.mem.Store(address+uint16(offset), v)
	}
	return nil
}

func (m *Monitor) registers(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.out.Println(m.mach.cpu.String())
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
	if err := checkLen(args, 0, 0); err != nil {
		return err
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
		if m.mach.cpu.Trace == nil {
			m.out.Println("trace off")
		} else {
			m.out.Println("trace on")
		}
		return nil
	}
	switch args[0] {
	case "on":
		m.mach.cpu.Trace = func(op Operation) {
			m.out.Println(op)
		}
	case "off":
		m.mach.cpu.Trace = nil
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
