package encoding

type Decoder interface {
	Decode(uint8) rune
	IsPrintable(uint8) bool
}
