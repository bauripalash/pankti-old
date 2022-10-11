package evaluator

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"pankti/ast"
	"pankti/lexer"
	"pankti/number"
	"pankti/object"
	"pankti/parser"
	"pankti/token"
	"strconv"
	"strings"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type ErrorHelper struct {
	Source string
}

func (e *ErrorHelper) GetLine(t token.Token) string {
	return strings.Split(e.Source, "\n")[t.LineNo]
}

func (e *ErrorHelper) MakeErrorLine(t token.Token) string {
	//    fmt.Println(t.LineNo , line)
	Lindex := t.Column - 1

	RIndex := t.Column + len(t.Literal) - 1

	if len(t.Literal) <= 1 {
		RIndex = Lindex + 1
	}

	newL := e.Source[:RIndex] + " <-- " + e.Source[RIndex:]
	newLine := newL[:Lindex] + " --> " + newL[Lindex:]
	return strconv.Itoa(t.LineNo) + "| " + newLine
}

func Eval(node ast.Node, env *object.Env, eh ErrorHelper) object.Obj {
	switch node := node.(type) {
	case *ast.Program:
		return evalProg(node, env, &eh)
	case *ast.ExprStmt:
		//fmt.Println("Eval Expr => ", node.Expr)
		return Eval(node.Expr, env, eh)
	case *ast.Boolean:
		return getBoolObj(node.Value)
	case *ast.NumberLit:
		return &object.Number{Value: node.Value, IsInt: node.IsInt, Token: node.Token}
	case *ast.PrefixExpr:
		r := Eval(node.Right, env, eh)
		if isErr(r) {
			return r
		}
		return evalPrefixExpr(node.Op, r, &eh)
	case *ast.InfixExpr:
		l := Eval(node.Left, env, eh)
		if isErr(l) {
			return l
		}
		r := Eval(node.Right, env, eh)
		if isErr(r) {
			return r
		}
		return evalInfixExpr(node.Op, l, r, &eh)
	case *ast.IfExpr:
		return evalIfExpr(node, env, &eh)
	case *ast.WhileExpr:
		return evalWhileExpr(node, env, &eh)
	case *ast.ReturnStmt:
		val := Eval(node.ReturnVal, env, eh)
		if isErr(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BlockStmt:
		return evalBlockStmt(node, env, &eh)
	case *ast.LetStmt:
		val := Eval(node.Value, env, eh)
		if isErr(val) {
			return val
		}

		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalId(node, env, &eh)
	case *ast.FunctionLit:
		pms := node.Params
		body := node.Body
		return &object.Function{Params: pms, Body: body, Env: env, Token: node.Token}
	case *ast.CallExpr:
		fnc := Eval(node.Func, env, eh)
		if isErr(fnc) {
			return fnc
		}
		//fmt.Println(node.Fun)
		args := evalExprs(node.Args, env, &eh)
		if len(args) == 1 && isErr(args[0]) {
			return args[0]
		}

		return applyFunc(fnc, args, &eh)

	case *ast.StringLit:
		return &object.String{Value: node.Value}
	case *ast.ArrLit:
		elms := evalExprs(node.Elms, env, &eh)
		if len(elms) == 1 && isErr(elms[0]) {
			return elms[0]
		}
		/*
		   func (p *Parser) parseCallArgs() []ast.Expr {
		   	args := []ast.Expr{}

		   	if p.isPeekToken(token.RPAREN) {
		   		p.nextToken()
		   		return args
		   	}

		   	p.nextToken()
		   	args = append(args, p.parseExpr(LOWEST))

		   	for p.isPeekToken(token.COMMA) {
		   		p.nextToken()
		   		p.nextToken()
		   		args = append(args, p.parseExpr(LOWEST))
		   	}

		   	if !p.peek(token.RPAREN) {
		   		return nil
		   	}

		   	return args
		   }
		*/
		return &object.Array{Elms: elms, Token: node.Token}

	case *ast.IndexExpr:
		left := Eval(node.Left, env, eh)
		if isErr(left) {
			return nil
		}

		index := Eval(node.Index, env, eh)
		if isErr(index) {
			return index
		}

		return evalIndexExpr(left, index, &eh)
	case *ast.HashLit:
		return evalHashLit(node, env, &eh)
	case *ast.IncludeStmt:
		//ImportMap.Env = *env
		//fmt.Println(env)
		newEnv, val := evalIncludeStmt(node, env, &eh) /*
			func (p *Parser) parseCallArgs() []ast.Expr {
				args := []ast.Expr{}

				if p.isPeekToken(token.RPAREN) {
					p.nextToken()
					return args
				}

				p.nextToken()
				args = append(args, p.parseExpr(LOWEST))

				for p.isPeekToken(token.COMMA) {
					p.nextToken()
					p.nextToken()
					args = append(args, p.parseExpr(LOWEST))
				}

				if !p.peek(token.RPAREN) {
					return nil
				}

				return args
			}
		*/
		//fmt.Println(env)
		if val.Type() != object.ERR_OBJ {
			*env = *object.NewEnclosedEnv(newEnv)
		} else {
			return val
		}
		//*env = *e
		//env = copy(env, e)
	}

	return nil
}

func evalHashLit(node *ast.HashLit, env *object.Env, eh *ErrorHelper) object.Obj {
	pairs := make(map[object.HashKey]object.HashPair)

	for kNode, vNode := range node.Pairs {

		key := Eval(kNode, env, *eh)

		if isErr(key) {
			return key
		}
		hashkey, ok := key.(object.Hashable)

		if !ok {
			return NewErr(node.Token, eh, "object cannot be used as hash key %s", key.Type())
		}

		val := Eval(vNode, env, *eh)

		if isErr(val) {
			return val
		}

		hashed := hashkey.HashKey()

		pairs[hashed] = object.HashPair{Key: key, Value: val}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpr(left, index object.Obj, eh *ErrorHelper) object.Obj {

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUM_OBJ:
		return evalArrIndexExpr(left, index, eh)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpr(left, index, eh)

	default:
		return NewErr(left.GetToken(), eh, "Unsupported Index Operator %s ", left.Type())
	}

}

func evalHashIndexExpr(hash, index object.Obj, eh *ErrorHelper) object.Obj {

	hashO := hash.(*object.Hash)

	key, ok := index.(object.Hashable)

	if !ok {
		return NewErr(index.GetToken(), eh, "This cannot be used as hash key %s", index.Type())
	}

	pair, ok := hashO.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}

	return pair.Value
}

func evalArrIndexExpr(arr, index object.Obj, eh *ErrorHelper) object.Obj {
	arrObj := arr.(*object.Array)
	id := index.(*object.Number).Value

	idx, noerr := number.GetAsInt(id)

	if !noerr {
		return NewBareErr("Arr Index Failed")
	}
	max := int64(len(arrObj.Elms) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrObj.Elms[idx]
}

func applyFunc(fn object.Obj, args []object.Obj, eh *ErrorHelper) object.Obj {

	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Params) == len(args) {
			eEnv := extendFuncEnv(fn, args)
			evd := Eval(fn.Body, eEnv, *eh)
			return unwrapRValue(evd)
		} else {

			return NewErr(fn.GetToken(), eh, "Function call doesn't have required arguments provided; wanted = %d but got %d", len(fn.Params), len(args))
		}
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return NewBareErr("%s is not a function", fn.Type())

	}
}

