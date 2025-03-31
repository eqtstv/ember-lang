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

func (e *Environment) Set(name string, val Object, mutable bool) (Object, bool) {
	// if e.VariableExists(name) && !e.IsMutable(name) {
	// 	return &Error{Message: "Cannot assign to immutable variable: " + name}, false
	// }

	e.store[name] = val
	e.mutable[name] = mutable
	return val, true
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

func (e *Environment) VariableExists(name string) bool {
	_, exists := e.store[name]
	if exists {
		return true
	}

	if e.outer != nil {
		return e.outer.VariableExists(name)
	}

	return false
}
