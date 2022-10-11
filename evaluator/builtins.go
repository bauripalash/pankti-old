package evaluator

import (
	"fmt"
	"pankti/object"
	"pankti/stdlib"
)

func lenFunc(args []object.Obj) object.Obj {
	if len(args) != 1 {
		return NewBareErr("wrong number of arguments. got %d but wanted 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return object.MakeIntNumber(int64(len(arg.Value)))
	case *object.Array:
		return object.MakeIntNumber(int64(len(arg.Elms)))
	default:
		return NewBareErr("argument type %s to `len` is not supported", args[0].Type())
	}
}

func firstFunc(args []object.Obj) object.Obj {

	if len(args) != 1 {
		return NewBareErr("wrong number of argument %d", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return NewBareErr("first cannot be used with %s", args[0].Type())
	}

	array := args[0].(*object.Array)
	if len(array.Elms) > 0 {
		return array.Elms[0]
	}
	return NULL
}

func lastFunc(args []object.Obj) object.Obj {

	if len(args) != 1 {
		return NewBareErr("wrong number of argument %d", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return NewBareErr("last cannot be used with %s", args[0].Type())
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
		return NewBareErr("wrong number of argument %d", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return NewBareErr("rest cannot be used with %s", args[0].Type())
	}

	array := args[0].(*object.Array)
	arrLen := len(array.Elms)
	if arrLen > 0 {
		newElms := make([]object.Obj, arrLen-1, arrLen-1)
		copy(newElms, array.Elms[1:arrLen])
		return &object.Array{Elms: newElms}
	}
	return NULL
}

func pushFunc(args []object.Obj) object.Obj {

	if len(args) != 2 {
		return NewBareErr("wrong number of argument %d", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return NewBareErr("push cannot be used with %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	arrLen := len(arr.Elms)

	newElms := make([]object.Obj, arrLen+1, arrLen+1)
	copy(newElms, arr.Elms)
	newElms[arrLen] = args[1]
	return &object.Array{Elms: newElms}
}

func showFunc(args []object.Obj) object.Obj {

	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

var builtins = map[string]*object.Builtin{

	"len": {
		Fn: func(args ...object.Obj) object.Obj {
			return lenFunc(args)
		},
	},

	"আয়তন": {
		Fn: func(args ...object.Obj) object.Obj {
			return lenFunc(args)
		},
	},

	"ayoton": {
		Fn: func(args ...object.Obj) object.Obj {
			return lenFunc(args)
		},
	},

	"প্রথম": {
		Fn: func(args ...object.Obj) object.Obj {
			return firstFunc(args)
		},
	},

	"first": {
		Fn: func(args ...object.Obj) object.Obj {
			return firstFunc(args)
		},
	},

	"prothom": {
		Fn: func(args ...object.Obj) object.Obj {
			return firstFunc(args)
		},
	},

	"শেষ": {
		Fn: func(args ...object.Obj) object.Obj {
			return lastFunc(args)
		},
	},

	"last": {
		Fn: func(args ...object.Obj) object.Obj {
			return lastFunc(args)
		},
	},

	"sesh": {
		Fn: func(args ...object.Obj) object.Obj {
			return lastFunc(args)
		},
	},

	"বাদবাকি": {
		Fn: func(args ...object.Obj) object.Obj {
			return restFunc(args)
		},
	},

	"rest": {
		Fn: func(args ...object.Obj) object.Obj {
			return restFunc(args)
		},
	},

	"badbaki": {
		Fn: func(args ...object.Obj) object.Obj {
			return restFunc(args)
		},
	},

	"যোগ": {
		Fn: func(args ...object.Obj) object.Obj {
			return pushFunc(args)
		},
	},

	"push": {
		Fn: func(args ...object.Obj) object.Obj {
			return pushFunc(args)
		},
	},

	"jog": {
		Fn: func(args ...object.Obj) object.Obj {
			return pushFunc(args)
		},
	},

	"দেখাও": {
		Fn: func(args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},

	"show": {
		Fn: func(args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},

	"dekhau": {
		Fn: func(args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},
	"ইপচ": {
		Fn: func(args ...object.Obj) object.Obj {
			return stdlib.UnixTimeFunc(args)
		},
	},
}
