package vm

import (
	"fmt"

	"go.cs.palashbauri.in/pankti/code"
	"go.cs.palashbauri.in/pankti/compiler"
	"go.cs.palashbauri.in/pankti/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Obj
	instructions code.Instructions

	stack []object.Obj
	sp    int
}

func NewVM(bc compiler.ByteCode) *VM {
	return &VM{
		instructions: bc.Instructions,
		constants:    bc.Constants,

		stack: make([]object.Obj, StackSize),
		sp:    0,
	}
}

func (vm *VM) StackTop() object.Obj {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.OpCode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return nil
			}
		}

	}

	return nil
}

func (v *VM) push(o object.Obj) error {
	if v.sp >= StackSize {
		return fmt.Errorf("Stack Overflow")
	}

	v.stack[v.sp] = o
	v.sp++

	return nil
}
