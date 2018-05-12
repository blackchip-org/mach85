package mach85

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Keyboard struct {
	mem   *Memory
	debug bool
}

func NewKeyboard(mem *Memory) *Keyboard {
	return &Keyboard{
		mem:   mem,
		debug: true,
	}
}

func (k *Keyboard) SDLEvent(event sdl.Event) error {
	e, ok := event.(*sdl.KeyboardEvent)
	if !ok {
		return nil
	}
	if e.Type != sdl.KEYUP {
		return nil
	}
	if k.debug {
		fmt.Printf("key: %+v\n", e)
	}
	if e.Keysym.Sym > 0xff { // ignore other keys for now
		return nil
	}
	keycode, _ := PetsciiEncoder(rune(e.Keysym.Sym))
	len := k.mem.Load(AddrKeyboardBufferLen)
	if len >= 10 { // max buffer len
		return nil
	}
	k.mem.Store(AddrKeyboardBuffer+uint16(len), keycode)
	len++
	k.mem.Store(AddrKeyboardBufferLen, len)
	return nil
}
