package mach85

import (
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Device interface {
	Service() error
}

type SDLInput interface {
	SDLEvent(sdl.Event) error
}

type Status int

const (
	Init Status = iota
	Halt
	Run
	Break
	Breakpoint
	Trap
)

func (s Status) String() string {
	switch s {
	case Init:
		return "halt"
	case Halt:
		return "halt"
	case Run:
		return "run"
	case Break:
		return "break"
	case Breakpoint:
		return "breakpoint"
	case Trap:
		return "trap"
	}
	return "???"
}

type Mach85 struct {
	Trace       func(op Operation)
	Breakpoints map[uint16]bool
	Memory      *Memory
	Status      Status
	Err         error
	StopOnBreak bool
	QuitOnStop  bool
	cpu         *CPU
	devices     []Device
	inputs      []SDLInput
	dasm        *Disassembler
	start       chan bool
	stop        chan bool
	reset       chan bool
}

func New() *Mach85 {
	mem := NewMemory(NewMemory64())
	cpu := New6510(mem)
	m := &Mach85{
		Memory:      mem,
		cpu:         cpu,
		dasm:        NewDisassembler(mem),
		Breakpoints: map[uint16]bool{},
		devices:     []Device{},
		start:       make(chan bool, 10),
		stop:        make(chan bool, 10),
		reset:       make(chan bool, 10),
	}
	return m
}

func (m *Mach85) Init() error {
	mem64 := m.Memory.Base.(*Memory64)
	if err := mem64.Init(); err != nil {
		log.Fatal(err)
	}

	video, err := NewVideo(m.Memory)
	if err != nil {
		log.Fatalf("unable to create window: %v", err)
	}
	m.AddDevice(video)
	m.AddDevice(NewJiffyClock(m.cpu))
	m.AddInput(NewKeyboard(m))

	m.cpu.PC = m.Memory.Load16(AddrResetVector) - 1
	return nil
}

func (m *Mach85) Run() {
	m.Status = Init
	lastUpdate := time.Now()
	for {
		if m.Status != Init && m.Status != Run && m.QuitOnStop {
			return
		}
		if m.Status != Run {
			<-m.start
			m.cpu.B = false
			m.Err = nil
		}
		m.Status = Run
		if _, ok := m.Breakpoints[m.cpu.PC+1]; ok {
			m.Status = Breakpoint
			continue
		}
		if m.cpu.B {
			if m.StopOnBreak {
				m.Status = Break
				continue
			}
			m.cpu.brk()
		}
		select {
		case <-m.stop:
			m.Status = Halt
			continue
		case <-m.reset:
			m.cpu.Reset()
			mem64 := m.Memory.Base.(*Memory64)
			mem64.Reset()
		default:
			m.cycle()
		}

		now := time.Now()
		if now.Sub(lastUpdate) > time.Millisecond {
			lastUpdate = now
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				if _, ok := event.(*sdl.QuitEvent); ok {
					os.Exit(0)
				}
				for _, input := range m.inputs {
					input.SDLEvent(event)
				}
			}
		}
	}
}

func (m *Mach85) cycle() {
	if m.Trace != nil && !m.cpu.inISR {
		m.dasm.PC = m.cpu.PC
		m.Trace(m.dasm.Next())
	}
	err := m.cpu.Next()
	if err != nil {
		m.Err = err
		m.Status = Trap
		return
	}
	for _, d := range m.devices {
		err := d.Service()
		if err != nil {
			m.Err = err
			m.Status = Trap
			return
		}
	}
}

func (m *Mach85) Start() {
	m.start <- true
}

func (m *Mach85) Stop() {
	m.stop <- true
}

func (m *Mach85) Reset() {
	m.reset <- true
}

func (m *Mach85) AddDevice(d Device) {
	m.devices = append(m.devices, d)
}

func (m *Mach85) AddInput(i SDLInput) {
	m.inputs = append(m.inputs, i)
}
