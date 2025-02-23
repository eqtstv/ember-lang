package object

import (
	"bytes"
	"ember_lang/ember_lang/ast"
	"fmt"
)

// ----------------------------------------------------------------------------
// Object Types
// ----------------------------------------------------------------------------

type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	STRING_OBJ       ObjectType = "STRING"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
)

// ----------------------------------------------------------------------------
// Object Interface
// ----------------------------------------------------------------------------

type Object interface {
	Type() ObjectType
	Inspect() string
}

// ----------------------------------------------------------------------------
// Integer Object
// ----------------------------------------------------------------------------
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// ----------------------------------------------------------------------------
// Boolean Object
// ----------------------------------------------------------------------------

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// ----------------------------------------------------------------------------
// Null Object
// ----------------------------------------------------------------------------

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

// ----------------------------------------------------------------------------
// ReturnValue Object
// ----------------------------------------------------------------------------

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

// ----------------------------------------------------------------------------
// Error Object
// ----------------------------------------------------------------------------

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "\033[31mERROR: " + e.Message + "\033[0m"
}

// ----------------------------------------------------------------------------
// Function Object
// ----------------------------------------------------------------------------

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("fn")
	out.WriteString("(")

	for i, param := range f.Parameters {
		out.WriteString(param.String())
		if i != len(f.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// ----------------------------------------------------------------------------
// String Object
// ----------------------------------------------------------------------------

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

// ----------------------------------------------------------------------------
// Builtin Object
// ----------------------------------------------------------------------------

type Builtin struct {
	Fn BuiltinFunction
}

type BuiltinFunction func(args ...Object) Object

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

// ------------------------------------- Array Object -------------------------------------

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	out.WriteString("[")
	for i, e := range a.Elements {
		out.WriteString(e.Inspect())
		if i != len(a.Elements)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("]")

	return out.String()
}
