package mach85

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width   = 320
	height  = 200
	screenW = 403
	screenH = 284
	borderW = (screenW - width) / 2
	borderH = (screenH - height) / 2
)

const (
	Black      = 0xff000000
	White      = 0xffffffff
	Red        = 0xff880000
	Cyan       = 0xffaaffee
	Purple     = 0xffcc44cc
	Green      = 0xff00cc55
	Blue       = 0xff0000aa
	Yellow     = 0xffeeee77
	Orange     = 0xffdd8855
	Brown      = 0xff664400
	LightRed   = 0xffff7777
	DarkGray   = 0xff333333
	Gray       = 0xff777777
	LightGreen = 0xffaaff66
	LightBlue  = 0xff0088ff
	LightGray  = 0xffbbbbbb
)

var colorMap = [...]uint{
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
	surface    *sdl.Surface
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
	surface, err := window.GetSurface()
	if err != nil {
		return nil, err
	}
	return &Video{
		mem:     mem,
		window:  window,
		surface: surface,
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
	v.window.UpdateSurface()
}

func (v *Video) drawBorder() {
	index := v.mem.Load(AddrBorderColor) & 0xf
	borderColor := uint32(colorMap[index])
	topBorder := sdl.Rect{
		X: 0,
		Y: 0,
		W: screenW,
		H: borderH,
	}
	v.surface.FillRect(&topBorder, borderColor)
	bottomBorder := sdl.Rect{
		X: 0,
		Y: borderH + height,
		W: screenW,
		H: borderH,
	}
	v.surface.FillRect(&bottomBorder, borderColor)
	leftBorder := sdl.Rect{
		X: 0,
		Y: borderH,
		W: borderW,
		H: height,
	}
	v.surface.FillRect(&leftBorder, borderColor)
	rightBorder := sdl.Rect{
		X: borderW + width,
		Y: borderH,
		W: borderW,
		H: height,
	}
	v.surface.FillRect(&rightBorder, borderColor)
}

func (v *Video) drawBackground() {
	index := v.mem.Load(AddrBackgroundColor) & 0xf
	backgroundColor := uint32(colorMap[index])
	background := sdl.Rect{
		X: borderW,
		Y: borderH,
		W: width,
		H: height,
	}
	v.surface.FillRect(&background, backgroundColor)
}
