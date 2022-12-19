package object

import "go.cs.palashbauri.in/pankti/number"

func MakeIntNumber(i int64) Obj {
	return &Number{
		Value: number.MakeInt(i),
	}
}

func MakeFloatNumber(i float64) Obj {
	return &Number{
		Value: number.MakeFloat(i),
	}
}
