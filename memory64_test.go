package mach85

import (
	"fmt"
	"testing"
)

func TestMemoryLoad(t *testing.T) {
	const (
		page00 = 0x0040
		page10 = 0x1040
		page80 = 0x8040
		pagea0 = 0xa040
		pagec0 = 0xc040
		paged0 = 0xd040
		pagee0 = 0xe040

		ram0   = 0xf0
		ram1   = 0xf1
		ram2   = 0xf2
		ram3   = 0xf3
		ram4   = 0xf4
		ram5   = 0xf5
		ram6   = 0xf6
		basic  = 0xba
		kernal = 0xea
		char   = 0xca
		cartLo = 0x0c
		cartHi = 0xc0
		io     = 0x10
		open   = 0x00
	)

	var tests = []struct {
		mode    uint8
		address uint16
		want    uint8
	}{
		{0, page00, ram0},
		{0, page10, ram1},
		{0, page80, ram2},
		{0, pagea0, ram3},
		{0, pagec0, ram4},
		{0, paged0, ram5},
		{0, pagee0, ram6},

		{1, page00, ram0},
		{1, page10, ram1},
		{1, page80, ram2},
		{1, pagea0, ram3},
		{1, pagec0, ram4},
		{1, paged0, ram5},
		{1, pagee0, ram6},

		{2, page00, ram0},
		{2, page10, ram1},
		{2, page80, ram2},
		{2, pagea0, cartHi},
		{2, pagec0, ram4},
		{2, paged0, char},
		{2, pagee0, kernal},

		{3, page00, ram0},
		{3, page10, ram1},
		{3, page80, cartLo},
		{3, pagea0, cartHi},
		{3, pagec0, ram4},
		{3, paged0, char},
		{3, pagee0, kernal},

		{4, page00, ram0},
		{4, page10, ram1},
		{4, page80, ram2},
		{4, pagea0, ram3},
		{4, pagec0, ram4},
		{4, paged0, ram5},
		{4, pagee0, ram6},

		{5, page00, ram0},
		{5, page10, ram1},
		{5, page80, ram2},
		{5, pagea0, ram3},
		{5, pagec0, ram4},
		{5, paged0, io},
		{5, pagee0, ram6},

		{6, page00, ram0},
		{6, page10, ram1},
		{6, page80, ram2},
		{6, pagea0, cartHi},
		{6, pagec0, ram4},
		{6, paged0, io},
		{6, pagee0, kernal},

		{7, page00, ram0},
		{7, page10, ram1},
		{7, page80, cartLo},
		{7, pagea0, cartHi},
		{7, pagec0, ram4},
		{7, paged0, io},
		{7, pagee0, kernal},

		{8, page00, ram0},
		{8, page10, ram1},
		{8, page80, ram2},
		{8, pagea0, ram3},
		{8, pagec0, ram4},
		{8, paged0, ram5},
		{8, pagee0, ram6},

		{9, page00, ram0},
		{9, page10, ram1},
		{9, page80, ram2},
		{9, pagea0, ram3},
		{9, pagec0, ram4},
		{9, paged0, char},
		{9, pagee0, ram6},

		{10, page00, ram0},
		{10, page10, ram1},
		{10, page80, ram2},
		{10, pagea0, ram3},
		{10, pagec0, ram4},
		{10, paged0, char},
		{10, pagee0, kernal},

		{11, page00, ram0},
		{11, page10, ram1},
		{11, page80, cartLo},
		{11, pagea0, basic},
		{11, pagec0, ram4},
		{11, paged0, char},
		{11, pagee0, kernal},

		{12, page00, ram0},
		{12, page10, ram1},
		{12, page80, ram2},
		{12, pagea0, ram3},
		{12, pagec0, ram4},
		{12, paged0, ram5},
		{12, pagee0, ram6},

		{13, page00, ram0},
		{13, page10, ram1},
		{13, page80, ram2},
		{13, pagea0, ram3},
		{13, pagec0, ram4},
		{13, paged0, io},
		{13, pagee0, ram6},

		{14, page00, ram0},
		{14, page10, ram1},
		{14, page80, ram2},
		{14, pagea0, ram3},
		{14, pagec0, ram4},
		{14, paged0, io},
		{14, pagee0, kernal},

		{15, page00, ram0},
		{15, page10, ram1},
		{15, page80, cartLo},
		{15, pagea0, basic},
		{15, pagec0, ram4},
		{15, paged0, io},
		{15, pagee0, kernal},

		{16, page00, ram0},
		{16, page10, open},
		{16, page80, cartLo},
		{16, pagea0, open},
		{16, pagec0, open},
		{16, paged0, io},
		{16, pagee0, cartHi},

		{17, page00, ram0},
		{17, page10, open},
		{17, page80, cartLo},
		{17, pagea0, open},
		{17, pagec0, open},
		{17, paged0, io},
		{17, pagee0, cartHi},

		{18, page00, ram0},
		{18, page10, open},
		{18, page80, cartLo},
		{18, pagea0, open},
		{18, pagec0, open},
		{18, paged0, io},
		{18, pagee0, cartHi},

		{19, page00, ram0},
		{19, page10, open},
		{19, page80, cartLo},
		{19, pagea0, open},
		{19, pagec0, open},
		{19, paged0, io},
		{19, pagee0, cartHi},

		{20, page00, ram0},
		{20, page10, open},
		{20, page80, cartLo},
		{20, pagea0, open},
		{20, pagec0, open},
		{20, paged0, io},
		{20, pagee0, cartHi},

		{21, page00, ram0},
		{21, page10, open},
		{21, page80, cartLo},
		{21, pagea0, open},
		{21, pagec0, open},
		{21, paged0, io},
		{21, pagee0, cartHi},

		{22, page00, ram0},
		{22, page10, open},
		{22, page80, cartLo},
		{22, pagea0, open},
		{22, pagec0, open},
		{22, paged0, io},
		{22, pagee0, cartHi},

		{23, page00, ram0},
		{23, page10, open},
		{23, page80, cartLo},
		{23, pagea0, open},
		{23, pagec0, open},
		{23, paged0, io},
		{23, pagee0, cartHi},

		{24, page00, ram0},
		{24, page10, ram1},
		{24, page80, ram2},
		{24, pagea0, ram3},
		{24, pagec0, ram4},
		{24, paged0, ram5},
		{24, pagee0, ram6},

		{25, page00, ram0},
		{25, page10, ram1},
		{25, page80, ram2},
		{25, pagea0, ram3},
		{25, pagec0, ram4},
		{25, paged0, char},
		{25, pagee0, ram6},

		{26, page00, ram0},
		{26, page10, ram1},
		{26, page80, ram2},
		{26, pagea0, ram3},
		{26, pagec0, ram4},
		{26, paged0, char},
		{26, pagee0, kernal},

		{27, page00, ram0},
		{27, page10, ram1},
		{27, page80, ram2},
		{27, pagea0, basic},
		{27, pagec0, ram4},
		{27, paged0, char},
		{27, pagee0, kernal},

		{28, page00, ram0},
		{28, page10, ram1},
		{28, page80, ram2},
		{28, pagea0, ram3},
		{28, pagec0, ram4},
		{28, paged0, ram5},
		{28, pagee0, ram6},

		{29, page00, ram0},
		{29, page10, ram1},
		{29, page80, ram2},
		{29, pagea0, ram3},
		{29, pagec0, ram4},
		{29, paged0, io},
		{29, pagee0, ram6},

		{30, page00, ram0},
		{30, page10, ram1},
		{30, page80, ram2},
		{30, pagea0, ram3},
		{30, pagec0, ram4},
		{30, paged0, io},
		{30, pagee0, kernal},

		{31, page00, ram0},
		{31, page10, ram1},
		{31, page80, ram2},
		{31, pagea0, basic},
		{31, pagec0, ram4},
		{31, paged0, io},
		{31, pagee0, kernal},
	}

	for _, test := range tests {
		label := fmt.Sprintf("mode %02d address %04x", test.mode, test.address)
		t.Run(label, func(t *testing.T) {
			mem := NewMemory64()
			mem.Chunks[RAM0].Store(0x40, ram0)
			mem.Chunks[RAM1].Store(0x40, ram1)
			mem.Chunks[RAM2].Store(0x40, ram2)
			mem.Chunks[RAM3].Store(0x40, ram3)
			mem.Chunks[RAM4].Store(0x40, ram4)
			mem.Chunks[RAM5].Store(0x40, ram5)
			mem.Chunks[RAM6].Store(0x40, ram6)
			mem.Chunks[BasicROM] = NullMemory{Value: basic}
			mem.Chunks[KernalROM] = NullMemory{Value: kernal}
			mem.Chunks[CharROM] = NullMemory{Value: char}
			mem.Chunks[CartLoROM] = NullMemory{Value: cartLo}
			mem.Chunks[CartHiROM] = NullMemory{Value: cartHi}
			mem.Chunks[IO].Store(0x40, io)

			mem.SetMode(test.mode)
			have := mem.Load(test.address)
			if test.want != have {
				t.Errorf("\n want: %02x \n have: %02x", test.want, have)
			}
		})
	}
}

func TestMemoryStoreRAM(t *testing.T) {
	mode := 0
	mem := NewMemory64()
	mem.SetMode(uint8(mode))
	for address := 2; address <= 0xffff; address++ {
		mem.Store(uint16(address), uint8(address))
	}
	mem.SetMode(0)
	for address := 2; address <= 0xffff; address++ {
		have := mem.Load(uint16(address))
		want := uint8(address)
		if want != have {
			t.Fatalf("at %04x \n want: %02x \n have: %02x \n", address, want, have)
		}
	}
}

func TestMemoryStoreThroughROM(t *testing.T) {
	var tests = []struct {
		mode    uint8
		address uint16
		label   string
	}{
		{03, 0x8000, "cart lo"},
		{03, 0xa000, "cart hi"},
		{02, 0xd000, "char"},
		{31, 0xa000, "basic"},
		{21, 0xe000, "kernal"},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			mem := NewMemory64()
			mem.SetMode(test.mode)
			want := uint8(0xab)
			mem.Store(test.address, want)
			mem.SetMode(0)
			have := mem.Load(test.address)
			if want != have {
				t.Fatalf("\n want: %02x \n have: %02x \n", want, have)
			}
		})
	}
}
