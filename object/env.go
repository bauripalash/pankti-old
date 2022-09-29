package object

type Env struct {
	str   map[string]Obj
	outer *Env
}

func NewEnv() *Env {
	s := make(map[string]Obj)
	return &Env{str: s, outer: nil}
}

func (e *Env) Get(n string) (Obj, bool) {
	val, ok := e.str[n]

	if !ok && e.outer != nil {
		val, ok = e.outer.Get(n)
	}

	return val, ok
}

func (e *Env) Set(n string, v Obj) Obj {
	e.str[n] = v
	return v
}

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}
