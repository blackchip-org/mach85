package rom

import (
	"flag"
)

var Path string

func init() {
	defaultPath := "rom"
	flag.StringVar(&Path, "rom-path", defaultPath, "path to ROM files")
}
