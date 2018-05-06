package mach85

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Mach85 struct {
	Trace       func(op Operation)
	Breakpoints map[uint16]bool
	Memory      *Memory
	ROMPath     string
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
	{"basic.rom", AddrBasicROM, "79015323128650c742a3694c9429aa91f355905e"},
	//{"chargen.rom", AddrCharacterROM, "adc7c31e18c7c7413d54802ef2f4193da14711aa"},
	{"kernal.rom", AddrKernalROM, "1d503e56df85a62fee696e7618dc5b4e781df1bb"},
}

func New() *Mach85 {
	mem := NewMemory64k()
	cpu := NewCPU(mem)
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
	if err := m.loadROM(); err != nil {
		log.Fatal(err)
	}
	video, err := NewVideo(m.Memory)
	if err != nil {
		log.Fatalf("unable to create window: %v", err)
	}
	m.AddDevice(video)
	m.AddDevice(NewHackDevice(m.Memory))
	return nil
}

func (m *Mach85) loadROM() error {
	for _, rom := range roms {
		file := filepath.Join(m.ROMPath, rom.file)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		checksum := fmt.Sprintf("%x", sha1.Sum(data))
		if checksum != rom.checksum {
			return fmt.Errorf("%v: invalid checksum", rom.file)
		}
		m.Memory.Import(rom.address, data)
	}
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
