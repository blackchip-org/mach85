package mach85

import (
	"testing"

	"github.com/blackchip-org/mach85/rom"
	"github.com/veandco/go-sdl2/sdl"
)

func BenchmarkMach(b *testing.B) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		b.Fatalf("unable to initialize sdl: %v", err)
	}
	defer sdl.Quit()
	sdl.GLSetSwapInterval(1)
	mach := New()
	rom.Path = "rom"
	if err := mach.Init(); err != nil {
		b.Fatalf("unable to initialize: %v", err)
	}
	b.Run("BenchmarkMach", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			mach.cycle()
		}
	})
}