func extendFuncEnv(fn *object.Function, args []object.Obj) *object.Env {
	env := object.NewEnclosedEnv(fn.Env)

	//if len(args) > 0 {
	for pId, param := range fn.Params {
		env.Set(param.Value, args[pId])
	}
	//}

	return env
}

func evalIncludeStmt(in *ast.IncludeStmt, e *object.Env, eh *ErrorHelper) (*object.Env, object.Obj) {
	rawFilename := Eval(in.Filename, e, *eh)
	enx := object.NewEnv()

	if rawFilename.Type() != object.STRING_OBJ {
		return enx, NewBareErr("include filename is invalid %s", rawFilename.Inspect())

	}

	includeFilename := rawFilename.(*object.String).Value

	_, err := os.Stat(includeFilename)

	if errors.Is(err, fs.ErrNotExist) {
		return enx, NewErr(in.Token, eh, "%s include file doesnot exists", includeFilename)

	}

	fdata, err := os.ReadFile(includeFilename)

	if err != nil {
		return enx, NewErr(in.Token, eh, "Failed to read include file %s", includeFilename)

	}

	l := lexer.NewLexer(string(fdata))
	p := parser.NewParser(&l)
	ex := object.NewEnv()
	prog := p.ParseProg()
	Eval(prog, ex, *eh)
	//fmt.Println(evd.Type())

	if len(p.GetErrors()) != 0 {
		for _, e := range p.GetErrors() {
			fmt.Println(e.String())
		}

		return enx, NewBareErr("Include file contains parsing errors")
	}

	return ex, &object.Null{}

}

func unwrapRValue(o object.Obj) object.Obj {
	if rv, ok := o.(*object.ReturnValue); ok {
		return rv.Value
	}

	return o

}

func evalExprs(es []ast.Expr, env *object.Env, eh *ErrorHelper) []object.Obj {
	var res []object.Obj

	for _, e := range es {
		ev := Eval(e, env, *eh)

		if isErr(ev) {
			return []object.Obj{ev}
		}

		res = append(res, ev)
	}

	return res
}

func evalId(node *ast.Identifier, env *object.Env, eh *ErrorHelper) object.Obj {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return NewBareErr(eh.MakeErrorLine(node.Token) + "\n" + "id not found : " + node.Value)
	//	return val
}

func NewErr(token token.Token, eh *ErrorHelper, format string, a ...interface{}) *object.Error {
	errMsg := eh.MakeErrorLine(token) + "\n" + fmt.Sprintf(format, a...)
	return &object.Error{Msg: errMsg}
}

