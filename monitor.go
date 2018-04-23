package mach85

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
)

type Monitor struct {
	PageLen int
	mach    *Mach85
	in      io.Reader
	out     io.Writer
	dasm    *Disassembler
	quit    bool
}

func NewMonitor(mach *Mach85) *Monitor {
	return &Monitor{
		PageLen: 0x10,
		mach:    mach,
		in:      os.Stdin,
		out:     os.Stdout,
		dasm:    NewDisassembler(mach.mem),
	}
}

func (m *Monitor) Run() {
	s := bufio.NewScanner(m.in)
	s.Split(bufio.ScanLines)
	for {
		// only show prompt when printing to the console
		if m.out == os.Stdout {
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
	s := bufio.NewScanner(strings.NewReader(line))
	s.Split(bufio.ScanWords)
	if !s.Scan() {
		return
	}
	cmd := s.Text()
	switch {
	case strings.HasPrefix(cmd, "c"):
		m.cpu(s)
	case strings.HasPrefix(cmd, "d"):
		m.disassemble(s)
	case strings.HasPrefix(cmd, "r"):
		m.run(s)
	case strings.HasPrefix(cmd, "q"):
		m.quit = true
		return
	case strings.HasPrefix(cmd, "z"):
		m.reset(s)
	default:
		m.print("unknown command: %v\n", cmd)
	}
}

func (m *Monitor) cpu(s *bufio.Scanner) {
	if s.Scan() {
		m.print("too many arguments\n")
		return
	}
	m.print(m.mach.cpu.String() + "\n")
}

func (m *Monitor) disassemble(s *bufio.Scanner) {
	if s.Scan() {
		m.print("too many arguments\n")
		return
	}
	m.dasm.PC = m.mach.cpu.PC
	for i := 0; i < m.PageLen; i++ {
		m.print("%v\n", m.dasm.Next().String())
	}
}

func (m *Monitor) reset(s *bufio.Scanner) {
	if s.Scan() {
		m.print("too many arguments\n")
		return
	}
	go m.signalHandler()
	m.mach.Reset()
}

func (m *Monitor) run(s *bufio.Scanner) {
	if s.Scan() {
		m.print("too many arguments\n")
		return
	}
	go m.signalHandler()
	m.mach.Run()
}

func (m *Monitor) print(format string, args ...interface{}) {
	fmt.Fprintf(m.out, format, args...)
}

func (m *Monitor) signalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	signal.Reset(os.Interrupt)
	m.mach.Stop()
}
