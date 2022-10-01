package evaluator

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"vabna/ast"
	"vabna/lexer"
	"vabna/number"
	"vabna/object"
	"vabna/parser"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Env) object.Obj {
	switch node := node.(type) {
	case *ast.Program:
		return evalProg(node, env)
	case *ast.ExprStmt:
		//fmt.Println("Eval Expr => ", node.Expr)
		return Eval(node.Expr, env)
	case *ast.Boolean:
		return getBoolObj(node.Value)
	case *ast.NumberLit:
		return &object.Number{Value: node.Value, IsInt: node.IsInt}
	case *ast.PrefixExpr:
		r := Eval(node.Right, env)
		if isErr(r) {
			return r
		}
		return evalPrefixExpr(node.Op, r)
	case *ast.InfixExpr:
		l := Eval(node.Left, env)
		if isErr(l) {
			return l
		}
		r := Eval(node.Right, env)
		if isErr(r) {
			return r
		}
		return evalInfixExpr(node.Op, l, r)
	case *ast.IfExpr:
		return evalIfExpr(node, env)
	case *ast.WhileExpr:
		return evalWhileExpr(node, env)
	case *ast.ReturnStmt:
		val := Eval(node.ReturnVal, env)
		if isErr(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BlockStmt:
		return evalBlockStmt(node, env)
	case *ast.LetStmt:
		val := Eval(node.Value, env)
		if isErr(val) {
			return val
		}

		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalId(node, env)
	case *ast.FunctionLit:
		pms := node.Params
		body := node.Body
		return &object.Function{Params: pms, Body: body, Env: env}
	case *ast.CallExpr:
		fnc := Eval(node.Func, env)
		if isErr(fnc) {
			return fnc
		}
		//fmt.Println(node.Fun)
		args := evalExprs(node.Args, env)
		if len(args) == 1 && isErr(args[0]) {
			return args[0]
		}

		return applyFunc(fnc, args)

	case *ast.StringLit:
		return &object.String{Value: node.Value}
	case *ast.ArrLit:
		elms := evalExprs(node.Elms, env)
		if len(elms) == 1 && isErr(elms[0]) {
			return elms[0]
		}

		return &object.Array{Elms: elms}

	case *ast.IndexExpr:
		left := Eval(node.Left, env)
		if isErr(left) {
			return nil
		}

		index := Eval(node.Index, env)
		if isErr(index) {
			return index
		}

		return evalIndexExpr(left, index)
	case *ast.HashLit:
		return evalHashLit(node, env)
	case *ast.IncludeStmt:
		//ImportMap.Env = *env
		//fmt.Println(env)
		newEnv,val := evalIncludeStmt(node, env)
		//fmt.Println(env)
        if val.Type() != object.ERR_OBJ{
            *env = *object.NewEnclosedEnv(newEnv)
        }else{
            return val
        }
		//*env = *e
		//env = copy(env, e)
	}

	return nil
}

func evalHashLit(node *ast.HashLit, env *object.Env) object.Obj {
	pairs := make(map[object.HashKey]object.HashPair)

	for kNode, vNode := range node.Pairs {

		key := Eval(kNode, env)

		if isErr(key) {
			return key
		}
		hashkey, ok := key.(object.Hashable)

		if !ok {
			return NewErr("object cannot be used as hash key %s", key.Type())
		}

		val := Eval(vNode, env)

		if isErr(val) {
			return val
		}

		hashed := hashkey.HashKey()

		pairs[hashed] = object.HashPair{Key: key, Value: val}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpr(left, index object.Obj) object.Obj {

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUM_OBJ:
		return evalArrIndexExpr(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpr(left, index)
	default:
		return NewErr("Unsupported Index Operator %s ", left.Type())
	}

}

func evalHashIndexExpr(hash, index object.Obj) object.Obj {

	hashO := hash.(*object.Hash)

	key, ok := index.(object.Hashable)

	if !ok {
		return NewErr("This cannot be used as hash key %s", index.Type())
	}

	pair, ok := hashO.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}

	return pair.Value
}

func evalArrIndexExpr(arr, index object.Obj) object.Obj {
	arrObj := arr.(*object.Array)
	id := index.(*object.Number).Value

	idx, noerr := number.GetAsInt(id)

	if !noerr {
		return NewErr("Arr Index Failed")
	}
	max := int64(len(arrObj.Elms) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrObj.Elms[idx]
}

func applyFunc(fn object.Obj, args []object.Obj) object.Obj {

	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Params) == len(args) {
			eEnv := extendFuncEnv(fn, args)
			evd := Eval(fn.Body, eEnv)
			return unwrapRValue(evd)
		} else {

			return NewErr("Function call doesn't have required arguments provided; wanted = %d but got %d", len(fn.Params), len(args))
		}
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return NewErr("%s is not a function", fn.Type())

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

func evalIncludeStmt(in *ast.IncludeStmt, e *object.Env) (*object.Env , object.Obj) {
	rawFilename := Eval(in.Filename, e)
	enx := object.NewEnv()

	if rawFilename.Type() != object.STRING_OBJ {
		return enx, NewErr("include filename is invalid %s", rawFilename.Inspect())
		
	}

	includeFilename := rawFilename.(*object.String).Value

	_, err := os.Stat(includeFilename)

	if errors.Is(err, fs.ErrNotExist) {
		return enx, NewErr("%s include file doesnot exists", includeFilename)
	
	}

	fdata, err := os.ReadFile(includeFilename)

	if err != nil {
		return enx, NewErr("Failed to read include file %s", includeFilename)
		
	}

	l := lexer.NewLexer(string(fdata))
	p := parser.NewParser(&l)
	ex := object.NewEnv()
	prog := p.ParseProg()
    Eval(prog, ex)
    //fmt.Println(evd.Type())
    
    if len(p.GetErrors()) != 0{
        for _,e := range p.GetErrors(){
            fmt.Println(e.String())
        }

        return enx , NewErr("Include file contains parsing errors")
    }

	return ex , &object.Null{}

}

func unwrapRValue(o object.Obj) object.Obj {
	if rv, ok := o.(*object.ReturnValue); ok {
		return rv.Value
	}

	return o

}

func evalExprs(es []ast.Expr, env *object.Env) []object.Obj {
	var res []object.Obj

	for _, e := range es {
		ev := Eval(e, env)

		if isErr(ev) {
			return []object.Obj{ev}
		}

		res = append(res, ev)
	}

	return res
}

func evalId(node *ast.Identifier, env *object.Env) object.Obj {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return NewErr("id not found : " + node.Value)
	//	return val
}

func NewErr(format string, a ...interface{}) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, a...)}
}

