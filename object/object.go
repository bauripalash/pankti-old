package object

import (
	"go.cs.palashbauri.in/pankti/token"
)

const (
	INT_OBJ        = "INTEGER"
	FLOAT_OBJ      = "FLOAT"
	BOOL_OBJ       = "BOOLEAN"
	RETURN_VAL_OBJ = "RETURN_VAL"
	NULL_OBJ       = "NIL"
	ERR_OBJ        = "ERROR"
	FUNC_OBJ       = "FUNCTION"
	STRING_OBJ     = "STRING"
	BUILTIN_OBJ    = "BUILTIN"
	ARRAY_OBJ      = "ARRAY"
	HASH_OBJ       = "HASH"
	NUM_OBJ        = "NUM"
	INCLUDE_OBJ    = "INCLUDE"
	SHOW_OBJ       = "SHOW"
)

type BuiltInFunc func(eh *ErrorHelper, caller token.Token, args ...Obj) Obj

type ObjType string

type Obj interface {
	Type() ObjType
	Inspect() string
	GetToken() token.Token
}

type Builtin struct {
	Fn    BuiltInFunc
	Token token.Token
}

func (_ *Builtin) Type() ObjType         { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string       { return "builtin function" }
func (b *Builtin) GetToken() token.Token { return b.Token }
