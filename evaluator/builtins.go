package evaluator

import (
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/stdlib"
	"go.cs.palashbauri.in/pankti/token"
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
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return lenFunc(args)
		},
	},

	"sethv": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.SetHashTableElm(eh, env, caller, args)
		},
	},
	"getkeys": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetAllKVsOfHashTable(true, eh, env, caller, args)
		},
	},
	"getvals": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetAllKVsOfHashTable(false, eh, env, caller, args)
		},
	},

	"__first": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return firstFunc(args)
		},
	},

	"__last": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return lastFunc(args)
		},
	},

	"__res": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return restFunc(args)
		},
	},

	"__push": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return pushFunc(args)
		},
	},

	"দেখাও": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},

	"show": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},

	"dekhau": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return showFunc(args)
		},
	},
	"__epoch": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.UnixTimeFunc(args)
		},
	},

	"__isonow": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.UtcDateISO(args)
		},
	},

	// Maths

	"__sqrt": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.DoSqrt(eh, args)
		},
	},

	"__log_ten": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.Log10(eh, args)
		},
	},
	"__list_sum": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.DoListSum(eh, args)
		},
	},

	"__gcd": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetGCD(eh, caller, args)
		},
	},

	"__lcm": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetLCM(eh, args)
		},
	},

	"__pow": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.DoPow(eh, args)
		},
	},

	"__log_e": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.LogE(eh, args)
		},
	},

	"__log_x": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.LogX(eh, args)
		},
	},

	"__cosine": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.Cosine(eh, args)
		},
	},

	"__sine": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.Sine(eh, args)
		},
	},

	"__acos": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.Acos(eh, args)
		},
	},

	"__asin": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.Asin(eh, args)
		},
	},

	"__tan": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.Tangent(eh, args)
		},
	},

	"__atan": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.Atan(eh, args)
		},
	},

	"__atan_two": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.Atan2(eh, args)
		},
	},

	"__to_deg": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ToDegree(eh, args)
		},
	},

	"__to_rad": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ToRadians(eh, args)
		},
	},

	"__get_pi": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetPI(args)
		},
	},

	"__get_e": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetE(args)
		},
	},
	"__to_number": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ToNumber(eh, args)
		},
	},

	"__to_number_float": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ConvertToFloat(eh, args)
		},
	},

	"__to_number_int": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ConvertToInt(eh, args)
		},
	},

	"__get_random_with_arg": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GenerateRandom(eh, args)
		},
	},

	"__string_split": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.SplitString(eh, args)
		},
	},

	"__string_join": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.JoinAsString(eh, args)
		},
	},

	"__string_convert": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ToString(eh, args)
		},
	},

	"__time_now": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.TimeNow()
		},
	},

	"__date_now": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.DateNow()
		},
	},

	"__time_format_local": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.FormatTimeLocal(eh, args)
		},
	},

	"__time_format_utc": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.FormatTimeUTC(eh, args)
		},
	},

	"__os_user_name": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetUserName()
		},
	},

	"__os_user_homedir": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetUserHomeDir()
		},
	},

	"__osname": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetOS()
		},
	},

	"__osarch": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetArch()
		},
	},

	"__array_pop_without_index": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ArrayPopWithoutIndex(eh, args)
		},
	},

	"__array_pop_index": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ArrayPopIndex(eh, args)
		},
	},
	"__array_join": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.JoinArrays(eh, args)
		},
	},

	"__array_insert": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.InsertToArray(eh, args)
		},
	},

	"__array_insert_asis": {

		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.InsertToArrayAsIs(eh, args)
		},
	},

	"__file_read": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ReadFile(eh, args)
		},
	},

	"__file_exists": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.FileDirExists(eh, args)
		},
	},

	"__file_create_empty": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.CreateEmptyFile(eh, args)
		},
	},

	"__file_write": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.WriteToFile(eh, args)
		},
	},
	"__file_delete": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.DeletePath(eh, args)
		},
	},

	"__file_rename": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.RenameFile(eh, args)
		},
	},

	"__file_is_file": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.IsAFile(eh, args)
		},
	},

	"__file_is_dir": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.IsADir(eh, args)
		},
	},

	"__file_append_line": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.AppendLineToFile(eh, args)
		},
	},

	"__file_list_dir": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ListDir(eh, args)
		},
	},

	"__return_error": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ReturnErrorString(eh, args)
		},
	},
	"__readline": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.ReadLine(eh, args)
		},
	},
	"__get_type": {
		Fn: func(eh *object.ErrorHelper, env *object.EnvMap, caller token.Token, args ...object.Obj) object.Obj {
			return stdlib.GetType(eh, args)
		},
	},
}
