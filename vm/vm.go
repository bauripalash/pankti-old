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

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.exeBinaryOp(op)
			if err != nil {
				return nil
			}
		case code.OpTrue:
			if err := vm.push(True); err != nil {
				return err
			}
		case code.OpFalse:
			if err := vm.push(False); err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGT:
			if err := vm.exeComparison(op); err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()
		}

	}

	return nil
}

func (vm *VM) exeComparison(op code.OpCode) error {
	r := vm.pop()
	l := vm.pop()

	if r.Type() == object.NUM_OBJ && r.Type() == object.NUM_OBJ {
		return vm.exeNumComparison(op, l, r)
	}

	switch op {
	case code.OpEqual:
		return vm.push(getBoolObj(l == r))
	case code.OpNotEqual:
		return vm.push(getBoolObj(l != r))
	default:
		return fmt.Errorf("unknown operator %d (%s %s)", op, l.Type(), r.Type())
	}

}

func (vm *VM) exeNumComparison(op code.OpCode, l, r object.Obj) error {
	lval := l.(*object.Number)
	rval := r.(*object.Number)
	var v bool
	switch op {
	case code.OpEqual:
		_, v, _ = number.NumberOperation(token.EQEQ, lval.Value, rval.Value)
		return vm.push(getBoolObj(v))
	case code.OpNotEqual:
		_, v, _ = number.NumberOperation(token.NOT_EQ, lval.Value, rval.Value)
		return vm.push(getBoolObj(v))
	case code.OpGT:
		_, v, _ = number.NumberOperation(token.GT, lval.Value, rval.Value)
		return vm.push(getBoolObj(v))
	default:
		return fmt.Errorf("Unknown Op : %d", op)
	}
}

func getBoolObj(b bool) *object.Boolean {
	if b {
		return True
	} else {
		return False
	}
}

func (vm *VM) exeBinaryOp(op code.OpCode) error {
	right := vm.pop()
	left := vm.pop()
	lType := left.Type()
	rType := right.Type()

	if lType == object.NUM_OBJ && rType == object.NUM_OBJ {
		return vm.exeNumBinaryOp(op, left, right)
	}

	return fmt.Errorf("Unsupported type for binary operation : %s %s", lType, rType)

}

func (vm *VM) exeNumBinaryOp(op code.OpCode, left, right object.Obj) error {

	lval := left.(*object.Number).Value
	rval := right.(*object.Number).Value
	var result number.Number

	switch op {
	case code.OpAdd:
		result, _, _ = number.NumberOperation(token.PLUS, lval, rval)
	case code.OpSub:
		result, _, _ = number.NumberOperation(token.MINUS, lval, rval)
	case code.OpMul:
		result, _, _ = number.NumberOperation(token.MUL, lval, rval)
	case code.OpDiv:
		result, _, _ = number.NumberOperation(token.DIV, lval, rval)
		fmt.Println(result)
	default:
		return fmt.Errorf("Unknown number operator : %d", op)

	}

	return vm.push(&object.Number{Value: result})
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

func (vm *VM) LastPoppedStackItem() object.Obj {
	return vm.stack[vm.sp]
}
