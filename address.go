package mach85

const (
	AddrStack           = uint16(0x0100)
	AddrBasicROM        = uint16(0xa000)
	AddrCharacterROM    = uint16(0xd000)
	AddrBorderColor     = uint16(0xd020)
	AddrKernalROM       = uint16(0xe000)
	AddrBackgroundColor = uint16(0xd021)
	AddrResetVector     = uint16(0xfffc)
)
