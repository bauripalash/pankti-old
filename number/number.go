package number

import (
	"math/big"
	"strings"
)

type NumberType string

type Num interface {
	String() string
	Type() NumberType
}

type FloatNumber struct {
	Value big.Float
}

func (f *FloatNumber) String() string {
	return f.Value.String()
}

func (*FloatNumber) Type() NumberType {
	return "F"
}

type IntNumber struct {
	Value big.Int
}

func (i *IntNumber) String() string {
	return i.Value.String()
}

func (*IntNumber) Type() NumberType {
	return "I"
}

type Number struct {
	Value Num
	IsInt bool
}

func IsFloat(inp string) bool {
	return strings.ContainsRune(inp, rune('.'))
}

func (n *Number) SetValue(v string) bool {
	if IsFloat(v) {
		temp := new(big.Float)
		f, noerr := temp.SetString(v)
		if noerr {
			n.Value = &FloatNumber{Value: *f}
			n.IsInt = false

		}

		return noerr
	} else {
		temp := new(big.Int)
		i, noerr := temp.SetString(v, 10)

		if noerr {
			n.Value = &IntNumber{Value: *i}
			n.IsInt = true
		}

		return noerr
	}

}

func (n *Number) GetType() string {
	if n.IsInt {
		return "INT"
	} else {
		return "FLOAT"
	}
}
