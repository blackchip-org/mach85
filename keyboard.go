package mach85

import (
	"flag"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var debugKeyboard bool

func init() {
	flag.BoolVar(&debugKeyboard, "debug-keyboard", false, "log keystroke events")
}

type Keyboard struct {
	mach *Mach85
	mem  *Memory
}

func NewKeyboard(mach *Mach85) *Keyboard {
	return &Keyboard{
		mach: mach,
		mem:  mach.Memory,
	}
}

func (k *Keyboard) SDLEvent(event sdl.Event) error {
	e, ok := event.(*sdl.KeyboardEvent)
	if !ok {
		return nil
	}
	if debugKeyboard {
		fmt.Printf("key: %+v\n", e.Keysym)
	}
	ch, ok := k.lookup(e)
	if !ok {
		return nil
	}
	len := k.mem.Load(AddrKeyboardBufferLen)
	if len >= 10 { // max buffer len
		return nil
	}
	k.mem.Store(AddrKeyboardBuffer+uint16(len), ch)
	len++
	k.mem.Store(AddrKeyboardBufferLen, len)
	return nil
}

// https://wiki.libsdl.org/SDLKeycodeLookup
type keymap map[sdl.Keycode]uint8

var keys = keymap{
	sdl.K_BACKSPACE:    0x14,
	sdl.K_RETURN:       0x0d,
	sdl.K_SPACE:        0x20,
	sdl.K_QUOTE:        0x27,
	sdl.K_PERIOD:       0x2e,
	sdl.K_COMMA:        0x2c,
	sdl.K_SLASH:        0x2f,
	sdl.K_0:            0x30,
	sdl.K_1:            0x31,
	sdl.K_2:            0x32,
	sdl.K_3:            0x33,
	sdl.K_4:            0x34,
	sdl.K_5:            0x35,
	sdl.K_6:            0x36,
	sdl.K_7:            0x37,
	sdl.K_8:            0x38,
	sdl.K_9:            0x39,
	sdl.K_SEMICOLON:    0x3b,
	sdl.K_EQUALS:       0x3d,
	sdl.K_LEFTBRACKET:  0x5b,
	sdl.K_BACKSLASH:    0x5c, // british pound
	sdl.K_RIGHTBRACKET: 0x5d,
	sdl.K_a:            0x41,
	sdl.K_b:            0x42,
	sdl.K_c:            0x43,
	sdl.K_d:            0x44,
	sdl.K_e:            0x45,
	sdl.K_f:            0x46,
	sdl.K_g:            0x47,
	sdl.K_h:            0x48,
	sdl.K_i:            0x49,
	sdl.K_j:            0x4a,
	sdl.K_k:            0x4b,
	sdl.K_l:            0x4c,
	sdl.K_m:            0x4d,
	sdl.K_n:            0x4e,
	sdl.K_o:            0x4f,
	sdl.K_p:            0x50,
	sdl.K_q:            0x51,
	sdl.K_r:            0x52,
	sdl.K_s:            0x53,
	sdl.K_t:            0x54,
	sdl.K_u:            0x55,
	sdl.K_v:            0x56,
	sdl.K_w:            0x57,
	sdl.K_x:            0x58,
	sdl.K_y:            0x59,
	sdl.K_z:            0x5a,
}

var shifted = keymap{
	sdl.K_QUOTE:     0x22,
	sdl.K_PERIOD:    0x3e,
	sdl.K_COMMA:     0x3c,
	sdl.K_SLASH:     0x3f,
	sdl.K_0:         0x29,
	sdl.K_1:         0x21,
	sdl.K_2:         0x40,
	sdl.K_3:         0x23,
	sdl.K_4:         0x24,
	sdl.K_5:         0x25,
	sdl.K_6:         0x5e,
	sdl.K_7:         0x26,
	sdl.K_8:         0x2a,
	sdl.K_9:         0x28,
	sdl.K_SEMICOLON: 0x3a,
	sdl.K_EQUALS:    0x2b,
	sdl.K_a:         0xc1,
	sdl.K_b:         0xc2,
	sdl.K_c:         0xc3,
	sdl.K_d:         0xc4,
	sdl.K_e:         0xc5,
	sdl.K_f:         0xc6,
	sdl.K_g:         0xc7,
	sdl.K_h:         0xc8,
	sdl.K_i:         0xc9,
	sdl.K_j:         0xca,
	sdl.K_k:         0xcb,
	sdl.K_l:         0xcc,
	sdl.K_m:         0xcd,
	sdl.K_n:         0xce,
	sdl.K_o:         0xcf,
	sdl.K_p:         0xd0,
	sdl.K_q:         0xd1,
	sdl.K_r:         0xd2,
	sdl.K_s:         0xd3,
	sdl.K_t:         0xd4,
	sdl.K_u:         0xd5,
	sdl.K_v:         0xd6,
	sdl.K_w:         0xd7,
	sdl.K_x:         0xd8,
	sdl.K_y:         0xd9,
	sdl.K_z:         0xda,
}

var keymaps = map[sdl.Keymod]keymap{
	sdl.KMOD_NONE:   keys,
	sdl.KMOD_LSHIFT: shifted,
	sdl.KMOD_RSHIFT: shifted,
}

func (k *Keyboard) lookup(e *sdl.KeyboardEvent) (uint8, bool) {
	keysym := e.Keysym
	switch {
	case keysym.Mod&sdl.KMOD_CTRL > 0 && keysym.Sym == sdl.K_ESCAPE:
		k.mach.Reset()
		return 0, false
	case keysym.Mod&sdl.KMOD_CTRL > 0 && keysym.Sym == sdl.K_BACKSPACE:
		if e.Type == sdl.KEYDOWN {
			k.mem.Store(AddrStopKey, 0x7f)
		} else if e.Type == sdl.KEYUP {
			k.mem.Store(AddrStopKey, 0xff)
		}
	}

	if e.Type != sdl.KEYDOWN {
		return 0, false
	}
	keymap0 := keymaps[sdl.KMOD_NONE]
	keymap, ok := keymaps[sdl.Keymod(keysym.Mod)]
	if !ok {
		keymap = keymap0
	}
	ch, ok := keymap[keysym.Sym]
	if !ok {
		ch, ok = keymap0[keysym.Sym]
	}
	if !ok {
		return 0, false
	}
	return ch, true
}

func (k *Keyboard) special(keysym sdl.Keysym) bool {
	return false
}
