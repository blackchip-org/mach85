package mach85

// http://dustlayer.com/index-vic-ii/
// http://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt

import (
	"flag"
	"image/color"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var scale int

func init() {
	flag.IntVar(&scale, "video-scale", 2, "scale video screen size")
}

const (
	width      = 320
	height     = 200
	screenW    = 404 // actually 403?
	screenH    = 284
	borderW    = (screenW - width) / 2
	borderH    = (screenH - height) / 2
	charSheetW = 32
	charSheetH = 16
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
	charSheet  *sdl.Texture
}

func NewVideo(mem *Memory) (*Video, error) {
	window, err := sdl.CreateWindow(
		"mach85",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(screenW*scale), int32(screenH*scale),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, err
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	renderer.SetScale(float32(scale), float32(scale))
	v := &Video{
		mem:      mem,
		window:   window,
		renderer: renderer,
	}
	v.Service()
	return v, nil
}

func (v *Video) Service() error {
	now := time.Now()
	if now.Sub(v.lastUpdate) < time.Millisecond*16 {
		return nil
	}
	v.lastUpdate = now
	v.mem.Store(0xd012, 00) // HACK: set raster line to zero
	if v.charSheet == nil {
		if err := v.genCharSheet(); err != nil {
			return err
		}
	}
	v.drawBorder()
	v.drawBackground()
	v.drawCharacters()
	v.renderer.Present()
	return nil
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

func (v *Video) drawCharacters() {
	mem64 := v.mem.Base.(*Memory64)
	prev := mem64.SetMode(0)
	defer mem64.SetMode(prev)

	io := mem64.Chunks[IO]
	addrScreenMem := uint16(0x0400)
	addrColorMem := uint16(0x0800)
	baseX := 0
	baseY := 0
	for baseY < height {
		ch := mem64.Load(addrScreenMem)
		color := colorMap[io.Load(addrColorMem)&0x0f]
		v.charSheet.SetColorMod(color.R, color.G, color.B)
		chx := int32(ch) % charSheetW * 8
		chy := int32(ch) / charSheetW * 8
		src := sdl.Rect{X: chx, Y: chy, W: 8, H: 8}
		dest := sdl.Rect{
			X: int32(baseX + borderW),
			Y: int32(baseY + borderH),
			W: 8,
			H: 8,
		}
		v.renderer.Copy(v.charSheet, &src, &dest)
		addrScreenMem++
		addrColorMem++
		baseX += 8
		if baseX >= width {
			baseX = 0
			baseY += 8
		}
	}
}

func (v *Video) genCharSheet() error {
	mem64 := v.mem.Base.(*Memory64)
	chargen := mem64.Chunks[CharROM]
	charSheet, err := CharGen(v.renderer, chargen)
	if err != nil {
		return err
	}
	v.charSheet = charSheet
	return nil
}

func CharGen(renderer *sdl.Renderer, chargen MemoryChunk) (*sdl.Texture, error) {
	w := charSheetW * 8
	h := charSheetH * 8
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
