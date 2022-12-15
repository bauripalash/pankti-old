package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/token"
)

type HashKey struct {
	Type  ObjType
	Value uint64
	Token token.Token
}

func (b *Boolean) HashKey() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 2
	}

	return HashKey{Type: b.Type(), Value: val}
}

func (s *String) HashKey() HashKey {

	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (n *Number) HashKey() HashKey {
	h := fnv.New64a()
	if n.Value.IsInt {
		k := n.Value.Value.(*number.IntNumber).Value
		h.Write(k.Bytes())
		return HashKey{Type: n.Type(), Value: h.Sum64()}
	} else {
		k := n.Value.Value.(*number.FloatNumber).Value
		temp, _ := k.Int64()
		h.Write([]byte(k.String() + fmt.Sprintf("%d", temp)))
		return HashKey{Type: n.Type(), Value: h.Sum64()}
		//        h.Write()
	}

}

//Hash Pair { a : b}

type HashPair struct {
	Key   Obj
	Value Obj
}

type Hash struct {
	Pairs map[HashKey]HashPair
	Token token.Token
}

func (h *Hash) Type() ObjType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, p := range h.Pairs {
		pairs = append(
			pairs,
			fmt.Sprintf("%s : %s", p.Key.Inspect(), p.Value.Inspect()),
		)

	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (h *Hash) GetToken() token.Token {
	return h.Token
}

type Hashable interface {
	HashKey() HashKey
}
