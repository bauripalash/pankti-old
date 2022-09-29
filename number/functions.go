package number

import (
	"math/big"
	"vabna/token"
)

func MakeInt(a int64) Number{
    return Number{ Value: &IntNumber{ Value: *big.NewInt(a) } }
}

func MakeFloat(a float64) Number{
    return Number{Value: &FloatNumber{ Value: *big.NewFloat(a) }  }
}

func MakeNeg(a Number) Number{
    
    if a.IsInt{
        ia := a.Value.(*IntNumber).Value
        r := new(big.Int).Neg(&ia)
        return Number{ Value: &IntNumber{ Value: *r } , IsInt: true }
    }else{
        fa := a.Value.(*FloatNumber).Value
        r := new(big.Float).Neg(&fa)
        return Number{Value: &FloatNumber{Value: *r} , IsInt: false}
    }

}

func GetAsInt(a Number) (int64, bool){

    if a.IsInt{
        ia := a.Value.(*IntNumber).Value

        if ia.IsInt64(){
            return ia.Int64(),true
        }
    }else{
        fa := a.Value.(*FloatNumber).Value

        a,_ := fa.Int64()

        return a,true
    }

    return int64(0), false

}

func FloatFloatCompare(op string, a big.Float, b big.Float) bool {
	switch op {
	case ">":
		res := a.Cmp(&b)
		if res == 1 {
			return true
		}
	case "<":
		res := a.Cmp(&b)
		if res == -1 {
			return true
		}
	case "==":
		res := a.Cmp(&b)
		if res == 0 {
			return true
		}
	case "!=":
		res := a.Cmp(&b)
		if res != 0 {
			return true
		}
	case ">=":
		res := a.Cmp(&b)
		if res == 1 || res == 0 {
			return true
		}
	case "<=":
		res := a.Cmp(&b)
		if res == -1 || res == 0 {
			return true
		}
	}

	return false
}

func IntIntCompare(op string, a big.Int, b big.Int) bool {
	switch op {
	case ">":
		res := a.Cmp(&b)
		if res == 1 {
			return true
		}
	case "<":
		res := a.Cmp(&b)
		if res == -1 {
			return true
		}
	case "==":
		res := a.Cmp(&b)
		if res == 0 {
			return true
		}
	case "!=":
		res := a.Cmp(&b)
		if res != 0 {
			return true
		}
	case ">=":
		res := a.Cmp(&b)
		if res == 1 || res == 0 {
			return true
		}
	case "<=":
		res := a.Cmp(&b)
		if res == -1 || res == 0 {
			return true
		}
	}

	return false
}

func NumberOperation(op string, n Number, x Number) (Number, bool, bool) {

	var fb *big.Float
	var fa *big.Float
	var val *big.Float
	switch n.GetType() {
	case "FLOAT":
		fa = &n.Value.(*FloatNumber).Value
		//var fb  *big.Float

		//var val *big.Float
		if x.IsInt {
			b := x.Value.(*IntNumber).Value

			fb = new(big.Float).SetInt(&b)
		} else {
			fb = &x.Value.(*FloatNumber).Value
		}

		switch op {
		case token.PLUS:
			val = new(big.Float).Add(fa, fb)
		case token.MINUS:
			val = new(big.Float).Sub(fa, fb)
		case token.MUL:
			val = new(big.Float).Mul(fa, fb)
		case token.DIV:
			val = new(big.Float).Quo(fa, fb)
		case token.GT:
			return Number{}, FloatFloatCompare(token.GT, *fa, *fb), true
		case token.GTE:
			return Number{}, FloatFloatCompare(token.GTE, *fa, *fb), true
		case token.LT:
			return Number{}, FloatFloatCompare(token.LT, *fa, *fb), true
		case token.LTE:
			return Number{}, FloatFloatCompare(token.LTE, *fa, *fb), true
		case token.NOT_EQ:
			return Number{}, FloatFloatCompare(token.NOT_EQ, *fa, *fb), true
		case token.EQEQ:
			return Number{}, FloatFloatCompare(token.EQEQ, *fa, *fb), true

		}
		return Number{Value: &FloatNumber{Value: *val}, IsInt: false}, false, true
	case "INT":
		ia := n.Value.(*IntNumber).Value

		if x.IsInt {
			ib := x.Value.(*IntNumber).Value
			var val *big.Int
			switch op {
			case token.PLUS:
				val = new(big.Int).Add(&ia, &ib)
			case token.MINUS:
				val = new(big.Int).Sub(&ia, &ib)
			case token.MUL:
				val = new(big.Int).Mul(&ia, &ib)
			case token.DIV:
				val = new(big.Int).Div(&ia, &ib)
			case token.GT:
				return Number{}, IntIntCompare(token.GT, ia, ib), true
			case token.GTE:
				return Number{}, IntIntCompare(token.GTE, ia, ib), true
			case token.LT:
				return Number{}, IntIntCompare(token.LT, ia, ib), true
			case token.LTE:
				return Number{}, IntIntCompare(token.LTE, ia, ib), true
			case token.NOT_EQ:
				return Number{}, IntIntCompare(token.NOT_EQ, ia, ib), true
			case token.EQEQ:
				return Number{}, IntIntCompare(token.EQEQ, ia, ib), true

			}

			return Number{Value: &IntNumber{Value: *val}, IsInt: true}, false, true

		}

		fb = &x.Value.(*FloatNumber).Value
		fa = new(big.Float).SetInt(&ia)
		switch op {
		case token.PLUS:
			val = new(big.Float).Add(fa, fb)
		case token.MINUS:
			val = new(big.Float).Sub(fa, fb)
		case token.MUL:
			val = new(big.Float).Mul(fa, fb)
		case token.DIV:
			val = new(big.Float).Quo(fa, fb)
		case token.GT:
			return Number{}, FloatFloatCompare(token.GT, *fa, *fb), true
		case token.GTE:
			return Number{}, FloatFloatCompare(token.GTE, *fa, *fb), true
		case token.LT:
			return Number{}, FloatFloatCompare(token.LT, *fa, *fb), true
		case token.LTE:
			return Number{}, FloatFloatCompare(token.LTE, *fa, *fb), true
		case token.NOT_EQ:
			return Number{}, FloatFloatCompare(token.NOT_EQ, *fa, *fb), true
		case token.EQEQ:
			return Number{}, FloatFloatCompare(token.EQEQ, *fa, *fb), true

		}
		return Number{Value: &FloatNumber{Value: *val}, IsInt: false}, false, true
	default:
		return Number{}, false, false
	}
}
