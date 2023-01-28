package compiler

import (
	"fmt"

	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/code"
	"go.cs.palashbauri.in/pankti/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Obj
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Obj
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Obj{},
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
	case *ast.NumberLit:
		i := &object.Number{Value: node.Value}
		c.emit(code.OpConstant, c.addConst(i))
		//fmt.Println(object.Number { Value: node.Value })
		//switch node := node.(type){

		//}
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}

	return nil

}

func (c *Compiler) addConst(o object.Obj) int {
	c.constants = append(c.constants, o)
	return len(c.constants) - 1
}

func (c *Compiler) addIns(ins []byte) int {
	pos := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return pos
}

func (c *Compiler) emit(op code.OpCode, oprs ...int) int {
	ins := code.Make(op, oprs...)
	pos := c.addIns(ins)
	return pos
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
