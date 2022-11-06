package object

import "bauri.palash/pankti/number"

func MakeIntNumber(i int64) Obj {
	return &Number{Value: number.MakeInt(i)}
}
