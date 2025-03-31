package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	mutable := make(map[string]bool)
	return &Environment{store: store, mutable: mutable}

}

type Environment struct {
	store   map[string]Object
	outer   *Environment
	mutable map[string]bool
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object, mutable bool) Object {
	e.store[name] = val
	e.mutable[name] = mutable
	return val
}

func (e *Environment) IsMutable(name string) bool {
	mutable, ok := e.mutable[name]
	if ok {
		return mutable
	}

	// If not found in current environment, check outer environment
	if e.outer != nil {
		return e.outer.IsMutable(name)
	}

	// Default to immutable if not found
	return false
}
