package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte
type OpCode byte

const (
	OpConstant OpCode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpTrue
	OpFalse
)

type Definition struct {
	Name     string
	OpWidths []int
}

var definitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
	OpAdd:      {"OpAdd", []int{}},
	OpSub:      {"OpSub", []int{}},
	OpMul:      {"OpMul", []int{}},
	OpDiv:      {"OpDiv", []int{}},
	OpPop:      {"OpPop", []int{}},
	OpTrue:     {"OpTrue", []int{}},
	OpFalse:    {"OpFalse", []int{}},
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0

	for i < len(ins) {
		def, err := Lookup(ins[i])

		if err != nil {
			fmt.Fprintf(&out, "ERR %s", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstructions(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstructions(def *Definition, ops []int) string {
	operandCount := len(def.OpWidths)

	if len(ops) != operandCount {
		return fmt.Sprintf("ERR: not enough operands for definition; Wanted %d got %d", operandCount, len(ops))
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, ops[0])
	}

	return fmt.Sprintf("ERR: operandCount unhandled for %s\n", def.Name)
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OpWidths))
	offset := 0

	for i, w := range def.OpWidths {
		switch w {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))

		}

		offset += w

	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func Lookup(op byte) (*Definition, error) {
	if d, ok := definitions[OpCode(op)]; ok {
		return d, nil
	}

	return nil, fmt.Errorf("Opcode %d not found", op)

}

func Make(op OpCode, operands ...int) []byte {
	def, ok := definitions[op]

	if !ok {
		return []byte{}
	}

	insLen := 1

	for _, w := range def.OpWidths {
		insLen += w
	}

	ins := make([]byte, insLen)
	ins[0] = byte(op)

	offset := 1

	for i, o := range operands {
		w := def.OpWidths[i]
		switch w {
		case 2:
			binary.BigEndian.PutUint16(ins[offset:], uint16(o))
		}

		offset += w
	}

	return ins
}
