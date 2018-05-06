package mach85

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/blackchip-org/mach85/rom"
)

// https://www.c64-wiki.com/wiki/Bank_Switching

type Chunk int

const (
	RAM0 Chunk = iota
	RAM1
	RAM2
	RAM3
	RAM4
	RAM5
	RAM6
	BasicROM
	KernalROM
	CharROM
	CartLoROM
	CartHiROM
	IO
	Open
)

var addrZones = [7]uint16{
	0x0000,
	0x1000,
	0x8000,
	0xa000,
	0xc000,
	0xd000,
	0xe000,
}

var zoneMap = [16]int{
	/* $0 */ 0,
	/* $1 */ 1,
	/* $2 */ 1,
	/* $3 */ 1,
	/* $4 */ 1,
	/* $5 */ 1,
	/* $6 */ 1,
	/* $7 */ 1,
	/* $8 */ 2,
	/* $9 */ 2,
	/* $a */ 3,
	/* $b */ 3,
	/* $c */ 4,
	/* $d */ 5,
	/* $e */ 6,
	/* $f */ 6,
}

type zones [7]Chunk

var modes = [32]zones{
	/* 00 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, RAM5, RAM6},
	/* 01 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, RAM5, RAM6},
	/* 02 */ zones{RAM0, RAM1, RAM2, CartHiROM, RAM4, CharROM, KernalROM},
	/* 03 */ zones{RAM0, RAM1, CartLoROM, CartHiROM, RAM4, CharROM, KernalROM},
	/* 04 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, RAM5, RAM6},
	/* 05 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, IO, RAM6},
	/* 06 */ zones{RAM0, RAM1, RAM2, CartHiROM, RAM4, IO, KernalROM},
	/* 07 */ zones{RAM0, RAM1, CartLoROM, CartHiROM, RAM4, IO, KernalROM},
	/* 08 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, RAM5, RAM6},
	/* 09 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, CharROM, RAM6},
	/* 10 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, CharROM, KernalROM},
	/* 11 */ zones{RAM0, RAM1, CartLoROM, BasicROM, RAM4, CharROM, KernalROM},
	/* 12 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, RAM5, RAM6},
	/* 13 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, IO, RAM6},
	/* 14 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, IO, KernalROM},
	/* 15 */ zones{RAM0, RAM1, CartLoROM, BasicROM, RAM4, IO, KernalROM},
	/* 16 */ zones{RAM0, Open, CartLoROM, Open, Open, IO, CartHiROM},
	/* 17 */ zones{RAM0, Open, CartLoROM, Open, Open, IO, CartHiROM},
	/* 18 */ zones{RAM0, Open, CartLoROM, Open, Open, IO, CartHiROM},
	/* 19 */ zones{RAM0, Open, CartLoROM, Open, Open, IO, CartHiROM},
	/* 20 */ zones{RAM0, Open, CartLoROM, Open, Open, IO, CartHiROM},
	/* 21 */ zones{RAM0, Open, CartLoROM, Open, Open, IO, CartHiROM},
	/* 22 */ zones{RAM0, Open, CartLoROM, Open, Open, IO, CartHiROM},
	/* 23 */ zones{RAM0, Open, CartLoROM, Open, Open, IO, CartHiROM},
	/* 24 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, RAM5, RAM6},
	/* 25 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, CharROM, RAM6},
	/* 26 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, CharROM, KernalROM},
	/* 27 */ zones{RAM0, RAM1, RAM2, BasicROM, RAM4, CharROM, KernalROM},
	/* 28 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, RAM5, RAM6},
	/* 29 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, IO, RAM6},
	/* 30 */ zones{RAM0, RAM1, RAM2, RAM3, RAM4, IO, KernalROM},
	/* 31 */ zones{RAM0, RAM1, RAM2, BasicROM, RAM4, IO, KernalROM},
}

type BankedMemory struct {
	Chunks [14]MemoryChunk
	Game   bool // pin 8
	ExROM  bool // pin 9
}

func NewBankedMemory() *BankedMemory {
	m := &BankedMemory{}

	m.Chunks[RAM0] = NewRAM(0x1000) // $0000 - $0fff
	m.Chunks[RAM1] = NewRAM(0x7000) // $1000 - $7fff
	m.Chunks[RAM2] = NewRAM(0x2000) // $8000 - $9fff
	m.Chunks[RAM3] = NewRAM(0x2000) // $a000 - $bfff
	m.Chunks[RAM4] = NewRAM(0x1000) // $c000 - $cfff
	m.Chunks[RAM5] = NewRAM(0x1000) // $d000 - $dfff
	m.Chunks[RAM6] = NewRAM(0x2000) // $e000 - $ffff

	m.Chunks[BasicROM] = NullMemory{}
	m.Chunks[KernalROM] = NullMemory{}
	m.Chunks[CharROM] = NullMemory{}
	m.Chunks[CartLoROM] = NullMemory{}
	m.Chunks[CartHiROM] = NullMemory{}
	m.Chunks[Open] = NullMemory{}

	m.Chunks[IO] = NewRAM(0x1000) // $d000 - $dfff

	m.SetMode(31)
	return m
}

func (m *BankedMemory) Mode() uint8 {
	mode := m.Chunks[RAM0].Load(AddrR6510) & 0x7 // bits 0 - 2
	if m.Game {
		mode |= 0x8 // bit 3
	}
	if m.ExROM {
		mode |= 0x10 // bit 4
	}
	return mode
}

func (m *BankedMemory) SetMode(value uint8) {
	prev := m.Chunks[RAM0].Load(AddrR6510)
	next := prev&0xf8 + value&0x07 // bits 0 - 2
	m.Chunks[RAM0].Store(AddrR6510, next)
	m.Game = value&0x08 > 0  // bit 3
	m.ExROM = value&0x10 > 0 // bit 4
}

func (m *BankedMemory) Load(address uint16) uint8 {
	zones := modes[m.Mode()]
	zone := zoneMap[address>>12]
	chunk := m.Chunks[zones[zone]]
	return chunk.Load(address - addrZones[zone])
}

func (m *BankedMemory) Store(address uint16, value uint8) {
	zones := modes[0]
	zone := zoneMap[address>>12]
	chunk := m.Chunks[zones[zone]]
	chunk.Store(address-addrZones[zone], value)
}

var roms = []struct {
	file     string
	chunk    Chunk
	checksum string
}{
	{"basic.rom", BasicROM, "79015323128650c742a3694c9429aa91f355905e"},
	{"chargen.rom", CharROM, "adc7c31e18c7c7413d54802ef2f4193da14711aa"},
	{"kernal.rom", KernalROM, "1d503e56df85a62fee696e7618dc5b4e781df1bb"},
}

func (m *BankedMemory) Init() error {
	for _, r := range roms {
		file := filepath.Join(rom.Path, r.file)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		checksum := fmt.Sprintf("%x", sha1.Sum(data))
		if checksum != r.checksum {
			return fmt.Errorf("%v: invalid checksum", r.file)
		}
		m.Chunks[r.chunk] = NewROM(data)
	}
	return nil
}
