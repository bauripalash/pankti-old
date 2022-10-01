package object

import "vabna/number"

func MakeIntNumber(i int64) Obj {
	return &Number{Value: number.MakeInt(i)}
}
