package mach85

import (
	"io/ioutil"
	"log"
)

type Mach85 struct {
	mem *Memory
	cpu *CPU
}

var roms = []struct {
	file    string
	address uint16
}{
	{"basic.rom", 0xa000},
	{"kernal.rom", 0xe000},
}

func New() *Mach85 {
	mem := NewMemory64k()
	cpu := NewCPU(mem)
	m := &Mach85{
		mem: mem,
		cpu: cpu,
	}
	return m
}

func (m *Mach85) LoadROM() error {
	for _, rom := range roms {
		data, err := ioutil.ReadFile(rom.file)
		if err != nil {
			return err
		}
		m.mem.Import(rom.address, data)
		log.Printf("$%04x: %v\n", rom.address, rom.file)
	}
	m.cpu.PC = m.mem.Load16(ResetVector) - 1
	return nil
}

func (m *Mach85) Run() {
	m.cpu.B = false
	m.cpu.Run()
}

func (m *Mach85) Stop() {
	m.cpu.Stop()
}

func (m *Mach85) Reset() {
	m.cpu.Reset()
	m.Run()
}
