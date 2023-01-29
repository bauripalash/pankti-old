package object

import (
	"fmt"

	"go.cs.palashbauri.in/pankti/code"
	"go.cs.palashbauri.in/pankti/token"
)

const (
	INT_OBJ           = "INTEGER"
	FLOAT_OBJ         = "FLOAT"
	BOOL_OBJ          = "BOOLEAN"
	RETURN_VAL_OBJ    = "RETURN_VAL"
	NULL_OBJ          = "NIL"
	ERR_OBJ           = "ERROR"
	FUNC_OBJ          = "FUNCTION"
	STRING_OBJ        = "STRING"
	BUILTIN_OBJ       = "BUILTIN"
	ARRAY_OBJ         = "ARRAY"
	HASH_OBJ          = "HASH"
	NUM_OBJ           = "NUM"
	INCLUDE_OBJ       = "INCLUDE"
	SHOW_OBJ          = "SHOW"
	BREAK_OBJ         = "BREAK"
	COMPILED_FUNC_OBJ = "COMPILED_FUNC_OBJ"
)

type BuiltInFunc func(eh *ErrorHelper, env *EnvMap, caller token.Token, args ...Obj) Obj

type ObjType string

type Obj interface {
	Type() ObjType
	Inspect() string
	GetToken() token.Token
}

type CompiledFunc struct {
	Instructions code.Instructions
}

func (*CompiledFunc) Type() ObjType { return COMPILED_FUNC_OBJ }
func (cf *CompiledFunc) Inspect() string {
	return fmt.Sprintf("COMPILED_FUNC[%p]", cf)
}
func (*CompiledFunc) GetToken() token.Token { return token.Token{} }

type Builtin struct {
	Fn    BuiltInFunc
	Token token.Token
}

func (*Builtin) Type() ObjType           { return BUILTIN_OBJ }
func (*Builtin) Inspect() string         { return "builtin function" }
func (b *Builtin) GetToken() token.Token { return b.Token }
