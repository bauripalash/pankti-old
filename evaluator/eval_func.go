package evaluator

import (
	"bytes"

	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/token"
)

func applyFunc(
	fn object.Obj,
	caller token.Token,
	args []object.Obj,
	isMod bool,
	env *object.EnvMap,
	eh *object.ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {

	//fmt.Println(caller)

	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Params) == len(args) {
			eEnv := extendFuncEnv(fn, args)

			eX := object.NewEnvMap()

			eX.Envs[object.DEFKEY] = *object.NewEnclosedEnv(eEnv)

			evd := Eval(fn.Body, eX, *eh, printBuff, isGui)
			return unwrapReturnValue(evd)
		} else {

			return object.NewErr(caller, eh, false, "Function call doesn't have required arguments provided; wanted = %d but got %d", len(fn.Params), len(args))
		}
	case *object.Builtin:
		return fn.Fn(eh, args...)
	default:
		return object.NewBareErr("%s is not a function", fn.Type())

	}
}

func extendFuncEnv(fn *object.Function, args []object.Obj) *object.Env {
	env := object.NewEnclosedEnv(fn.Env)

	//if len(args) > 0 {
	for pId, param := range fn.Params {
		env.Set(param.Value, args[pId])
	}
	//}

	return env
}

func unwrapReturnValue(o object.Obj) object.Obj {
	if rv, ok := o.(*object.ReturnValue); ok {
		return rv.Value
	}

	return o

}
