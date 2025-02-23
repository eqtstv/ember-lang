package parser

import (
	"ember_lang/ast"
	"ember_lang/lexer"
	"ember_lang/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

func (parser *Parser) peekPrecedence() int {
	precedence, ok := precedences[parser.peekToken.Type]

	if !ok {
		return LOWEST
	}
	return precedence
}

func (parser *Parser) curPrecedence() int {
	precedence, ok := precedences[parser.curToken.Type]

	if !ok {
		return LOWEST
	}
	return precedence
}

type (
	PrefixParseFn func() ast.Expression
	InfixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]PrefixParseFn
	infixParseFns  map[token.TokenType]InfixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	// Prefix parse functions
	parser.prefixParseFns = make(map[token.TokenType]PrefixParseFn)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBooleanLiteral)
	parser.registerPrefix(token.FALSE, parser.parseBooleanLiteral)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)

	// Infix parse functions
	parser.infixParseFns = make(map[token.TokenType]InfixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NEQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: parser.curToken}

	value, err := strconv.ParseInt(parser.curToken.Literal, 0, 64)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", parser.curToken.Literal)
		parser.errors = append(parser.errors, message)
		return nil
	}

	literal.Value = value
	return literal
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
	}

	parser.nextToken()

	expression.Right = parser.parseExpression(PREFIX)

	return expression
}

func (parser *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.Boolean{Token: parser.curToken, Value: parser.curTokenIs(token.TRUE)}
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()

	expression := parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
	}

	return expression
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.curToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	parser.nextToken()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.peekTokenIs(token.ELSE) {
		parser.nextToken()

		if !parser.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = parser.parseBlockStatement()
	}

	return expression

}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: parser.curToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = parser.parseBlockStatement()

	return literal
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return identifiers
	}

	parser.nextToken()

	identifiers = append(identifiers, &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal})

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()

		identifiers = append(identifiers, &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal})
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
		Left:     left,
	}

	precedence := parser.curPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: parser.curToken, Function: function}

	expression.Arguments = parser.parseExpressionList(token.RPAREN)

	return expression
}

func (parser *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if parser.peekTokenIs(end) {
		parser.nextToken()
		return args
	}

	parser.nextToken()
	args = append(args, parser.parseExpression(LOWEST))

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()

		args = append(args, parser.parseExpression(LOWEST))
	}

	if parser.peekTokenIs(end) {
		parser.nextToken()
	}

	return args
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, message)
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}

}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn PrefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn InfixParseFn) {
	parser.infixParseFns[tokenType] = fn
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.curToken.Type != token.EOF {
		statement := parser.parseStatement()

		program.Statements = append(program.Statements, statement)

		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}

}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: parser.curToken,
		Value: parser.curToken.Literal,
	}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We are skipping the expressions until we encounter a semicolon
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: parser.curToken}

	parser.nextToken()

	// TODO: We are skipping the expressions until we encounter a semicolon
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser.curToken}

	statement.Expression = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFns[parser.curToken.Type]

	if prefix == nil {
		parser.noPrefixParseFnError(parser.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFns[parser.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		parser.nextToken()

		leftExp = infix(leftExp)

	}

	return leftExp
}

func (parser *Parser) noPrefixParseFnError(t token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s found", t)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.curToken}
	block.Statements = []ast.Statement{}

	parser.nextToken()

	for !parser.curTokenIs(token.RBRACE) && parser.curToken.Type != token.EOF {
		statement := parser.parseStatement()
		block.Statements = append(block.Statements, statement)

		parser.nextToken()
	}

	return block
}
