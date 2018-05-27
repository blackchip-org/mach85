package mach85

import "fmt"

type Watchdog struct {
	mach     *Mach85
	lastPC   uint16
	pcRepeat int
}

func NewWatchdog(mach *Mach85) *Watchdog {
	return &Watchdog{mach: mach}
}

func (w *Watchdog) Service() error {
	if w.lastPC == w.mach.cpu.PC {
		w.pcRepeat++
		if w.pcRepeat == 3 {
			return fmt.Errorf("loop")
		}
	} else {
		w.pcRepeat = 0
	}
	w.lastPC = w.mach.cpu.PC
	return nil
}
