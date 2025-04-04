package parser

import (
	"ember_lang/ember_lang/ast"
	"ember_lang/ember_lang/lexer"
	"ember_lang/ember_lang/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	ASSIGN      // =
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INCREMENT   // i++
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:        EQUALS,
	token.NEQ:       EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.LTE:       LESSGREATER,
	token.GTE:       LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTERISK:  PRODUCT,
	token.LPAREN:    CALL,
	token.INCREMENT: INCREMENT,
	token.LBRACKET:  INDEX,
	token.ASSIGN:    ASSIGN,
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
	parser.registerPrefix(token.PLUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBooleanLiteral)
	parser.registerPrefix(token.FALSE, parser.parseBooleanLiteral)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)
	parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral)
	parser.registerPrefix(token.LBRACE, parser.parseHashLiteral)
	parser.registerPrefix(token.WHILE, parser.parseWhileExpression)
	parser.registerPrefix(token.FOR, parser.parseForExpression)
	parser.registerPrefix(token.AMPERSAND, parser.parsePointerReferenceExpression)
	parser.registerPrefix(token.ASTERISK, parser.parsePointerDereferenceExpression)

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
	parser.registerInfix(token.LTE, parser.parseInfixExpression)
	parser.registerInfix(token.GTE, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)
	parser.registerInfix(token.INCREMENT, parser.parseIncrementExpression)
	parser.registerInfix(token.ASSIGN, parser.parseAssignmentExpression)

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
		message := fmt.Sprintf("(line %d) could not parse %q as integer", parser.curToken.LineNumber, parser.curToken.Literal)
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

func (parser *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: parser.curToken}

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

	expression.Body = parser.parseBlockStatement()

	return expression

}

func (parser *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: parser.curToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	parser.nextToken()
	expression.LetStatement = parser.parseForLoopLetStatement()

	if !parser.expectPeek(token.SEMICOLON) {
		return nil
	}

	parser.nextToken()

	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.SEMICOLON) {
		return nil
	}

	parser.nextToken()

	left := parser.parseIdentifier()
	parser.nextToken()
	expression.Increment = parser.parseIncrementExpression(left)

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = parser.parseBlockStatement()

	return expression
}

func (parser *Parser) parseForLoopLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	letStmt.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()

	letStmt.Value = parser.parseExpression(LOWEST)

	return letStmt
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

func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: parser.curToken, Value: parser.curToken.Literal}
}

func (parser *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: parser.curToken}

	array.Elements = parser.parseExpressionList(token.RBRACKET)

	return array
}

func (parser *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: parser.curToken}

	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !parser.peekTokenIs(token.RBRACE) {
		parser.nextToken()

		key := parser.parseExpression(LOWEST)
		if !parser.expectPeek(token.COLON) {
			return nil
		}

		parser.nextToken()

		value := parser.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !parser.peekTokenIs(token.RBRACE) && !parser.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !parser.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (parser *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: parser.curToken, Left: left}

	parser.nextToken()
	expression.Index = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RBRACKET) {
		return nil
	}

	return expression
}

func (parser *Parser) parseIncrementExpression(left ast.Expression) ast.Expression {
	expression := &ast.IncrementExpression{Token: parser.curToken, Left: left}

	parser.nextToken()

	return expression
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
	list := []ast.Expression{}

	if parser.peekTokenIs(end) {
		parser.nextToken()
		return list
	}

	parser.nextToken()
	list = append(list, parser.parseExpression(LOWEST))

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()

		list = append(list, parser.parseExpression(LOWEST))
	}

	if parser.peekTokenIs(end) {
		parser.nextToken()
	}

	return list
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	message := fmt.Sprintf("\x1b[31m (line %d) expected next token to be: %s, got: %s (%s) instead.\x1b[0m",
		p.peekToken.LineNumber, t, p.peekToken.Type, p.peekToken.Literal)
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

	// Skip comment tokens
	if parser.curToken.Type == token.COMMENT {
		parser.curToken = parser.peekToken
		parser.peekToken = parser.lexer.NextToken()
	}
	if parser.peekToken.Type == token.COMMENT {
		parser.peekToken = parser.lexer.NextToken()
	}
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
	statement := &ast.LetStatement{Token: parser.curToken}
	mutable := false

	if parser.peekTokenIs(token.MUT) {
		mutable = true
		parser.nextToken()
	}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token:   parser.curToken,
		Value:   parser.curToken.Literal,
		Mutable: mutable,
	}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()

	statement.Value = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: parser.curToken}

	parser.nextToken()

	statement.ReturnValue = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
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
	message := fmt.Sprintf("\x1b[31m (line %d) no prefix parse function for %s found\x1b[0m", parser.curToken.LineNumber, t)
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

func (parser *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	// Check if the left side is a valid assignment target
	isValidTarget := false

	switch left.(type) {
	case *ast.Identifier, *ast.IndexExpression, *ast.PointerDereferenceExpression:
		isValidTarget = true
	}

	if !isValidTarget {
		parser.errors = append(parser.errors, fmt.Sprintf("(line %d) invalid assignment target: %s", parser.curToken.LineNumber, left.TokenLiteral()))
	}

	expression := &ast.AssignmentExpression{
		Token: parser.curToken,
		Left:  left,
	}

	precedence := parser.curPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}

func (parser *Parser) parsePointerReferenceExpression() ast.Expression {
	// Save the current token for the expression
	currentToken := parser.curToken

	// Move to the next token (after the &)
	parser.nextToken()

	// Create the expression with the correct token
	expression := &ast.PointerReferenceExpression{
		Token: currentToken,
		Right: parser.parseExpression(PREFIX),
	}

	return expression
}

func (parser *Parser) parsePointerDereferenceExpression() ast.Expression {
	// Save the current token for the expression
	currentToken := parser.curToken

	// Move to the next token (after the *)
	parser.nextToken()

	// Create the expression with the correct token
	expression := &ast.PointerDereferenceExpression{
		Token: currentToken,
		Right: parser.parseExpression(PREFIX),
	}

	return expression
}
