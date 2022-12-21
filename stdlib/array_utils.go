package stdlib

import (
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
)

func getArrayFromArg(arg object.Obj) ([]object.Obj, bool) {
	if arg.Type() != object.ARRAY_OBJ {
		return []object.Obj{}, false
	}
	xa := arg.(*object.Array).Elms
	return xa, true
}

func getIntFromArg(arg object.Obj) (int64, bool) {
	if arg.Type() != object.NUM_OBJ {
		return 0, false
	}

	xa := arg.(*object.Number).Value

	if xa.IsInt {
		temp := xa.Value.(*number.IntNumber)
		return temp.Value.Int64(), true
	} else {
		temp := xa.Value.(*number.FloatNumber).Value
		if temp.IsInt() {
			t, _ := temp.Int64()
			return t, true
		} else {
			return 0, false
		}

	}
}

func ArrayPopWithoutIndex(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	elms, ok := getArrayFromArg(args[0])

	if !ok {
		object.NewErr(args[0].GetToken(), eh, true, "Argument must be list")
	}

	return &object.Array{Elms: elms[:len(elms)-1]}
}

func ArrayPopIndex(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	elms, ok := getArrayFromArg(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, "Argument must be list")
	}

	i, ok2 := getIntFromArg(args[1])

	if ok2 {
		return object.NewErr(args[1].GetToken(), eh, true, "Index Argument must be list")
	}
	if i >= int64(len(elms)) {
		return object.NewErr(args[1].GetToken(), eh, true, "Index out of range")
	}
	s := append(elms[:i], elms[i+1:]...)

	return &object.Array{
		Elms: s,
	}

}

func JoinArrays(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	x_elms, ok := getArrayFromArg(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, "Argument must be a list")
	}

	y_elms, ok2 := getArrayFromArg(args[2])

	if !ok2 {

		return object.NewErr(args[1].GetToken(), eh, true, "Argument must be a list")
	}

	x_elms = append(x_elms, y_elms...)

	return &object.Array{Elms: x_elms}
}

func InsertToArray(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	x_elms, ok := getArrayFromArg(args[0])

	if !ok {

		return object.NewErr(args[0].GetToken(), eh, true, "Argument must be a list")
	}

	y := args[1]
	index, ok2 := getIntFromArg(args[2])

	if !ok2 {

		return object.NewErr(args[2].GetToken(), eh, true, "Index Argument must be list")
	}
	if index >= int64(len(x_elms)) {
		return object.NewErr(y.GetToken(), eh, true, "Index out of range")
	}
	if y.Type() == object.ARRAY_OBJ {
		y_elms, oky := getArrayFromArg(y)

		if !oky {

			return object.NewErr(y.GetToken(), eh, true, "Argument must be a list")
		}

		result := append(x_elms[:index], append(y_elms, x_elms[index+1:]...)...)

		return &object.Array{
			Elms: result,
		}

	} else {
		result := x_elms[:index]
		result = append(result, y)
		result = append(result, x_elms[index+1:]...)
		return &object.Array{Elms: result}
	}

}

func InsertToArrayAsIs(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	x_elms, ok := getArrayFromArg(args[0])

	if !ok {

		return object.NewErr(args[0].GetToken(), eh, true, "Argument must be a list")
	}

	y := args[1]
	index, ok2 := getIntFromArg(args[2])

	if !ok2 {

		return object.NewErr(args[2].GetToken(), eh, true, "Index Argument must be list")
	}

	if index >= int64(len(x_elms)) {
		return object.NewErr(y.GetToken(), eh, true, "Index out of range")
	}

	result := x_elms[:index]
	result = append(result, y)
	result = append(result, x_elms[index+1:]...)
	return &object.Array{Elms: result}
}
