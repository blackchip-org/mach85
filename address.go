package mach85

// https://archive.org/details/Compute_s_Mapping_the_Commodore_64

const (
	AddrProcessorPort     = uint16(0x0001)
	AddrStopKey           = uint16(0x0091)
	AddrKeyboardBufferLen = uint16(0x00c6)
	AddrStack             = uint16(0x0100)
	AddrKeyboardBuffer    = uint16(0x0277)
	AddrBorderColor       = uint16(0xd020)
	AddrBackgroundColor   = uint16(0xd021)
	AddrResetVector       = uint16(0xfffc)
	AddrISR               = uint16(0xff48)
	AddrNmiVector         = uint16(0xfffa)
	AddrIrqVector         = uint16(0xfffe)
)
