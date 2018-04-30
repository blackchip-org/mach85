package mach85

func fromBCD(v uint8) uint8 {
	low := v & 0x0f
	high := v >> 4
	return high*10 + low
}

func toBCD(v uint8) uint8 {
	low := v % 10
	high := (v / 10) % 10
	return high<<4 | low
}
