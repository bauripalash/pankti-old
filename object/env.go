package object

const DEFKEY = "__default"

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

type EnvMap struct {
	Envs map[string]Env
}

func NewEnvMap() *EnvMap {
	e := EnvMap{Envs: make(map[string]Env)}
	e.Envs[DEFKEY] = *NewEnv()
	return &e
}

func (em *EnvMap) GetDefaultEnv() (Env, bool) {
	e, ok := em.Envs[DEFKEY]
	return e, ok
}

func (em *EnvMap) GetEnv(envName string) (Env, bool) {

	e, ok := em.Envs[envName]
	return e, ok
}

func (em *EnvMap) GetFromDefault(key string) (Obj, bool) {
	x, ok := em.Envs[DEFKEY]

	if !ok {
		return &Null{}, ok
	}

	return x.Get(key)
}

func (em *EnvMap) GetFrom(envName, key string) (Obj, bool) {
	x, ok := em.Envs[envName]

	if !ok {
		return &Null{}, ok
	}

	return x.Get(key)
}

func (em *EnvMap) SetToDefault(key string, value Obj) Obj {
	x, ok := em.Envs[DEFKEY]

	if !ok {
		return &Null{}
	}

	return x.Set(key, value)
}

func (em *EnvMap) SetTo(envName string, key string, value Obj) Obj {
	x, ok := em.Envs[envName]

	if !ok {
		return &Null{}
	}

	return x.Set(key, value)
}

func (em *EnvMap) CreateEmptyEnv(envName string) *Env {
	e := NewEnv()
	em.Envs[envName] = *e
	return e
}

func (em *EnvMap) MergeEnv(envName string, env *Env) {
	_, ok := em.Envs[envName]

	if !ok {
		em.CreateEmptyEnv(envName)
	}

	em.Envs[envName] = *env
}

func (em *EnvMap) EnvExists(envName string) bool {
	_, ok := em.Envs[envName]
	return ok
}
