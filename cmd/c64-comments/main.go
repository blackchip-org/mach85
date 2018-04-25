package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/blackchip-org/mach85"
)

func main() {
	log.SetFlags(0)
	in, err := os.Open("c64rom_en.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	s := bufio.NewScanner(in)
	s.Split(bufio.ScanLines)

	comments := []mach85.Comment{}
	commentCanContinue := false
	for s.Scan() {
		line := s.Text()
		if len(line) < 33 {
			commentCanContinue = false
			continue
		}
		if line[0] == ' ' && commentCanContinue {
			comment := strings.TrimSpace(line[32:])
			comments[len(comments)-1].Text += "\n" + comment
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
		comment := mach85.Comment{
			Address: uint16(address),
			Text:    strings.TrimSpace(line[32:]),
		}
		comments = append(comments, comment)
		commentCanContinue = true
	}

	out, err := os.Create("../mach85/c64rom.debug")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	encoder := json.NewEncoder(out)
	encoder.Encode(comments)
}
