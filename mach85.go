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

type Mach85 struct {
	Trace       func(op Operation)
	Breakpoints map[uint16]bool
	Memory      *Memory
	cpu         *CPU
	devices     []Device
	inputs      []SDLInput
	dasm        *Disassembler
	stop        chan bool
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
		stop:        make(chan bool),
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
	m.cpu.B = false
	lastUpdate := time.Now()
	for {
		if _, ok := m.Breakpoints[m.cpu.PC+1]; ok {
			return
		}
		if m.cpu.B && m.cpu.StopOnBreak {
			return
		}
		select {
		case <-m.stop:
			return
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
	m.cpu.Next()
	for _, d := range m.devices {
		d.Service()
	}
}

func (m *Mach85) Stop() {
	m.stop <- true
}

func (m *Mach85) Reset() {
	m.cpu.Reset()
}

func (m *Mach85) AddDevice(d Device) {
	m.devices = append(m.devices, d)
}

func (m *Mach85) AddInput(i SDLInput) {
	m.inputs = append(m.inputs, i)
}
