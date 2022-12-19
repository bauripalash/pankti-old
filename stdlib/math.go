package stdlib

import (
	"math"

	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
)

func getFloat(arg object.Obj) (float64, bool) {
	rawinput := arg

	if rawinput.Type() != object.NUM_OBJ {
		return 0.0, false
	}
	fValue, ok := rawinput.(*object.Number).Value.GetAsFloat()

	if !ok {
		return 0.0, false
	}

	return fValue, ok
}

func DoSqrt(args []object.Obj) object.Obj {
	fValue, ok := getFloat(args[0])

	if !ok {
		return &object.Error{Msg: "Arguments can not be parsed as decimal number"}
	}

	result := math.Sqrt(fValue)

	//bFloat := new(big.Float).SetFloat64(result)

	return &object.Number{Value: number.MakeFloat(result)}

}

func Log10(args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {

		return &object.Error{Msg: "Arguments can not be parsed as decimal number"}
	}
	return &object.Number{Value: number.MakeFloat(math.Log10(fv))}
}
