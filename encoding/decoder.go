package encoding

type Decoder func(uint8) (rune, bool)
