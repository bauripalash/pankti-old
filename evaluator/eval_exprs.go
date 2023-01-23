package evaluator

import (
	"bytes"

	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/token"
)

func evalInfixExpr(
	optok token.Token,
	l, r object.Obj,
	eh *object.ErrorHelper,
) object.Obj {

	op := optok.Literal

	//fmt.Println(l.GetToken(), r.Type())
	switch {
	case l.Type() == object.NUM_OBJ && r.Type() == object.NUM_OBJ:
		return evalNumInfixExpr(op, l, r, eh)
		//}
		//fmt.Println("FI-> ", l , r)
		//return NewErr("has Float")
	case l.Type() == object.STRING_OBJ && r.Type() == object.STRING_OBJ:
		return evalStringInfixExpr(op, l, r, eh)
	case op == "==":
		return getBoolObj(l == r)
	case op == "!=":
		return getBoolObj(l != r)
	case optok.Type == token.AND:
		lval := getBoolTypeFromObj(l)
		rval := getBoolTypeFromObj(r)
		return getBoolObj(lval && rval)

	case optok.Type == token.OR:
		lval := getBoolTypeFromObj(l)
		rval := getBoolTypeFromObj(r)
		return getBoolObj(lval || rval)
	case l.Type() != r.Type():
		return object.NewErr(
			l.GetToken(),
			eh,
			false,
			"Type mismatch:  %s %s %s ",
			l.Type(),
			op,
			r.Type(),
		)

	default:
		return object.NewErr(
			l.GetToken(),
			eh,
			false,
			"unknown Operator : %s %s %s",
			l.Type(),
			op,
			r.Type(),
		)
	}
}

func evalPrefixExpr(
	op string,
	right object.Obj,
	eh *object.ErrorHelper,
) object.Obj {
	switch op {
	case "!":
		return evalBangOp(right)
	case "-":
		return evalMinusPrefOp(right, eh)
	default:
		return object.NewBareErr("Unknown Operator : %s%s", op, right.Type())

	}
}

func evalExprs(
	es []ast.Expr,
	env *object.EnvMap,
	eh *object.ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) []object.Obj {
	var res []object.Obj

	for _, e := range es {
		ev := Eval(e, env, *eh, printBuff, isGui)

		if object.IsErr(ev) {
			return []object.Obj{ev}
		}

		res = append(res, ev)
	}

	return res
}
