package mach85

import "github.com/veandco/go-sdl2/sdl"

const (
	width   = 320
	height  = 200
	borderW = 403
	borderH = 284
)

const (
	Black      = 0x000000
	White      = 0xffffff
	Red        = 0x880000
	Cyan       = 0xaaffee
	Purple     = 0xcc44cc
	Green      = 0x00cc55
	Blue       = 0x0000aa
	Yellow     = 0xeeee77
	Orange     = 0xdd8855
	Brown      = 0x664400
	LightRed   = 0xff7777
	DarkGray   = 0x333333
	Gray       = 0x777777
	LightGreen = 0xaaff66
	LightBlue  = 0x0088ff
	LightGray  = 0xbbbbbb
)

var colorMap = [...]int{
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
	window  *sdl.Window
	surface *sdl.Surface
}

func NewVideo() (*Video, error) {
	window, err := sdl.CreateWindow(
		"mach85",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(borderW), int32(borderH),
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
		window:  window,
		surface: surface,
	}, nil
}

func (v *Video) Service() {}
