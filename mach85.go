package mach85

import (
	"fmt"
	"io/ioutil"
)

type Mach85 struct {
	mem *Memory
	cpu *CPU
}

func New() (*Mach85, error) {
	mem := NewMemory64k()
	cpu := NewCPU(mem)
	m := &Mach85{
		mem: mem,
		cpu: cpu,
	}
	basic, err := ioutil.ReadFile("./basic.rom")
	if err != nil {
		return nil, fmt.Errorf("unable to load basic.rom: %v", err)
	}
	mem.Import(0xa000, basic)
	kernal, err := ioutil.ReadFile("./kernal.rom")
	if err != nil {
		return nil, fmt.Errorf("unable to load kernal.rom: %v", err)
	}
	mem.Import(0xe000, kernal)
	return m, nil
}

func (m *Mach85) Run() {
	m.cpu.Trace = true
	m.cpu.Reset()
	m.cpu.Run()
}
