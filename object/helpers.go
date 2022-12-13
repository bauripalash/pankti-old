package object

import "go.cs.palashbauri.in/pankti/number"

func MakeIntNumber(i int64) Obj {
	return &Number{
		Value: number.MakeInt(i),
	}
}
