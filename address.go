package mach85

// https://archive.org/details/Compute_s_Mapping_the_Commodore_64

const (
	// AddrD6510 is the 6510 on-chip IO DATA direction register
	AddrD6510 = uint16(0x0000)

	// AddrR6510 is the memory configuration and cassette register
	AddrR6510 = uint16(0x0001)

	AddrStack           = uint16(0x0100)
	AddrBorderColor     = uint16(0xd020)
	AddrBackgroundColor = uint16(0xd021)
	AddrResetVector     = uint16(0xfffc)
)
