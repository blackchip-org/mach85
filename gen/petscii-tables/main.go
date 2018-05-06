// +build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(0)
	in, err := os.Open("petscii.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	var table [2][0x100]rune

	s := bufio.NewScanner(in)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		if line[0] == '#' {
			continue
		}
		fields := strings.Split(line, "|")
		cp := strings.TrimSpace(fields[1])
		charCode, err := strconv.ParseUint(cp[1:], 16, 8)
		if err != nil {
			panic(err)
		}
		ucp := strings.TrimSpace(fields[2])
		ucps := regexp.MustCompile("\\s+").Split(ucp, -1)
		if len(ucps) == 1 {
			ucps = append(ucps, ucps[0])
		}

		for tableN, uni := range ucps {
			if uni == "-" {
				table[tableN][charCode] = 0xfffd
			} else {
				value, err := strconv.ParseUint(uni[2:], 16, 16)
				if err != nil {
					panic(err)
				}
				table[tableN][int(charCode)] = rune(value)
			}
		}
	}
	if s.Err() != nil {
		panic(s.Err())
	}

	var out bytes.Buffer
	out.WriteString("package mach85\n")
	out.WriteString("var petsciiUnshifted = [...]rune {\n")
	for _, val := range table[0] {
		out.WriteString(fmt.Sprintf("0x%04x,\n", val))
	}
	out.WriteString("}\n")
	out.WriteString("var petsciiShifted = [...]rune {\n")
	for _, val := range table[1] {
		out.WriteString(fmt.Sprintf("0x%04x,\n", val))
	}
	out.WriteString("}\n")

	outfile := filepath.Join("..", "..", "petscii.go")
	err = ioutil.WriteFile(outfile, out.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
