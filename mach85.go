package mach85

import (
	"log"
	"time"
)

type Device interface {
	Service() error
}

type Mach85 struct {
	Trace       func(op Operation)
	Breakpoints map[uint16]bool
	Memory      *Memory
	cpu         *CPU
	devices     []Device
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
		stop:        make(chan bool),
		Breakpoints: map[uint16]bool{},
		devices:     []Device{},
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
	m.cpu.PC = m.Memory.Load16(AddrResetVector) - 1
	return nil
}

func (m *Mach85) Run() {
	m.cpu.B = false
	for {
		if _, ok := m.Breakpoints[m.cpu.PC+1]; ok {
			return
		}
		if m.cpu.B {
			return
		}
		select {
		case <-m.stop:
			return
		default:
			m.cycle()
		}
	}
}

func (m *Mach85) cycle() {
	if m.Trace != nil {
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
	m.Run()
}

func (m *Mach85) AddDevice(d Device) {
	m.devices = append(m.devices, d)
}

// https://www.c64-wiki.com/wiki/Jiffy_Clock

type JiffyClock struct {
	lastUpdate time.Time
	cpu        *CPU
}

func NewJiffyClock(cpu *CPU) *JiffyClock {
	return &JiffyClock{cpu: cpu}
}

func (c *JiffyClock) Service() error {
	now := time.Now()
	if now.Sub(c.lastUpdate) < 16800000 { // 16.8 ms
		return nil
	}
	c.lastUpdate = now
	c.cpu.IRQ()
	return nil
}
