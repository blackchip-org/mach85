package mach85

import "time"

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
	if now.Sub(c.lastUpdate) < 16800000 { // 16.8 ms, NTSC
		return nil
	}
	c.lastUpdate = now
	c.cpu.IRQ()
	return nil
}
