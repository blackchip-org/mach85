package petscii

const replacementChar = 0xfffd

var UnshiftedDecoder = func(code uint8) (rune, bool) {
	ch := unshifted[code]
	return ch, isPrintable(code) && ch != replacementChar
}

var ShiftedDecoder = func(code uint8) (rune, bool) {
	ch := shifted[code]
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
