package ast

import (
	"bytes"
	"ember_lang/ember_lang/token"
	"strings"
)

// ------------------------------------- Base Types -------------------------------------

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// ------------------------------------- Identifier -------------------------------------

type Identifier struct {
	Token   token.Token // token.IDENTIFIER token
	Value   string
	Mutable bool
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

// ------------------------------------- Program -------------------------------------

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ------------------------------------- LetStatement -------------------------------------

type LetStatement struct {
	Token token.Token // token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.Value)
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ------------------------------------- ReturnStatement -------------------------------------

type ReturnStatement struct {
	Token       token.Token // token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ------------------------------------- ExpressionStatement -------------------------------------

type ExpressionStatement struct {
	Token      token.Token // token.IDENTIFIER token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// ------------------------------------- BlockStatement -------------------------------------

type BlockStatement struct {
	Token      token.Token // token.LBRACE token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ------------------------------------- IntegerLiteral -------------------------------------

type IntegerLiteral struct {
	Token token.Token // token.INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// ------------------------------------- PrefixExpression -------------------------------------

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. ! or -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// ------------------------------------- InfixExpression -------------------------------------

type InfixExpression struct {
	Token    token.Token // The infix token, e.g. + or -
	Operator string
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// ------------------------------------- Boolean -------------------------------------

type Boolean struct {
	Token token.Token // token.TRUE or token.FALSE
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

// ------------------------------------- IfExpression -------------------------------------

type IfExpression struct {
	Token       token.Token // token.IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.TokenLiteral() + " ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" : ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// ------------------------------------- FunctionLiteral -------------------------------------

type FunctionLiteral struct {
	Token      token.Token // token.FUNCTION token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// ------------------------------------- CallExpression -------------------------------------

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// ------------------------------------- StringLiteral -------------------------------------

type StringLiteral struct {
	Token token.Token // token.STRING token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// ------------------------------------- ArrayLiteral -------------------------------------

type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	for i, e := range al.Elements {
		out.WriteString(e.String())
		if i != len(al.Elements)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("]")

	return out.String()
}

// ------------------------------------- IndexExpression -------------------------------------

type IndexExpression struct {
	Token token.Token // token.LBRACKET token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// ------------------------------------- HashLiteral -------------------------------------

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {

}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}

	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// ------------------------------------- IncrementExpression -------------------------------------

type IncrementExpression struct {
	Token token.Token // token.INCREMENT token
	Left  Expression
}

func (ie *IncrementExpression) expressionNode() {}

func (ie *IncrementExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IncrementExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.TokenLiteral())
	out.WriteString(ie.Left.String())

	return out.String()
}

// ------------------------------------- WhileExpression -------------------------------------

type WhileExpression struct {
	Token     token.Token // token.WHILE token
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode() {}

func (we *WhileExpression) TokenLiteral() string {
	return we.Token.Literal
}

func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString(we.TokenLiteral() + " ")
	out.WriteString(we.Condition.String())
	out.WriteString(" {")
	out.WriteString(we.Body.String())
	out.WriteString("}")

	return out.String()
}

// ------------------------------------- ForExpression -------------------------------------

type ForExpression struct {
	Token        token.Token // token.FOR token
	LetStatement *LetStatement
	Condition    Expression
	Increment    Expression
	Body         *BlockStatement
}

func (fe *ForExpression) expressionNode() {}

func (fe *ForExpression) TokenLiteral() string {
	return fe.Token.Literal
}

func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString(fe.TokenLiteral() + " ")
	out.WriteString(fe.LetStatement.String())
	out.WriteString("; ")
	out.WriteString(fe.Condition.String())
	out.WriteString("; ")
	out.WriteString(fe.Increment.String())
	out.WriteString(" {")
	out.WriteString(fe.Body.String())
	out.WriteString("}")

	return out.String()
}

// ------------------------------------- AssignmentExpression -------------------------------------
type AssignmentExpression struct {
	Token token.Token // The '=' token
	Left  Expression
	Right Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Left.String())
	out.WriteString(" = ")
	out.WriteString(ae.Right.String())

	return out.String()
}