func isErr(obj object.Obj) bool {
	if obj != nil {
		return obj.Type() == object.ERR_OBJ
	}

	return false
}

func evalBlockStmt(block *ast.BlockStmt, env *object.Env, eh *ErrorHelper) object.Obj {

	var res object.Obj

	for _, stmt := range block.Stmts {
		res = Eval(stmt, env, *eh)

		//fmt.Println("E_BS=> " , res)

		if res != nil {
			rtype := res.Type()
			if rtype == object.RETURN_VAL_OBJ || rtype == object.ERR_OBJ {
				//fmt.Println("RET => " ,  res)
				return res
			}
		}
	}
	//fmt.Println("EBS 2=>" ,res)
	return res
}

func evalProg(prog *ast.Program, env *object.Env, eh *ErrorHelper) object.Obj {
	var res object.Obj

	for _, stmt := range prog.Stmts {
		res = Eval(stmt, env, *eh)

		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}

	return res
}

func evalIfExpr(iex *ast.IfExpr, env *object.Env, eh *ErrorHelper) object.Obj {
	cond := Eval(iex.Cond, env, *eh)

	if isErr(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(iex.TrueBlock, env, *eh)
	} else if iex.ElseBlock != nil {
		return Eval(iex.ElseBlock, env, *eh)
	} else {
		return NULL
	}

}

func evalWhileExpr(wx *ast.WhileExpr, env *object.Env, eh *ErrorHelper) object.Obj {
	cond := Eval(wx.Cond, env, *eh)
	var result object.Obj
	if isErr(cond) {
		return cond
	}

	for isTruthy(cond) {
		result = Eval(wx.StmtBlock, env, *eh)
		cond = Eval(wx.Cond, env, *eh)
	}

	return result
}

func isTruthy(obj object.Obj) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpr(op string, l, r object.Obj, eh *ErrorHelper) object.Obj {
	//fmt.Println(l.Type() , r.Type())
	switch {
	case l.Type() == object.NUM_OBJ && r.Type() == object.NUM_OBJ:
		return evalNumInfixExpr(op, l, r, eh)
		//}
		//fmt.Println("FI-> ", l , r)
		//return NewErr("has Float")
	case l.Type() == object.STRING_OBJ && r.Type() == object.STRING_OBJ:
		return evalStringInfixExpr(op, l, r, eh)
	case op == "==":
		return getBoolObj(l == r)
	case op == "!=":
		return getBoolObj(l != r)
	case l.Type() != r.Type():
		return NewBareErr("Type mismatch:  %s %s %s ", l.Type(), op, r.Type())
	default:
		return NewBareErr("unknown Operator : %s %s %s", l.Type(), op, r.Type())
	}
}

func evalStringInfixExpr(op string, l, r object.Obj, eh *ErrorHelper) object.Obj {
	lval := l.(*object.String).Value
	rval := r.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: lval + rval}
	case "==":
		return getBoolObj(lval == rval)
	case "!=":
		return getBoolObj(lval != rval)
	default:
		return NewBareErr("Unknown Operator %s %s %s", l.Type(), op, r.Type())

	}

}

func evalNumInfixExpr(op string, l, r object.Obj, eh *ErrorHelper) object.Obj {

	lval := l.(*object.Number).Value
	rval := r.(*object.Number).Value

	//fmt.Println(lval.GetType() , rval.GetType())

	val, cval, noerr := number.NumberOperation(op, lval, rval)
	if val.Value != nil && noerr {
		return &object.Number{Value: val, IsInt: val.IsInt}
	} else if val.Value == nil && noerr {
		return getBoolObj(cval)
	} else {
		return NewBareErr("Unknown Operator for Numbers %s", op)
	}

}

func evalPrefixExpr(op string, right object.Obj, eh *ErrorHelper) object.Obj {
	switch op {
	case "!":
		return evalBangOp(right)
	case "-":
		return evalMinusPrefOp(right, eh)
	default:
		return NewBareErr("Unknown Operator : %s%s", op, right.Type())

	}
}

func evalMinusPrefOp(right object.Obj, eh *ErrorHelper) object.Obj {
	if right.Type() != object.NUM_OBJ {
		return NewBareErr("unknown Operator : -%s", right.Type())
	}
	num := right.(*object.Number)
	return &object.Number{Value: number.MakeNeg(num.Value), IsInt: num.IsInt}
}

func evalBangOp(r object.Obj) object.Obj {
	switch r {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func getBoolObj(inp bool) *object.Boolean {
	if inp {
		return TRUE
	} else {
		return FALSE
	}
}

func evalStmts(stmts []ast.Stmt, env *object.Env, eh *ErrorHelper) object.Obj {
	var res object.Obj

	for _, stmt := range stmts {
		res = Eval(stmt, env, *eh)

		if rvalue, ok := res.(*object.ReturnValue); ok {
			return rvalue.Value
		}
	}

	return res
}

func NewBareErr(format string, a ...interface{}) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, a...)}
}
