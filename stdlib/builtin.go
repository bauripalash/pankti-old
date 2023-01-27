package stdlib

import (
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/token"
)

func getHashKey(i object.Obj) (object.HashKey, bool) {

	switch i := i.(type) {
	case *object.String:
		return i.HashKey(), true
	case *object.Boolean:
		return i.HashKey(), true
	case *object.Number:
		return i.HashKey(), true
	default:
		return object.HashKey{}, false

	}
}

/*
func GetObjFromHashKey(i object.HashKey) (object.HashKey , bool){
	return object.String{} , false
}
*/

func SetHashTableElm(eh *object.ErrorHelper, env *object.EnvMap, t token.Token, args []object.Obj) object.Obj {
	if len(args) != 3 {
		return object.NewErr(t, eh, true, "setting hash table value required three arguments")
	}

	rawHashTble := args[0]
	hashKey := args[1]
	newValue := args[2]

	if rawHashTble.Type() != object.HASH_OBJ {
		return object.NewErr(t, eh, true, "this function only works in HashTable")
	}

	hashTble := rawHashTble.(*object.Hash)
	key, isokay := getHashKey(hashKey)

	if !isokay {
		return object.NewErr(t, eh, true, "Key must be hashable -> string/boolean/number")
	}
	newHp := object.HashPair{Key: hashKey, Value: newValue}
	newHashTable := &object.Hash{Token: hashTble.Token, Pairs: map[object.HashKey]object.HashPair{}}

	for k, v := range hashTble.Pairs {
		newHashTable.Pairs[k] = v
	}

	newHashTable.Pairs[key] = newHp

	return newHashTable
}

func GetAllKVsOfHashTable(forKeys bool, eh *object.ErrorHelper, env *object.EnvMap, t token.Token, args []object.Obj) object.Obj {

	rawHashTble := args[0]
	ht := rawHashTble.(*object.Hash)
	arr := []object.Obj{}
	for _, v := range ht.Pairs {
		if forKeys {
			arr = append(arr, v.Key)
		} else {
			arr = append(arr, v.Value)
		}
	}

	return &object.Array{Elms: arr}

}
