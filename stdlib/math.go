package stdlib

import (
	"math"

	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
)

var ARG_NOT_FLOAT = &object.Error{Msg: "Arguments must be Numbers as the provided can not be parsed as decimal number"}

var NOT_ALL_INT = &object.Error{Msg: "All arguments to GCD must be integers"}

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

func getInt(arg object.Obj) (int64, bool) {
	rawInput := arg

	if rawInput.Type() != object.NUM_OBJ {
		return 0, false
	}

	ival := rawInput.(*object.Number).Value

	if ival.IsInt {
		iv := ival.Value.(*number.IntNumber).Value
		return iv.Int64(), true
	} else {
		fv := ival.Value.(*number.FloatNumber).Value

		if fv.IsInt() {
			return 0, false
		}

		v, _ := fv.Int64()
		return v, true
	}

	//return 0,false
}

func DoListSum(args []object.Obj) object.Obj {
	if args[0].Type() != object.ARRAY_OBJ {
		return &object.Error{Msg: "Sum can be only done on Arrays"}
	}

	inputList := args[0].(*object.Array).Elms
	result := float64(0)
	for _, item := range inputList {
		if fv, ok := getFloat(item); ok {
			result += fv
		} else {
			return ARG_NOT_FLOAT
		}
	}

	return object.MakeFloatNumber(result)
}

func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// Unstable
func GetGCD(args []object.Obj) object.Obj {
	if len(args) < 2 {
		return &object.Error{Msg: "Must provide at least two arguments to calculate GCD"}
	}
	temp, ok := getInt(args[0])

	if !ok {
		return NOT_ALL_INT
	}

	for _, item := range args[1:] {
		b, ok2 := getInt(item)
		if !ok2 {
			return NOT_ALL_INT
		}
		temp = gcd(temp, b)
	}

	return object.MakeIntNumber(temp)
}

func lcm(a, b int64) int64 {
	return (a * b / gcd(a, b))
}

func GetLCM(args []object.Obj) object.Obj {
	if len(args) < 2 {
		return &object.Error{Msg: "Must provide at least two arguments to calculate GCD"}
	}

	tempA, ok := getInt(args[0])
	if !ok {
		return NOT_ALL_INT
	}
	tempB, okB := getInt(args[1])
	if !okB {
		return NOT_ALL_INT
	}

	result := tempA * tempB / gcd(tempA, tempB)

	for _, item := range args[2:] {
		b, ok := getInt(item)
		if !ok {
			return NOT_ALL_INT
		}
		result = lcm(result, b)
	}

	return object.MakeIntNumber(result)
}

func DoSqrt(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	fValue, ok := getFloat(args[0])
	if !ok {
		//fmt.Println(args[0].GetToken())
		return object.NewErr(args[0].GetToken(), eh, true, "Provided Argument must be a Number")
	}
	return object.MakeFloatNumber(math.Sqrt(fValue))

}

func DoPow(args []object.Obj) object.Obj {
	if len(args) != 2 {
		return &object.Error{Msg: "Pow must have two arguments"}
	}
	if fv, ok := getFloat(args[0]); ok {
		if fv2, ok2 := getFloat(args[1]); ok2 {
			return object.MakeFloatNumber(math.Pow(fv, fv2))
		}
	}

	return ARG_NOT_FLOAT
}

func Log10(args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {
		return ARG_NOT_FLOAT
	}
	return object.MakeFloatNumber(math.Log10(fv))
}

func LogE(args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {
		return ARG_NOT_FLOAT
	}
	return object.MakeFloatNumber(math.Log(fv))
}

func LogX(args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		if len(args) == 1 {
			return object.MakeFloatNumber(math.Log(fv))
		} else if len(args) >= 2 {
			if base, okay := getFloat(args[1]); okay {
				return object.MakeFloatNumber(math.Log(fv) / math.Log(base))
			}
		}

	}

	return ARG_NOT_FLOAT
}

func Cosine(args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {
		return ARG_NOT_FLOAT
	}

	return object.MakeFloatNumber(math.Cos(fv))
}

func Acos(args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		return object.MakeFloatNumber(math.Acos(fv))
	}
	return ARG_NOT_FLOAT
}

func Sine(args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {
		return ARG_NOT_FLOAT
	}

	return object.MakeFloatNumber(math.Sin(fv))
}

func Asin(args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		return object.MakeFloatNumber(math.Asin(fv))
	}
	return ARG_NOT_FLOAT
}

func Tangent(args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		return object.MakeFloatNumber(math.Tan(fv))
	}

	return ARG_NOT_FLOAT
}

func Atan(args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		return object.MakeFloatNumber(math.Atan(fv))
	}
	return ARG_NOT_FLOAT
}

func Atan2(args []object.Obj) object.Obj {
	if len(args) != 2 {
		return &object.Error{Msg: "There must be 2 arguments to Atan2"}
	}

	if fv, ok := getFloat(args[0]); ok {
		if fv2, ok2 := getFloat(args[1]); ok2 {
			return object.MakeFloatNumber(math.Atan2(fv, fv2))
		}
	}

	return ARG_NOT_FLOAT
}

func ToDegree(args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		dg := fv * (180 / math.Pi)
		return object.MakeFloatNumber(dg)
	}
	return ARG_NOT_FLOAT
}

func ToRadians(args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		rad := fv * (math.Pi / 180)
		return object.MakeFloatNumber(rad)
	}

	return ARG_NOT_FLOAT
}

func GetPI(_ []object.Obj) object.Obj {
	return object.MakeFloatNumber(float64(math.Pi))
}

func GetE(_ []object.Obj) object.Obj {
	return object.MakeFloatNumber(float64(math.E))
}
