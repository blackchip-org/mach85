package mach85

type HackDevice struct {
	mem *Memory
}

func NewHackDevice(m *Memory) *HackDevice {
	return &HackDevice{mem: m}
}

func (d *HackDevice) Service() {
	d.mem.Store(0xd012, 00) // set raster line to zero
}
