package mach85

import (
	"image/color"
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
	Black      = color.RGBA{0x00, 0x00, 0x00, 0xff}
	White      = color.RGBA{0xff, 0xff, 0xff, 0xff}
	Red        = color.RGBA{0x88, 0x00, 0x00, 0xff}
	Cyan       = color.RGBA{0xaa, 0xff, 0xee, 0xff}
	Purple     = color.RGBA{0xcc, 0x44, 0xcc, 0xff}
	Green      = color.RGBA{0x00, 0xcc, 0x55, 0xff}
	Blue       = color.RGBA{0x00, 0x00, 0xaa, 0xff}
	Yellow     = color.RGBA{0xee, 0xee, 0x77, 0xff}
	Orange     = color.RGBA{0xdd, 0x88, 0x55, 0xff}
	Brown      = color.RGBA{0x66, 0x44, 0x00, 0xff}
	LightRed   = color.RGBA{0xff, 0x77, 0x77, 0xff}
	DarkGray   = color.RGBA{0x33, 0x33, 0x33, 0xff}
	Gray       = color.RGBA{0x77, 0x77, 0x77, 0xff}
	LightGreen = color.RGBA{0xaa, 0xff, 0x66, 0xff}
	LightBlue  = color.RGBA{0x00, 0x88, 0xff, 0xff}
	LightGray  = color.RGBA{0xbb, 0xbb, 0xbb, 0xff}
)

var colorMap = [...]color.RGBA{
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
	c := colorMap[index]
	v.renderer.SetDrawColor(c.R, c.G, c.B, c.A)
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
	c := colorMap[index]
	v.renderer.SetDrawColor(c.R, c.G, c.B, c.A)
	background := sdl.Rect{
		X: borderW,
		Y: borderH,
		W: width,
		H: height,
	}
	v.renderer.FillRect(&background)
}

func CharGen(renderer *sdl.Renderer, chargen *ROM) (*sdl.Texture, error) {
	w := 32 * 8
	h := 16 * 8
	t, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET, int32(w), int32(h))
	if err != nil {
		return nil, err
	}
	renderer.SetRenderTarget(t)
	renderer.SetDrawColorArray(0xff, 0xff, 0xff, 0xff)
	baseX := 0
	baseY := 0
	addr := uint16(0)
	for baseY < h {
		for y := baseY; y < baseY+8; y++ {
			line := chargen.Load(addr)
			addr++
			for x := baseX; x < baseX+8; x++ {
				bit := line & 0x80
				line = line << 1
				if bit != 0 {
					renderer.DrawPoint(int32(x), int32(y))
				}
			}
		}
		baseX += 8
		if baseX >= w {
			baseX = 0
			baseY += 8
		}
	}
	t.SetBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetRenderTarget(nil)
	return t, nil
}
