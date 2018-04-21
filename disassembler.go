package mach85

import "fmt"

type Operation struct {
	Address     uint16
	Instruction Instruction
	Mode        Mode
	Operand     uint16
	Bytes       []uint8
}

type Disassembler struct {
	PC  uint16
	mem *Memory
}

func NewDisassembler(mem *Memory) *Disassembler {
	return &Disassembler{
		PC:  0xffff,
		mem: mem,
	}
}

func (d *Disassembler) Next() Operation {
	d.PC++
	opcode := d.mem.Load(d.PC)
	result := Operation{
		Address:     d.PC,
		Instruction: Illegal,
		Mode:        Implied,
		Operand:     uint16(opcode),
		Bytes:       []uint8{opcode},
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
		operand = " " + fmt.Sprintf(format, o.Operand)
	}
	return fmt.Sprintf("$%04x: %v %v %v %v%v", o.Address, b0, b1, b2, o.Instruction, operand)
}
