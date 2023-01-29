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
var Null = &object.Null{}

const GlobalsSize = 65536

type VM struct {
	constants    []object.Obj
	instructions code.Instructions

	stack []object.Obj
	sp    int

	globals []object.Obj
}

func NewVM(bc compiler.ByteCode) *VM {
	return &VM{
		instructions: bc.Instructions,
		constants:    bc.Constants,

		stack:   make([]object.Obj, StackSize),
		sp:      0,
		globals: make([]object.Obj, GlobalsSize),
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
		case code.OpBang:
			if err := vm.exeBangOp(); err != nil {
				return err
			}
		case code.OpMinus:
			if err := vm.exeMinusOp(); err != nil {
				return err
			}
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case code.OpJumpNotTruthy:
			//pos :=

			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			cond := vm.pop()

			if !isTruthy(cond) {
				ip = pos - 1
			}
		case code.OpNull:
			if err := vm.push(Null); err != nil {
				return err
			}
		case code.OpSetGlobal:
			gIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[gIndex] = vm.pop()
		case code.OpGetGlobal:
			gIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			if err := vm.push(vm.globals[gIndex]); err != nil {
				return err
			}
		case code.OpArray:
			numElms := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			arr := vm.buildArray(vm.sp-numElms, vm.sp)
			vm.sp = vm.sp - numElms
			if err := vm.push(arr); err != nil {
				return err
			}
		case code.OpHash:
			numElms := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			hash, err := vm.buildHash(vm.sp-numElms, vm.sp)

			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElms
			err = vm.push(hash)
			if err != nil {
				return err
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			if err := vm.exeIndexExpr(left, index); err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()
		}

	}

	return nil
}

func (vm *VM) exeIndexExpr(left, index object.Obj) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUM_OBJ:
		return vm.exeArrIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.exeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported %s", left.Type())
	}
}

func (vm *VM) exeArrIndex(arr, index object.Obj) error {
	arrObj := arr.(*object.Array)
	fi, _ := index.(*object.Number).Value.GetAsFloat()
	i := int64(fi)
	max := int64(len(arrObj.Elms) - 1)

	if i < 0 || i > max {
		return vm.push(Null)
	}

	return vm.push(arrObj.Elms[i])
}

func (vm *VM) exeHashIndex(hash, index object.Obj) error {
	hashObj := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("Unsupported Hash Key with Type %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Obj, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		k := vm.stack[i]
		val := vm.stack[i+1]

		pair := object.HashPair{Key: k, Value: val}

		hashKey, ok := k.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("Unusable as hash key : %s", k.Type())
		}

		hashedPairs[hashKey.HashKey()] = pair

		//fmt.Println(hashedPairs)
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Obj {
	elms := make([]object.Obj, endIndex-startIndex)
	for i := startIndex; i < endIndex; i++ {
		elms[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elms: elms}
}

func (vm *VM) exeBangOp() error {
	op := vm.pop()
	switch op {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) exeMinusOp() error {
	op := vm.pop()
	if op.Type() != object.NUM_OBJ {
		return fmt.Errorf("Unsupported type for neg %s", op.Type())
	}

	val := op.(*object.Number).Value
	return vm.push(&object.Number{Value: number.MakeNeg(val)})
}

func isTruthy(obj object.Obj) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
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

	switch {
	case lType == object.NUM_OBJ && rType == object.NUM_OBJ:
		return vm.exeNumBinaryOp(op, left, right)
	case lType == object.STRING_OBJ && rType == object.STRING_OBJ:
		return vm.exeStrBinaryOp(op, left, right)

	default:
		return fmt.Errorf("Unsupported type for binary operation : %s %s", lType, rType)
	}

}

func (vm *VM) exeStrBinaryOp(op code.OpCode, l, r object.Obj) error {
	if op != code.OpAdd {

		return fmt.Errorf("unknown string operator %d", op)
	}

	lval := l.(*object.String).Value
	rval := r.(*object.String).Value

	return vm.push(&object.String{Value: lval + rval})
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
