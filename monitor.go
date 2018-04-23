package mach85

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
)

type Monitor struct {
	PageLen     int
	mach        *Mach85
	in          io.Reader
	out         *log.Logger
	dasm        *Disassembler
	interactive bool
	quit        bool
}

func NewMonitor(mach *Mach85) *Monitor {
	mon := &Monitor{
		PageLen:     0x10,
		mach:        mach,
		in:          os.Stdin,
		out:         log.New(os.Stdout, "", 0),
		dasm:        NewDisassembler(mach.mem),
		interactive: true,
	}
	return mon
}

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
	case strings.HasPrefix(cmd, "c"):
		err = m.cpu(args)
	case strings.HasPrefix(cmd, "d"):
		err = m.disassemble(args)
	case strings.HasPrefix(cmd, "r"):
		err = m.run(args)
	case strings.HasPrefix(cmd, "q"):
		m.quit = true
		return
	case strings.HasPrefix(cmd, "t"):
		err = m.trace(args)
	case strings.HasPrefix(cmd, "z"):
		err = m.reset(args)
	default:
		err = fmt.Errorf("unknown command: %v", cmd)
	}
	if err != nil {
		m.out.Println(err)
	}
}

func (m *Monitor) cpu(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.out.Println(m.mach.cpu.String())
	return nil
}

func (m *Monitor) disassemble(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.dasm.PC = m.mach.cpu.PC
	for i := 0; i < m.PageLen; i++ {
		m.out.Println(m.dasm.Next().String())
	}
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

func (m *Monitor) run(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	go m.signalHandler()
	m.mach.Run()
	return nil
}

func (m *Monitor) trace(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	if m.mach.cpu.Trace == nil {
		m.mach.cpu.Trace = func(op Operation) {
			m.out.Println(op)
		}
		m.out.Println("tracing enabled")
	} else {
		m.mach.cpu.Trace = nil
		m.out.Println("tracing disabled")
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
