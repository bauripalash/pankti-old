package evaluator

import (
	"bytes"

	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
)

func evalHashLit(
	node *ast.HashLit,
	env *object.EnvMap,
	eh *object.ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {

	pairs := make(map[object.HashKey]object.HashPair)

	for kNode, vNode := range node.Pairs {

		key := Eval(kNode, env, *eh, printBuff, isGui)

		if object.IsErr(key) {
			return key
		}
		hashkey, ok := key.(object.Hashable)

		if !ok {
			return object.NewErr(
				node.Token,
				eh,
				true,
				"object cannot be used as hash key %s",
				key.Type(),
			)
		}

		val := Eval(vNode, env, *eh, printBuff, isGui)

		if object.IsErr(val) {
			return val
		}

		hashed := hashkey.HashKey()

		pairs[hashed] = object.HashPair{Key: key, Value: val}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpr(left, index object.Obj, eh *object.ErrorHelper) object.Obj {

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUM_OBJ:
		return evalArrIndexExpr(left, index, eh)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpr(left, index, eh)

	default:
		return object.NewErr(
			left.GetToken(),
			eh,
			true,
			"Unsupported Index Operator %s ",
			left.Type(),
		)
	}

}

func evalHashIndexExpr(
	hash, index object.Obj,
	eh *object.ErrorHelper,
) object.Obj {

	hashO := hash.(*object.Hash)

	key, ok := index.(object.Hashable)

	if !ok {
		return object.NewErr(
			index.GetToken(),
			eh,
			true,
			"This cannot be used as hash key %s",
			index.Type(),
		)
	}

	pair, ok := hashO.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}

	return pair.Value
}

func evalArrIndexExpr(arr, index object.Obj, eh *object.ErrorHelper) object.Obj {
	arrObj := arr.(*object.Array)
	id := index.(*object.Number).Value

	idx, noerr := number.GetAsInt(id)

	if !noerr {
		return object.NewBareErr("Arr Index Failed")
	}
	max := int64(len(arrObj.Elms) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrObj.Elms[idx]
}
