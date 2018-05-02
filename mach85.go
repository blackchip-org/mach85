package mach85

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
)

type Mach85 struct {
	Trace       func(op Operation)
	Breakpoints map[uint16]bool
	mem         *Memory
	cpu         *CPU
	devices     []Device
	dasm        *Disassembler
	stop        chan bool
}

type Device interface {
	Service()
}

var roms = []struct {
	file     string
	address  uint16
	checksum string
}{
	{"basic.rom", 0xa000, "79015323128650c742a3694c9429aa91f355905e"},
	{"kernal.rom", 0xe000, "1d503e56df85a62fee696e7618dc5b4e781df1bb"},
}

func New() *Mach85 {
	mem := NewMemory64k()
	cpu := NewCPU(mem)
	m := &Mach85{
		mem:         mem,
		cpu:         cpu,
		dasm:        NewDisassembler(mem),
		stop:        make(chan bool),
		Breakpoints: map[uint16]bool{},
		devices: []Device{
			newHackDevice(mem),
		},
	}
	return m
}

func (m *Mach85) LoadROM() error {
	for _, rom := range roms {
		data, err := ioutil.ReadFile(rom.file)
		if err != nil {
			return err
		}
		checksum := fmt.Sprintf("%x", sha1.Sum(data))
		if checksum != rom.checksum {
			return fmt.Errorf("%v: invalid checksum", rom.file)
		}
		m.mem.Import(rom.address, data)
		log.Printf("$%04x: %v\n", rom.address, rom.file)
	}
	m.cpu.PC = m.mem.Load16(ResetVector) - 1
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
			if m.Trace != nil {
				m.dasm.PC = m.cpu.PC
				m.Trace(m.dasm.Next())
			}
			m.cpu.Next()
		}
		for _, d := range m.devices {
			d.Service()
		}
	}
}

func (m *Mach85) Stop() {
	m.stop <- true
}

func (m *Mach85) Reset() {
	m.cpu.Reset()
	m.Run()
}
