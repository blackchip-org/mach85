// +build ignore

package main

// https://github.com/mist64/c64rom/blob/master/c64rom_en.txt

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/blackchip-org/mach85"
)

type comment struct {
	address uint16
	text    string
}

func main() {
	log.SetFlags(0)
	in, err := os.Open("c64rom_en.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	s := bufio.NewScanner(in)
	s.Split(bufio.ScanLines)

	comments := []comment{}
	commentCanContinue := false
	for s.Scan() {
		line := s.Text()
		if len(line) < 33 {
			commentCanContinue = false
			continue
		}
		if line[0] == ' ' && commentCanContinue {
			comment := strings.TrimSpace(line[32:])

			comments[len(comments)-1].text += "\n" + comment
			continue
		}
		if !strings.HasPrefix(line, ".,") {
			commentCanContinue = false
			continue
		}
		address, err := strconv.ParseUint(line[2:6], 16, 16)
		if err != nil {
			log.Printf("unable to parse address: %v", address)
			continue
		}
		comment := comment{
			address: uint16(address),
			text:    strings.TrimSpace(line[32:]),
		}
		comments = append(comments, comment)
		commentCanContinue = true
	}

	source := mach85.NewSource()
	for _, c := range comments {
		source.Comments[c.address] = c.text
	}
	outfile := filepath.Join("..", "..", "rom", "c64rom_en.source")
	out, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	encoder := json.NewEncoder(out)
	encoder.Encode(source)
}
