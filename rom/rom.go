package rom

import (
	"flag"
	"path/filepath"
)

var Path string

func init() {
	defaultPath := filepath.Join("..", "..", "rom")
	flag.StringVar(&Path, "rom-path", defaultPath, "path to ROM files")
}
