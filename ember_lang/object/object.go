package object

import (
	"bytes"
	"ember_lang/ember_lang/ast"
	"fmt"
	"hash/fnv"
	"strings"
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
	HASH_OBJ         ObjectType = "HASH"
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

// ----------------------------------------------------------------------------
// Array Object
// ----------------------------------------------------------------------------

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

// ----------------------------------------------------------------------------
// HashKey Object
// ----------------------------------------------------------------------------

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var result uint64

	if b.Value {
		result = 1
	}

	return HashKey{Type: b.Type(), Value: result}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: hash.Sum64()}
}

// ----------------------------------------------------------------------------
// Hash Object
// ----------------------------------------------------------------------------

type Hash struct {
	Pairs map[HashKey]HashPair
}

type HashPair struct {
	Key   Object
	Value Object
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}
