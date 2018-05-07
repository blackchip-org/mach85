package mach85

// https://archive.org/details/Compute_s_Mapping_the_Commodore_64

const (
	AddrProcessorPort   = uint16(0x0001)
	AddrStack           = uint16(0x0100)
	AddrBorderColor     = uint16(0xd020)
	AddrBackgroundColor = uint16(0xd021)
	AddrResetVector     = uint16(0xfffc)
)
