package stdlib

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"go.cs.palashbauri.in/pankti/constants"
	"go.cs.palashbauri.in/pankti/errs"
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/token"
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

func DoListSum(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["SUM_ONLY_LISTS"])
	}

	inputList := args[0].(*object.Array)
	result := float64(0)
	for _, item := range inputList.Elms {
		if fv, ok := getFloat(item); ok {
			result += fv
		} else {
			return object.NewErr(item.GetToken(), eh, true, errs.Errs["SUM_ARRAY_ALL_NUM"])
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
func GetGCD(eh *object.ErrorHelper, _ token.Token, args []object.Obj) object.Obj {

	temp, ok := getInt(args[0])

	if !ok {
		return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["NOT_ALL_ARE_INT"], constants.FNames["gcd"])
	}

	for index, item := range args[1:] {
		b, ok2 := getInt(item)
		if !ok2 {
			return object.NewErr(args[index].GetToken(), eh, true, errs.Errs["NOT_ALL_ARE_INT"], constants.FNames["gcd"])
		}
		temp = gcd(temp, b)
	}

	return object.MakeIntNumber(temp)
}

func lcm(a, b int64) int64 {
	return (a * b / gcd(a, b))
}

func GetLCM(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	tempA, ok := getInt(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ALL_INT"])
	}
	tempB, okB := getInt(args[1])
	if !okB {
		//eturn NOT_ALL_INT

		return object.NewErr(args[1].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ALL_INT"])
	}

	result := tempA * tempB / gcd(tempA, tempB)

	for _, item := range args[2:] {
		b, ok := getInt(item)
		if !ok {
			return object.NewErr(item.GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ALL_INT"])
		}
		result = lcm(result, b)
	}

	return object.MakeIntNumber(result)
}

func DoSqrt(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	fValue, ok := getFloat(args[0])
	if !ok {
		//fmt.Println(args[0].GetToken())
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ALL_INT"], constants.FNames["sqrt"], constants.FNames["int"])
	}
	return object.MakeFloatNumber(math.Sqrt(fValue))

}

func DoPow(eh *object.ErrorHelper, args []object.Obj) object.Obj {

	if fv, ok := getFloat(args[0]); ok {
		if fv2, ok2 := getFloat(args[1]); ok2 {
			return object.MakeFloatNumber(math.Pow(fv, fv2))
		} else {

			return object.NewErr(args[1].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ALL_INT"], constants.FNames["pow"], constants.FNames["num"])
		}

	} else {

		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ALL_INT"], constants.FNames["pow"], constants.FNames["num"])
	}

}

func Log10(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["log"], constants.FNames["num"])
	}
	return object.MakeFloatNumber(math.Log10(fv))
}

func LogE(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["log"], constants.FNames["num"])
	}
	return object.MakeFloatNumber(math.Log(fv))
}

func LogX(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		if len(args) == 1 {
			return object.MakeFloatNumber(math.Log(fv))
		} else if len(args) >= 2 {
			if base, okay := getFloat(args[1]); okay {
				return object.MakeFloatNumber(math.Log(fv) / math.Log(base))
			} else {

				return object.NewErr(args[1].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["log"], constants.FNames["num"])
			}

		}

	} else {

		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["log"], constants.FNames["num"])
	}

	return object.NewErr(args[0].GetToken(), eh, false, errs.Errs["TEMPLATE_NOT_ALL_INT"], constants.FNames["log"], constants.FNames["num"])
}

func Cosine(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {
		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["cos"], constants.FNames["num"])
	}

	return object.MakeFloatNumber(math.Cos(fv))
}

func Acos(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		return object.MakeFloatNumber(math.Acos(fv))
	}

	return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["cos"], constants.FNames["num"])
}

func Sine(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	fv, ok := getFloat(args[0])
	if !ok {

		return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["sin"], constants.FNames["num"])
	}

	return object.MakeFloatNumber(math.Sin(fv))
}

func Asin(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		return object.MakeFloatNumber(math.Asin(fv))
	}
	return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["sin"], constants.FNames["num"])
}

func Tangent(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		return object.MakeFloatNumber(math.Tan(fv))
	}

	return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["tan"], constants.FNames["num"])
}

func Atan(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		return object.MakeFloatNumber(math.Atan(fv))
	}

	return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["tan"], constants.FNames["num"])

}

func Atan2(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		if fv2, ok2 := getFloat(args[1]); ok2 {
			return object.MakeFloatNumber(math.Atan2(fv, fv2))
		} else {

			return object.NewErr(args[1].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ALL_INT"], constants.FNames["tan"], constants.FNames["num"])
		}
	}

	return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ALL_INT"], constants.FNames["tan"], constants.FNames["num"])

}

func ToDegree(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		dg := fv * (180 / math.Pi)
		return object.MakeFloatNumber(dg)
	}

	return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["degree"], constants.FNames["num"])

}

func ToRadians(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	if fv, ok := getFloat(args[0]); ok {
		rad := fv * (math.Pi / 180)
		return object.MakeFloatNumber(rad)
	}

	return object.NewErr(args[0].GetToken(), eh, true, errs.Errs["TEMPLATE_NOT_ONE_TEMPALTE"], constants.FNames["rad"], constants.FNames["num"])

}

func ToNumber(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	result := float64(0)
	target := args[0]
	switch target.Type() {
	case object.BOOL_OBJ:
		if t := target.(*object.Boolean).Value; t {
			result = 1.0
		} else {
			result = 0.0
		}
	case object.STRING_OBJ:
		t := target.(*object.String).Value

		v, err := strconv.ParseFloat(t, 64)

		if err != nil {
			return object.NewErr(target.GetToken(), eh, true, errs.Errs["CANNOT_PARSE_STRING_AS_NUM"])
		}

		result = v
	case object.NUM_OBJ:
		v, _ := getFloat(target)
		result = v

	default:
		return object.NewErr(target.GetToken(), eh, true, errs.Errs["CANNOT_PARSE_AS_NUM"])

	}

	return object.MakeFloatNumber(result)
}

func ConvertToFloat(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	return ToNumber(eh, args)
}

func ConvertToInt(eh *object.ErrorHelper, args []object.Obj) object.Obj {
	result := int64(0)
	target := args[0]
	switch target.Type() {
	case object.BOOL_OBJ:
		if t := target.(*object.Boolean).Value; t {
			result = 1
		} else {
			result = 0
		}
	case object.STRING_OBJ:
		t := target.(*object.String).Value

		v, err := strconv.Atoi(t)

		if err != nil {
			return object.NewErr(target.GetToken(), eh, true, errs.Errs["CANNOT_PARSE_STRING_AS_NUM"])
		}

		result = int64(v)
	case object.NUM_OBJ:
		v, _ := getInt(target)
		result = v

	default:
		return object.NewErr(target.GetToken(), eh, true, errs.Errs["CANNOT_PARSE_AS_NUM"])

	}

	return object.MakeIntNumber(result)

}

func GenerateRandom(eh *object.ErrorHelper, args []object.Obj) object.Obj {

	rand.Seed(time.Now().UnixNano())
	n, ok := getIntFromArg(args[0])
	if !ok {
		object.NewErr(args[0].GetToken(), eh, true, errs.Errs["CANNOT_PARSE_AS_NUM"])
	}

	return object.MakeIntNumber(int64(rand.Intn(int(n))))
	//rand.Intn()

}

func GetPI(_ []object.Obj) object.Obj {
	return object.MakeFloatNumber(float64(math.Pi))
}

func GetE(_ []object.Obj) object.Obj {
	return object.MakeFloatNumber(float64(math.E))
}
