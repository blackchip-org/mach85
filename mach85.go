package mach85

import (
	"fmt"
	"io/ioutil"
)

type Mach85 struct {
	mem *Memory
	cpu *CPU
}

const (
	addrBasic  = 0xa000
	addrKernal = 0xe000
)

func New() *Mach85 {
	mem := NewMemory64k()
	cpu := NewCPU(mem)
	m := &Mach85{
		mem: mem,
		cpu: cpu,
	}
	cpu.PC = mem.Load16(ResetVector) - 1
	return m
}

func (m *Mach85) LoadROM() error {
	basic, err := ioutil.ReadFile("./basic.rom")
	if err != nil {
		return fmt.Errorf("unable to load basic.rom: %v", err)
	}
	m.mem.Import(addrBasic, basic)
	kernal, err := ioutil.ReadFile("./kernal.rom")
	if err != nil {
		return fmt.Errorf("unable to load kernal.rom: %v", err)
	}
	m.mem.Import(addrKernal, kernal)
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
