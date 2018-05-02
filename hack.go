package mach85

type hackDevice struct {
	mem *Memory
}

func newHackDevice(m *Memory) *hackDevice {
	return &hackDevice{mem: m}
}

func (d *hackDevice) Service() {
	d.mem.Store(0xd012, 00) // set raster line to zero
}
