package mach85

import (
	"fmt"
	"strings"
)

type Operation struct {
	Address     uint16
	Instruction Instruction
	Mode        Mode
	Operand     uint16
	Bytes       []uint8
	Comment     string
}

type Comment struct {
	Address uint16 `json:"address"`
	Text    string `json:"text"`
}

type Disassembler struct {
	PC       uint16
	mem      *Memory
	comments map[uint16]string
}

func NewDisassembler(mem *Memory) *Disassembler {
	return &Disassembler{
		PC:       0xffff,
		mem:      mem,
		comments: map[uint16]string{},
	}
}

func (d *Disassembler) Next() Operation {
	d.PC++
	opcode := d.mem.Load(d.PC)
	address := d.PC
	result := Operation{
		Address:     address,
		Instruction: Illegal,
		Mode:        Implied,
		Operand:     uint16(opcode),
		Bytes:       []uint8{opcode},
		Comment:     d.comments[address],
	}
	op, ok := opcodes[opcode]
	if !ok {
		return result
	}
	len := operandLengths[op.mode]
	operand := uint16(0)
	switch len {
	case 1:
		d.PC++
		b := d.mem.Load(d.PC)
		operand = uint16(b)
		result.Bytes = append(result.Bytes, b)
	case 2:
		d.PC++
		low := d.mem.Load(d.PC)
		d.PC++
		high := d.mem.Load(d.PC)
		operand = uint16(low) + (uint16(high) << 8)
		result.Bytes = append(result.Bytes, low, high)
	}
	result.Instruction = op.inst
	result.Mode = op.mode
	result.Operand = operand
	return result
}

func (o Operation) String() string {
	b0 := fmt.Sprintf("%02x", o.Bytes[0])
	b1, b2 := "  ", "  "
	if len(o.Bytes) > 1 {
		b1 = fmt.Sprintf("%02x", o.Bytes[1])
	}
	if len(o.Bytes) > 2 {
		b2 = fmt.Sprintf("%02x", o.Bytes[2])
	}
	format, ok := operandFormats[o.Mode]
	operand := ""
	if ok {
		// If this is a branch instruction, the value of the operand needs to be
		// added to the current addresss. Add two as it is relative after consuming
		// the instruction
		value := o.Operand
		if o.Mode == Relative {
			value8 := int8(value)
			if value8 >= 0 {
				value = o.Address + uint16(value8) + 2
			} else {
				value = o.Address - uint16(value8*-1) + 2
			}
		}
		// If the format does not contain a formatting directive, just use as is.
		// For example: "asl a"
		if strings.Contains(format, "%") {
			operand = " " + fmt.Sprintf(format, value)
		} else {
			operand = " " + format
		}
	}
	line := fmt.Sprintf("$%04x: %v %v %v  %v%v", o.Address, b0, b1, b2, o.Instruction, operand)
	if o.Comment != "" {
		comments := strings.Split(o.Comment, "\n")
		spaces := 30 - len(line)
		line = line + strings.Repeat(" ", spaces) + comments[0]
		for _, c := range comments[1:] {
			line += "\n" + strings.Repeat(" ", 30) + c
		}
	}
	return line
}

func (d *Disassembler) LoadComments(comments []Comment) {
	for _, c := range comments {
		d.comments[c.Address] = c.Text
	}
}
