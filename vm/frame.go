package vm

import (
	"go.cs.palashbauri.in/pankti/code"
	"go.cs.palashbauri.in/pankti/object"
)

type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

func NewFrame(fn *object.Closure, bp int) *Frame {
	f := &Frame{
		cl:          fn,
		ip:          -1,
		basePointer: bp,
	}
	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
