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