func isErr(obj object.Obj) bool {
	if obj != nil {
		return obj.Type() == object.ERR_OBJ
	}

	return false
}

func evalBlockStmt(block *ast.BlockStmt, env *object.Env) object.Obj {

	var res object.Obj

	for _, stmt := range block.Stmts {
		res = Eval(stmt, env)

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

func evalProg(prog *ast.Program, env *object.Env) object.Obj {
	var res object.Obj

	for _, stmt := range prog.Stmts {
		res = Eval(stmt, env)

		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}

	return res
}

func evalIfExpr(iex *ast.IfExpr, env *object.Env) object.Obj {
	cond := Eval(iex.Cond, env)

	if isErr(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(iex.TrueBlock, env)
	} else if iex.ElseBlock != nil {
		return Eval(iex.ElseBlock, env)
	} else {
		return NULL
	}

}

func evalWhileExpr(wx *ast.WhileExpr, env *object.Env) object.Obj {
	cond := Eval(wx.Cond, env)
	var result object.Obj
	if isErr(cond) {
		return cond
	}

	for isTruthy(cond) {
		result = Eval(wx.StmtBlock, env)
		cond = Eval(wx.Cond, env)
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

func evalInfixExpr(op string, l, r object.Obj) object.Obj {
	//fmt.Println(l.Type() , r.Type())
	switch {
	case l.Type() == object.NUM_OBJ && r.Type() == object.NUM_OBJ:
		return evalNumInfixExpr(op, l, r)
		//}
		//fmt.Println("FI-> ", l , r)
		//return NewErr("has Float")
	case l.Type() == object.STRING_OBJ && r.Type() == object.STRING_OBJ:
		return evalStringInfixExpr(op, l, r)
	case op == "==":
		return getBoolObj(l == r)
	case op == "!=":
		return getBoolObj(l != r)
	case l.Type() != r.Type():
		return NewErr("Type mismatch:  %s %s %s ", l.Type(), op, r.Type())
	default:
		return NewErr("unknown Operator : %s %s %s", l.Type(), op, r.Type())
	}
}

func evalStringInfixExpr(op string, l, r object.Obj) object.Obj {
	if op != "+" {
		return NewErr("Unknown Operator %s %s %s", l.Type(), op, r.Type())
	}

	lval := l.(*object.String).Value
	rval := r.(*object.String).Value
	return &object.String{Value: lval + rval}
}

func evalNumInfixExpr(op string, l, r object.Obj) object.Obj {

	lval := l.(*object.Number).Value
	rval := r.(*object.Number).Value

	//fmt.Println(lval.GetType() , rval.GetType())

	val, cval, noerr := number.NumberOperation(op, lval, rval)
	if val.Value != nil && noerr {
		return &object.Number{Value: val, IsInt: val.IsInt}
	} else if val.Value == nil && noerr {
		return getBoolObj(cval)
	} else {
		return NewErr("Unknown Operator for Numbers %s", op)
	}

}

func evalPrefixExpr(op string, right object.Obj) object.Obj {
	switch op {
	case "!":
		return evalBangOp(right)
	case "-":
		return evalMinusPrefOp(right)
	default:
		return NewErr("Unknown Operator : %s%s", op, right.Type())

	}
}

func evalMinusPrefOp(right object.Obj) object.Obj {
	if right.Type() != object.NUM_OBJ {
		return NewErr("unknown Operator : -%s", right.Type())
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

func evalStmts(stmts []ast.Stmt, env *object.Env) object.Obj {
	var res object.Obj

	for _, stmt := range stmts {
		res = Eval(stmt, env)

		if rvalue, ok := res.(*object.ReturnValue); ok {
			return rvalue.Value
		}
	}

	return res
}
