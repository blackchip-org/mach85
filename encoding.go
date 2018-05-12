package mach85

type Decoder func(uint8) (rune, bool)

const replacementChar = 0xfffd

var PetsciiUnshiftedDecoder = func(code uint8) (rune, bool) {
	ch := petsciiUnshifted[code]
	return ch, isPrintable(code) && ch != replacementChar
}

var PetsciiShiftedDecoder = func(code uint8) (rune, bool) {
	ch := petsciiShifted[code]
	return ch, isPrintable(code) && ch != replacementChar
}

func isPrintable(v uint8) bool {
	if v < 0x20 {
		return false
	}
	if v >= 0x80 && v <= 0xa0 {
		return false
	}
	return true
}

// http://sta.c64.org/cbm64pettoscr.html

var ScreenUnshiftedDecoder = func(code uint8) (rune, bool) {
	return decoder(code, PetsciiUnshiftedDecoder)
}

var ScreenShiftedDecoder = func(code uint8) (rune, bool) {
	return decoder(code, PetsciiShiftedDecoder)
}

func decoder(code uint8, decode Decoder) (rune, bool) {
	switch {
	case code == 0x5e:
		return decode(0xff)
	case code >= 0x00 && code <= 0x1f:
		return decode(code + 64)
	case code >= 0x20 && code <= 0x3f:
		return decode(code)
	case code >= 0x40 && code <= 0x5f:
		return decode(code + 32)
	case code >= 0x60 && code <= 0x7f:
		return decode(code + 64)
	case code >= 0x80 && code <= 0x9f:
		return decode(code - 128)
	case code >= 0xc0 && code <= 0xdf:
		return decode(code - 64)
	}
	return decode(code)
}

func PetsciiEncoder(ch rune) (uint8, bool) {
	if ch >= 0x61 && ch <= 0x7a {
		return uint8(ch - 0x20), true
	}
	return uint8(ch), true
}
