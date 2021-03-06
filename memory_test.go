package mach85

import (
	"strings"
	"testing"
)

func TestLoad16(t *testing.T) {
	m := NewMemory(NewRAM(0x10000))
	m.Store(0x00, 0xcd)
	m.Store(0x01, 0xab)
	want := uint16(0xabcd)
	have := m.Load16(0)
	if want != have {
		t.Errorf("\n want: %x \n have: %x \n", want, have)
	}
}

func TestStore16(t *testing.T) {
	m := NewMemory(NewRAM(0x10000))
	m.Store16(0x00, 0xabcd)
	want := uint8(0xcd)
	have := m.Load(0)
	if want != have {
		t.Errorf("\n want: %x \n have: %x \n", want, have)
	}
	want = uint8(0xab)
	have = m.Load(1)
	if want != have {
		t.Errorf("\n want: %x \n have: %x \n", want, have)
	}
}

var dumpTests = []struct {
	name     string
	start    int
	data     func() []int
	showFrom int
	showTo   int
	want     string
}{
	{
		"one line", 0x10,
		func() []int { return []int{} },
		0x10, 0x20, "" +
			"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
	}, {
		"two lines", 0x10,
		func() []int { return []int{} },
		0x10, 0x30, "" +
			"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................\n" +
			"$0020 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
	}, {
		"jagged top", 0x10,
		func() []int { return []int{} },
		0x14, 0x30, "" +
			"$0010             00 00 00 00  00 00 00 00 00 00 00 00     ............\n" +
			"$0020 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
	}, {
		"jagged bottom", 0x10,
		func() []int { return []int{} },
		0x10, 0x2b, "" +
			"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................\n" +
			"$0020 00 00 00 00 00 00 00 00  00 00 00 00             ............",
	},
	{
		"single value", 0x10,
		func() []int { return []int{0, 0x41} },
		0x11, 0x11, "" +
			"$0010    41                                             A",
	},
	{
		"$40-$5f", 0x10,
		func() []int {
			data := make([]int, 0)
			for i := 0x40; i < 0x60; i++ {
				data = append(data, i)
			}
			return data
		},
		0x10, 0x30, "" +
			"$0010 40 41 42 43 44 45 46 47  48 49 4a 4b 4c 4d 4e 4f @ABCDEFGHIJKLMNO\n" +
			"$0020 50 51 52 53 54 55 56 57  58 59 5a 5b 5c 5d 5e 5f PQRSTUVWXYZ[£]↑←",
	},
}

func TestDump(t *testing.T) {
	m := NewMemory(NewRAM(0x10000))
	for _, test := range dumpTests {
		t.Run(test.name, func(t *testing.T) {
			for i, value := range test.data() {
				m.Store(uint16(test.start+i), uint8(value))
			}
			have := m.Dump(uint16(test.showFrom), uint16(test.showTo),
				PetsciiUnshiftedDecoder)
			have = strings.TrimSpace(have)
			if test.want != have {
				t.Errorf("\n want: \n%v\n have:\n%v", test.want, have)
			}
		})
	}
}
