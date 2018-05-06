package mach85

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width   = 320
	height  = 200
	screenW = 404 // actually 403?
	screenH = 284
	borderW = (screenW - width) / 2
	borderH = (screenH - height) / 2
)

var (
	Black      = []uint8{0x00, 0x00, 0x00, 0xff}
	White      = []uint8{0xff, 0xff, 0xff, 0xff}
	Red        = []uint8{0x88, 0x00, 0x00, 0xff}
	Cyan       = []uint8{0xaa, 0xff, 0xee, 0xff}
	Purple     = []uint8{0xcc, 0x44, 0xcc, 0xff}
	Green      = []uint8{0x00, 0xcc, 0x55, 0xff}
	Blue       = []uint8{0x00, 0x00, 0xaa, 0xff}
	Yellow     = []uint8{0xee, 0xee, 0x77, 0xff}
	Orange     = []uint8{0xdd, 0x88, 0x55, 0xff}
	Brown      = []uint8{0x66, 0x44, 0x00, 0xff}
	LightRed   = []uint8{0xff, 0x77, 0x77, 0xff}
	DarkGray   = []uint8{0x33, 0x33, 0x33, 0xff}
	Gray       = []uint8{0x77, 0x77, 0x77, 0xff}
	LightGreen = []uint8{0xaa, 0xff, 0x66, 0xff}
	LightBlue  = []uint8{0x00, 0x88, 0xff, 0xff}
	LightGray  = []uint8{0xbb, 0xbb, 0xbb, 0xff}
)

var colorMap = [...][]uint8{
	Black,
	White,
	Red,
	Cyan,
	Purple,
	Green,
	Blue,
	Yellow,
	Orange,
	Brown,
	LightRed,
	DarkGray,
	Gray,
	LightGreen,
	LightBlue,
	LightGray,
}

type Video struct {
	mem        *Memory
	window     *sdl.Window
	renderer   *sdl.Renderer
	lastUpdate time.Time
}

func NewVideo(mem *Memory) (*Video, error) {
	window, err := sdl.CreateWindow(
		"mach85",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(screenW), int32(screenH),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, err
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	return &Video{
		mem:      mem,
		window:   window,
		renderer: renderer,
	}, nil
}

func (v *Video) Service() {
	now := time.Now()
	if now.Sub(v.lastUpdate) < time.Millisecond*16 {
		return
	}
	v.lastUpdate = now
	v.drawBorder()
	v.drawBackground()
	v.renderer.Present()
}

func (v *Video) drawBorder() {
	index := v.mem.Load(AddrBorderColor) & 0xf
	borderColor := []uint8(colorMap[index])
	v.renderer.SetDrawColorArray(borderColor...)
	topBorder := sdl.Rect{
		X: 0,
		Y: 0,
		W: screenW,
		H: borderH,
	}
	v.renderer.FillRect(&topBorder)
	bottomBorder := sdl.Rect{
		X: 0,
		Y: borderH + height,
		W: screenW,
		H: borderH,
	}
	v.renderer.FillRect(&bottomBorder)
	leftBorder := sdl.Rect{
		X: 0,
		Y: borderH,
		W: borderW,
		H: height,
	}
	v.renderer.FillRect(&leftBorder)
	rightBorder := sdl.Rect{
		X: borderW + width,
		Y: borderH,
		W: borderW,
		H: height,
	}
	v.renderer.FillRect(&rightBorder)
}

func (v *Video) drawBackground() {
	index := v.mem.Load(AddrBackgroundColor) & 0xf
	backgroundColor := []uint8(colorMap[index])
	v.renderer.SetDrawColorArray(backgroundColor...)
	background := sdl.Rect{
		X: borderW,
		Y: borderH,
		W: width,
		H: height,
	}
	v.renderer.FillRect(&background)
}
