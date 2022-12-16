package evaluator

import (
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
)

func evalNumInfixExpr(
	op string,
	l, r object.Obj,
	eh *ErrorHelper,
) object.Obj {

	lval := l.(*object.Number).Value
	rval := r.(*object.Number).Value

	//fmt.Println(lval.GetType() , rval.GetType())

	val, cval, noerr := number.NumberOperation(op, lval, rval)
	if val.Value != nil && noerr {
		return &object.Number{Value: val, IsInt: val.IsInt}
	} else if val.Value == nil && noerr {
		return getBoolObj(cval)
	} else {
		return NewBareErr("Unknown Operator for Numbers %s", op)
	}

}

func evalStringInfixExpr(
	op string,
	l, r object.Obj,
	eh *ErrorHelper,
) object.Obj {
	lval := l.(*object.String).Value
	rval := r.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: lval + rval}
	case "==":
		return getBoolObj(lval == rval)
	case "!=":
		return getBoolObj(lval != rval)
	default:
		return NewErr(
			l.GetToken(),
			eh,
			false,
			"Unknown Operator %s %s %s",
			l.Type(),
			op,
			r.Type(),
		)

	}

}
