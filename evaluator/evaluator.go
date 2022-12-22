package evaluator

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/constants"
	"go.cs.palashbauri.in/pankti/lexer"
	"go.cs.palashbauri.in/pankti/number"
	"go.cs.palashbauri.in/pankti/object"
	"go.cs.palashbauri.in/pankti/parser"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(
	node ast.Node,
	env *object.EnvMap,
	eh object.ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {
	switch node := node.(type) {
	case *ast.Program: //Entry point of a Program AST
		return evalProg(node, env, &eh, printBuff, isGui)
	case *ast.ExprStmt:
		return Eval(node.Expr, env, eh, printBuff, isGui)
	case *ast.Boolean:
		return getBoolObj(node.Value)
	case *ast.NumberLit:
		return &object.Number{Value: node.Value, IsInt: node.IsInt, Token: node.Token}
	case *ast.PrefixExpr:
		//
		// Prefix ->
		// <Operator> Expression
		//
		r := Eval(node.Right, env, eh, printBuff, isGui) // Evaluate the Expression to the smallest possible value
		if object.IsErr(r) {
			return r
		}
		return evalPrefixExpr(node.Op, r, &eh)
	case *ast.InfixExpr:
		//
		// Infix ->
		// Left_Expression <Operator> Right_Expression
		//
		l := Eval(node.Left, env, eh, printBuff, isGui) // Evaluate the <Left_Expression> to the smallest possible value
		if object.IsErr(l) {
			return l
		}
		r := Eval(node.Right, env, eh, printBuff, isGui) // Evaluate the <Right_Expression> to the smallest possible value
		if object.IsErr(r) {
			return r
		}
		return evalInfixExpr(node.Op, l, r, &eh)
	case *ast.IfExpr:
		return evalIfExpr(node, env, &eh, printBuff, isGui)
	case *ast.WhileExpr:
		return evalWhileExpr(node, env, &eh, printBuff, isGui)
	case *ast.ReturnStmt:
		val := Eval(node.ReturnVal, env, eh, printBuff, isGui)
		if object.IsErr(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.ShowStmt:
		//
		// Show / Print Statement
		// here lies the real usage of `printBuff` and `isGui`
		//
		args := evalExprs(node.Value, env, &eh, printBuff, isGui)
		return evalShowStmt(args, printBuff, isGui)
	case *ast.BlockStmt:
		return evalBlockStmt(node, env, &eh, printBuff, isGui)
	case *ast.LetStmt:
		return evalLetStmt(node, env, &eh, printBuff, isGui)
	case *ast.Identifier:
		return evalId(node, env, &eh)
	//case *ast.IncludeId:
	//	return evalIncludeId(node, env, &eh)
	case *ast.FunctionLit:
		// Function Declaration
		// dhori X = ekti kaj() ...... sesh
		pms := node.Params
		body := node.Body
		e, _ := env.GetDefaultEnv()
		return &object.Function{Name: node.TokenLit(), Params: pms, Body: body, Env: &e, Token: node.Token}
	case *ast.CallExpr:
		//
		// Function Call
		// function()
		//
		isMod := len(strings.Split(node.Func.String(), ".")) == 2

		fnc := Eval(node.Func, env, eh, printBuff, isGui)
		//fmt.Println(isMod)
		if object.IsErr(fnc) {
			return fnc
		}
		args := evalExprs(node.Args, env, &eh, printBuff, isGui)
		if len(args) == 1 && object.IsErr(args[0]) {
			return args[0]
		}

		return applyFunc(fnc, node.Token, args, isMod, env, &eh, printBuff, isGui)

	case *ast.StringLit:
		return &object.String{Value: node.Value, Token: node.Token}
	case *ast.ArrLit:
		elms := evalExprs(node.Elms, env, &eh, printBuff, isGui)
		if len(elms) == 1 && object.IsErr(elms[0]) {
			return elms[0]
		}

		return &object.Array{Elms: elms, Token: node.Token}

	case *ast.IndexExpr:
		// node.Left ==> The Array --> ARRAY[index]
		left := Eval(node.Left, env, eh, printBuff, isGui)
		if object.IsErr(left) {
			return nil
		}
		// node.Index ==> The Array Index --> array[INDEX]
		index := Eval(node.Index, env, eh, printBuff, isGui)
		if object.IsErr(index) {
			return index
		}

		return evalIndexExpr(left, index, &eh)
	case *ast.HashLit:
		return evalHashLit(node, env, &eh, printBuff, isGui)
	case *ast.IncludeExpr:
		return &object.IncludeObj{Filename: node.Filename.String()}
	}

	return nil
}

func evalId(
	node *ast.Identifier,
	env *object.EnvMap,
	eh *object.ErrorHelper,
) object.Obj {
	if val, ok := env.GetFromDefault(node.Value); ok {
		return val
	}

	if node.IsMod {
		keys := strings.Split(node.Value, ".")
		envName := keys[0]
		envId := keys[1]
		if val, ok := env.GetFrom(envName, envId); ok {
			return val
		}
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return object.NewErr(node.Token, eh, true, "id not found : "+node.Value)
	//	return val
}

func evalProg(
	prog *ast.Program,
	env *object.EnvMap,
	eh *object.ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {
	var res object.Obj

	for _, stmt := range prog.Stmts {
		res = Eval(stmt, env, *eh, printBuff, isGui)

		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}

	return res
}

func evalLetStmt(
	node *ast.LetStmt,
	env *object.EnvMap,
	eh *object.ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {

	if node.Name.IsMod {
		return object.NewBareErr("Dot notation can not be used directly")
	}

	val := Eval(node.Value, env, *eh, printBuff, isGui)
	//fmt.Println(val)

	if val.Type() == object.FUNC_OBJ {
		v := val.(*object.Function)
		v.Name = node.Name.Value
	}

	if val.Type() == object.INCLUDE_OBJ {
		//fmt.Println(val.Inspect())
		iobj := val.(*object.IncludeObj)
		evaluateInclude(env, eh, printBuff, isGui, node.Name.Value, iobj.Filename)

		val = &object.String{Value: iobj.Filename}
	}

	//fmt.Println(node.Value.String())
	if object.IsErr(val) {
		return val
	}

	env.SetToDefault(node.Name.Value, val)
	return &object.Null{}
}

func evaluateInclude(env *object.EnvMap,
	eh *object.ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
	key string, filename string) {
	//e := object.NewEnv()
	fname := filename
	if filepath.IsAbs(filename) {
		fname = filename
	} else {
		cwd, _ := os.Getwd()
		p := filepath.Join(cwd, filename)
		if _, err := os.Stat(p); !errors.Is(err, fs.ErrNotExist) {
			//fmt.Println("xx->" , p)
			fname = p
		}

		import_env := os.Getenv(constants.IMPORT_PATH_ENV)
		//fmt.Println(import_env , filename)
		if len(import_env) >= 1 && fname != p {
			enName, ok := constants.GetStdName(filename)
			if ok {
				filename = enName
			}
			//fmt.Println(os.Stat(filepath.Join(import_env, filename)))
			if _, err := os.Stat(filepath.Join(import_env, filename)); !errors.Is(err, fs.ErrNotExist) {
				fname = filepath.Join(import_env, filename)
			}
		}

	}

	//fmt.Println("IP=>>" + fname)
	_, err := os.Stat(fname)

	if errors.Is(err, fs.ErrNotExist) {

		fmt.Println("Not exists file")

	}

	fdata, _ := os.ReadFile(fname)

	l := lexer.NewLexer(string(fdata))
	p := parser.NewParser(&l)
	ex := object.NewEnvMap()
	prog := p.ParseProg()
	Eval(prog, ex, *eh, printBuff, isGui)
	x, _ := ex.GetDefaultEnv()

	env.MergeEnv(key, &x)
	//fmt.Println(key, filename)
}

func evalMinusPrefOp(right object.Obj, eh *object.ErrorHelper) object.Obj {
	if right.Type() != object.NUM_OBJ {
		//return object.NewBareErr("unknown Operator : -%s", right.Type())
		return object.NewErr(right.GetToken(), eh, true, "Unknown operator")
	}
	num := right.(*object.Number)
	return &object.Number{
		Value: number.MakeNeg(num.Value),
		IsInt: num.IsInt,
	}
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
