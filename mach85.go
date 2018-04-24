package mach85

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
)

type Mach85 struct {
	mem *Memory
	cpu *CPU
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
	m.cpu.Run()
}

func (m *Mach85) Stop() {
	m.cpu.Stop()
}

func (m *Mach85) Reset() {
	m.cpu.Reset()
	m.Run()
}
