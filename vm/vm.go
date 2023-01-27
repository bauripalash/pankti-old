package vm

import (
	"fmt"

	"go.cs.palashbauri.in/pankti/code"
	"go.cs.palashbauri.in/pankti/compiler"
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/token"
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
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			lval := left.(*object.Number).Value
			rval := right.(*object.Number).Value
			r, _, _ := number.NumberOperation(token.PLUS, lval, rval)
			vm.push(&object.Number{Value: r})
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

func (v *VM) pop() object.Obj {
	o := v.stack[v.sp-1]
	v.sp--
	return o
}
