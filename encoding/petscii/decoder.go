package petscii

type UnshiftedDecoder struct{}

func (d UnshiftedDecoder) Decode(index uint8) rune {
	return unshifted[index]
}

func (d UnshiftedDecoder) IsPrintable(v uint8) bool {
	return isPrintable(v)
}

type ShiftedDecoder struct{}

func (d ShiftedDecoder) Decode(index uint8) rune {
	return shifted[index]
}

func (d ShiftedDecoder) IsPrintable(v uint8) bool {
	return isPrintable(v)
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
