package stdlib

import (
	"strings"

	"go.cs.palashbauri.in/pankti/errs"
	"go.cs.palashbauri.in/pankti/object"
)

func SplitString(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	target, ok := getStringFromArgs(args[0])

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["NOT_ALL_STRING"])
	}

	delim, ok2 := getStringFromArgs(args[1])

	if !ok2 {

		return object.NewErr(args[1].GetToken(), eh, true, errs.Errs["NOT_ALL_STRING"])
	}

	var result object.Array

	for _, item := range strings.Split(target, delim) {
		result.Elms = append(result.Elms, &object.String{Value: item})
	}

	return &result

}

func JoinAsString(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	target, ok := getArrayFromArg(args[0])

	if !ok {

		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["NOT_ALL_ARE_LIST"])
	}

	delim, ok2 := getStringFromArgs(args[1])

	if !ok2 {

		return object.NewErr(args[1].GetToken(), eh, true, errs.Errs["NOT_ALL_STRING"])
	}

	result := ""

	for index, item := range target {
		result += item.Inspect()
		if (index + 1) != len(target) {
			result += delim
		}
	}

	return &object.String{Value: result}

}

func ToString(_ *object.ErrorHelper, args []object.Obj) object.Obj {
	//if args[0].Type() != object.NULL_OBJ{
	return &object.String{Value: args[0].Inspect()}
	//}

}
