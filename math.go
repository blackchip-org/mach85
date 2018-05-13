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

// delete if not used
func signed8(v uint8) int8 {
	if v > 0 && v < 0x7f {
		return int8(v)
	}
	return int8(int16(v) - 0x100)
}
