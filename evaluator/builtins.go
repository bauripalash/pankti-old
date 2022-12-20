package evaluator

import (
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/stdlib"
)

func lenFunc(args []object.Obj) object.Obj {
	if len(args) != 1 {
		return object.NewBareErr(
			"wrong number of arguments. got %d but wanted 1",
			len(args),
		)
	}

	switch arg := args[0].(type) {
	case *object.String:
		return object.MakeIntNumber(int64(len(arg.Value)))
	case *object.Array:
		return object.MakeIntNumber(int64(len(arg.Elms)))
	default:
		return object.NewBareErr("argument type %s to `len` is not supported", args[0].Type())
	}
}

func firstFunc(args []object.Obj) object.Obj {

	if len(args) != 1 {
		return object.NewBareErr("wrong number of argument %d", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewBareErr("first cannot be used with %s", args[0].Type())
	}

	array := args[0].(*object.Array)
	if len(array.Elms) > 0 {
		return array.Elms[0]
	}
	return NULL
}

func lastFunc(args []object.Obj) object.Obj {

	if len(args) != 1 {
		return object.NewBareErr("wrong number of argument %d", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewBareErr("last cannot be used with %s", args[0].Type())
	}

	array := args[0].(*object.Array)
	arr_len := len(array.Elms)
	if arr_len > 0 {
		return array.Elms[arr_len-1]
	}
	return NULL
}

func restFunc(args []object.Obj) object.Obj {

	if len(args) != 1 {
		return object.NewBareErr("wrong number of argument %d", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewBareErr("rest cannot be used with %s", args[0].Type())
	}

	array := args[0].(*object.Array)
	arrLen := len(array.Elms)
	if arrLen > 0 {
		newElms := make([]object.Obj, arrLen-1)
		copy(newElms, array.Elms[1:arrLen])
		return &object.Array{Elms: newElms}
	}
	return NULL
}

func pushFunc(args []object.Obj) object.Obj {

	if len(args) != 2 {
		return object.NewBareErr("wrong number of argument %d", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewBareErr("push cannot be used with %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	arrLen := len(arr.Elms)

	newElms := make([]object.Obj, arrLen+1)
	copy(newElms, arr.Elms)
	newElms[arrLen] = args[1]
	return &object.Array{Elms: newElms}
}

func showFunc(args []object.Obj) object.Obj {
	output := []string{}
	for _, arg := range args {
		//fmt.Println(arg.Inspect())
		output = append(output, arg.Inspect())
	}
	return &object.ShowObj{Value: output, Token: NULL.GetToken()}
}

var builtins = map[string]*object.Builtin{

	"__len": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return lenFunc(args)
		},
	},

	"__first": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return firstFunc(args)
		},
	},

	"__last": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return lastFunc(args)
		},
	},

	"__res": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return restFunc(args)
		},
	},

	"__push": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return pushFunc(args)
		},
	},

	"দেখাও": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},

	"show": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},

	"dekhau": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},
	"__epoch": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.UnixTimeFunc(args)
		},
	},

	"__isonow": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.UtcDateISO(args)
		},
	},

	"__osname": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.GetOS(args)
		},
	},

	"__osarch": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.GetArch(args)
		},
	},

	"__readfile": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.ReadFile(args)
		},
	},

	"__exists": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.FileDirExists(args)
		},
	},

	"__create_empty": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.CreateEmptyFile(args)
		},
	},

	"__write_file": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.WriteToFile(args)
		},
	},

	// Maths

	"__sqrt": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.DoSqrt(eh, args)
		},
	},

	"__log_ten": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.Log10(args)
		},
	},
	"__list_sum": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.DoListSum(args)
		},
	},

	"__gcd": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.GetGCD(args)
		},
	},

	"__lcm": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.GetLCM(args)
		},
	},

	"__pow": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.DoPow(args)
		},
	},

	"__log_e": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.LogE(args)
		},
	},

	"__log_x": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.LogX(args)
		},
	},

	"__cosine": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.Cosine(args)
		},
	},

	"__sine": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.Sine(args)
		},
	},

	"__acos": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.Acos(args)
		},
	},

	"__asin": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.Asin(args)
		},
	},

	"__tan": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.Tangent(args)
		},
	},

	"__atan": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.Atan(args)
		},
	},

	"__atan_two": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.Atan2(args)
		},
	},

	"__to_deg": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.ToDegree(args)
		},
	},

	"__to_rad": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.ToRadians(args)
		},
	},

	"__get_pi": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.GetPI(args)
		},
	},

	"__get_e": {
		Fn: func(eh *object.ErrorHelper, args ...object.Obj) object.Obj {
			return stdlib.GetE(args)
		},
	},
}
