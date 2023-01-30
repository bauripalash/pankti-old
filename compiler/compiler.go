package compiler

import (
	"fmt"
	"sort"

	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/code"
	"go.cs.palashbauri.in/pankti/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Obj

	//lastIns EmittedIns
	//prevIns EmittedIns

	symTable *SymbolTable

	scopes     []CompScope
	scopeIndex int
}

type CompScope struct {
	instructions code.Instructions
	lastIns      EmittedIns
	prevIns      EmittedIns
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Obj
}

type EmittedIns struct {
	OpCode code.OpCode
	Pos    int
}

func NewCompiler() *Compiler {
	defScope := CompScope{
		instructions: code.Instructions{},
		lastIns:      EmittedIns{},
		prevIns:      EmittedIns{},
	}
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Obj{},
		symTable:     NewSymbolTable(),
		scopes:       []CompScope{defScope},
		scopeIndex:   0,
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Stmts {
			err := c.Compile(s)

			if err != nil {
				return err
			}
		}
	case *ast.ExprStmt:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.PrefixExpr:
		if err := c.Compile(node.Right); err != nil {
			return err
		}

		switch node.Op {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("Unknown operator %s", node.Op)
		}
	case *ast.InfixExpr:

		if node.Op.Literal == "<" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			if err := c.Compile(node.Left); err != nil {
				return err

			}

			c.emit(code.OpGT)
			return nil
		}

		if err := c.Compile(node.Left); err != nil {
			return err
		}

		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Op.Literal {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGT)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)

		default:
			return fmt.Errorf("Unknown operator %s", node.Op.Literal)

		}

	case *ast.IfExpr:
		//fmt.Println(node.TrueBlock)
		if err := c.Compile(node.Cond); err != nil {
			return err
		}

		jntpos := c.emit(code.OpJumpNotTruthy, 9999)

		if err := c.Compile(node.TrueBlock); err != nil {
			return err
		}

		if c.lastInsIs(code.OpPop) {
			c.removeLastPop()
		}

		jmpPos := c.emit(code.OpJump, 9999)
		afteTBPos := len(c.currentIns())
		c.changeOperand(jntpos, afteTBPos)
		//fmt.Println(node.ElseBlock)
		if len(node.ElseBlock.Stmts) < 1 {
			c.emit(code.OpNull)
			//fmt.Println("ELSE+>NULL")
			//afterTBPOS := len(c.instructions)
			//c.changeOperand(jntpos , afterTBPOS)
		} else {
			//afteTBPos := len(c.instructions)
			//c.changeOperand(jntpos , afteTBPos)

			if err := c.Compile(node.ElseBlock); err != nil {
				return err
			}

			if c.lastInsIs(code.OpPop) {
				c.removeLastPop()
			}

		}
		afterEBPos := len(c.currentIns())
		c.changeOperand(jmpPos, afterEBPos)

	case *ast.LetStmt:
		if err := c.Compile(node.Value); err != nil {
			return err
		}

		sm := c.symTable.Define(node.Name.Value)
		if sm.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, sm.Index)
		} else {
			c.emit(code.OpSetLocal, sm.Index)
		}
		//c.emit(code.OpSetGlobal, sm.Index)
	case *ast.BlockStmt:
		for _, st := range node.Stmts {
			if err := c.Compile(st); err != nil {
				return err
			}
		}
	case *ast.Identifier:
		s, ok := c.symTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		//	c.emit(code.OpGetGlobal, s.Index)
		/*if s.Scope == GlobalScope {
			c.emit(code.OpGetGlobal, s.Index)
		} else {
			c.emit(code.OpGetLocal, s.Index)
		}*/
		c.loadSymbol(s)
	case *ast.NumberLit:
		i := &object.Number{Value: node.Value}
		c.emit(code.OpConstant, c.addConst(i))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.ArrLit:
		for _, el := range node.Elms {
			if err := c.Compile(el); err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elms))
	case *ast.StringLit:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConst(str))
	case *ast.HashLit:
		keys := []ast.Expr{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}

		//		fmt.Println(keys)

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			if err := c.Compile(k); err != nil {
				return err
			}

			if err := c.Compile(node.Pairs[k]); err != nil {
				return err
			}

		}

		c.emit(code.OpHash, len(node.Pairs)*2)
	case *ast.IndexExpr:
		if err := c.Compile(node.Left); err != nil {
			return err
		}

		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(code.OpIndex)

	case *ast.FunctionLit:
		c.enterScope()

		for _, p := range node.Params {
			c.symTable.Define(p.Value)
		}

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		if c.lastInsIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInsIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}
		nL := c.symTable.numDef
		ins := c.exitScope()
		cFn := &object.CompiledFunc{Instructions: ins, NumLocals: nL, NumParams: len(node.Params)}
		fnIndex := c.addConst(cFn)
		c.emit(code.OpClosure, fnIndex, 0)
	case *ast.ReturnStmt:
		if err := c.Compile(node.ReturnVal); err != nil {
			return err
		}
		c.emit(code.OpReturnValue)
	case *ast.CallExpr:
		if err := c.Compile(node.Func); err != nil {
			return err
		}
		for _, a := range node.Args {
			if err := c.Compile(a); err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(node.Args))

	}

	return nil

}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	}
}

func (c *Compiler) currentIns() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastP := c.scopes[c.scopeIndex].lastIns.Pos
	c.replaceIns(lastP, code.Make(code.OpReturnValue))
	c.scopes[c.scopeIndex].lastIns.OpCode = code.OpReturnValue
}

func (c *Compiler) replaceIns(pos int, newIns []byte) {
	ins := c.currentIns()
	for i := 0; i < len(newIns); i++ {
		ins[pos+i] = newIns[i]
	}
}

func (c *Compiler) changeOperand(pos int, op int) {
	o := code.OpCode(c.currentIns()[pos])
	newIns := code.Make(o, op)
	c.replaceIns(pos, newIns)
}

func (c *Compiler) addConst(o object.Obj) int {
	c.constants = append(c.constants, o)
	return len(c.constants) - 1
}

func (c *Compiler) addIns(ins []byte) int {
	posNew := len(c.currentIns())
	updatedIns := append(c.currentIns(), ins...)
	c.scopes[c.scopeIndex].instructions = updatedIns
	// pos := len(c.instructions)
	// c.instructions = append(c.instructions, ins...)
	// return pos
	return posNew
}

func (c *Compiler) emit(op code.OpCode, oprs ...int) int {
	ins := code.Make(op, oprs...)
	pos := c.addIns(ins)
	c.setLastIns(op, pos)
	return pos
}

func (c *Compiler) setLastIns(op code.OpCode, pos int) {
	prev := c.scopes[c.scopeIndex].lastIns
	last := EmittedIns{OpCode: op, Pos: pos}
	c.scopes[c.scopeIndex].prevIns = prev
	c.scopes[c.scopeIndex].lastIns = last
}

func (c *Compiler) lastInsIs(op code.OpCode) bool {
	if len(c.currentIns()) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].lastIns.OpCode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastIns
	prev := c.scopes[c.scopeIndex].prevIns

	old := c.currentIns()
	nw := old[:last.Pos]

	c.scopes[c.scopeIndex].instructions = nw
	c.scopes[c.scopeIndex].lastIns = prev

	// c.instructions = c.instructions[:c.lastIns.Pos]
	// c.lastIns = c.prevIns
}

func (c *Compiler) enterScope() {
	scope := CompScope{
		instructions: code.Instructions{},
		lastIns:      EmittedIns{},
		prevIns:      EmittedIns{},
	}

	c.symTable = NewEncolsedSymbolTable(c.symTable)
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
}

func (c *Compiler) exitScope() code.Instructions {
	ins := c.currentIns()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symTable = c.symTable.Outer
	return ins
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.currentIns(),
		Constants:    c.constants,
	}
}
