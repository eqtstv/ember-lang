package logger

import (
	"ember_lang/ember_lang/ast"
	"ember_lang/ember_lang/lexer"
	"ember_lang/ember_lang/object"
	"fmt"
	"strings"
)

const (
	verticalLine   = "│"
	horizontalLine = "──"
	cornerLine     = "└──"
	teeLine        = "├──"

	// Atom One Dark theme colors
	purple = "\033[38;2;198;120;221m" // Keywords, control flow
	blue   = "\033[38;2;97;175;239m"  // Function keywords
	red    = "\033[38;2;224;108;117m" // Errors, important tokens
	cyan   = "\033[38;2;86;182;194m"  // Built-in types, numbers
	orange = "\033[38;2;209;154;102m" // Strings
	white  = "\033[38;2;171;178;191m" // Default text, identifiers
	gray   = "\033[38;2;92;99;112m"   // Tree structure, delimiters
	green  = "\033[38;2;152;195;121m" // Special keywords (true/false)
	reset  = "\033[0m"
)

func LogSourceCode(code []byte) {
	fmt.Printf("\n%s=== Source Code ===%s\n%s\n", blue, reset, string(code))
}

func LogTokens(l *lexer.Lexer) {
	fmt.Printf("\n%s=== Tokens ===%s\n", blue, reset)
	for tok := l.NextToken(); tok.Type != "EOF"; tok = l.NextToken() {
		fmt.Println(tok)
	}
}

func LogAST(node ast.Node) {
	fmt.Printf("\n%s=== Abstract Syntax Tree ===%s\n", blue, reset)
	printNode(node, "", true)
}

func LogResult(result object.Object) {
	fmt.Printf("\n%s=== Result ===%s\n", blue, reset)
}

func printNode(node ast.Node, prefix string, isLast bool) {
	if node == nil {
		return
	}

	connector := teeLine
	newPrefix := prefix + verticalLine + " "
	if isLast {
		connector = cornerLine
		newPrefix = prefix + "  "
	}

	nodeInfo := getNodeInfo(node)
	fmt.Printf("%s%s%s %s%s\n", gray, prefix, connector, nodeInfo, reset)

	switch n := node.(type) {
	case *ast.Program:
		printChildren(n.Statements, newPrefix)
	case *ast.LetStatement:
		printNode(n.Name, newPrefix, false)
		printNode(n.Value, newPrefix, true)
	case *ast.ReturnStatement:
		printNode(n.ReturnValue, newPrefix, true)
	case *ast.ExpressionStatement:
		printNode(n.Expression, newPrefix, true)
	case *ast.BlockStatement:
		printChildren(n.Statements, newPrefix)
	case *ast.IfExpression:
		printNode(n.Condition, newPrefix, false)
		printNode(n.Consequence, newPrefix, n.Alternative == nil)
		if n.Alternative != nil {
			printNode(n.Alternative, newPrefix, true)
		}
	case *ast.FunctionLiteral:
		for i, param := range n.Parameters {
			printNode(param, newPrefix, i == len(n.Parameters)-1 && n.Body == nil)
		}
		if n.Body != nil {
			printNode(n.Body, newPrefix, true)
		}
	case *ast.CallExpression:
		printNode(n.Function, newPrefix, len(n.Arguments) == 0)
		for i, arg := range n.Arguments {
			printNode(arg, newPrefix, i == len(n.Arguments)-1)
		}
	case *ast.InfixExpression:
		printNode(n.Left, newPrefix, false)
		printNode(n.Right, newPrefix, true)
	case *ast.PrefixExpression:
		printNode(n.Right, newPrefix, true)
	case *ast.ArrayLiteral:
		for i, elem := range n.Elements {
			printNode(elem, newPrefix, i == len(n.Elements)-1)
		}
	case *ast.IndexExpression:
		printNode(n.Left, newPrefix, false)
		printNode(n.Index, newPrefix, true)
	}
}

func printChildren(nodes []ast.Statement, prefix string) {
	for i, node := range nodes {
		printNode(node, prefix, i == len(nodes)-1)
	}
}

func getNodeInfo(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Program:
		return purple + "Program"
	case *ast.LetStatement:
		return purple + "Let Statement"
	case *ast.ReturnStatement:
		return purple + "Return Statement"
	case *ast.ExpressionStatement:
		return white + "Expression Statement"
	case *ast.BlockStatement:
		return white + "Block Statement"
	case *ast.IfExpression:
		return purple + "If Expression"
	case *ast.Identifier:
		return white + "Identifier: " + cyan + n.Value
	case *ast.IntegerLiteral:
		return cyan + fmt.Sprintf("Integer: %d", n.Value)
	case *ast.Boolean:
		return green + fmt.Sprintf("Boolean: %t", n.Value)
	case *ast.FunctionLiteral:
		params := make([]string, len(n.Parameters))
		for i, p := range n.Parameters {
			params[i] = p.String()
		}
		return blue + fmt.Sprintf("Function: fn(%s)", strings.Join(params, ", "))
	case *ast.CallExpression:
		return blue + "Call Expression"
	case *ast.InfixExpression:
		return white + fmt.Sprintf("Infix: %s%s", orange, n.Operator)
	case *ast.PrefixExpression:
		return white + fmt.Sprintf("Prefix: %s%s", orange, n.Operator)
	case *ast.StringLiteral:
		return orange + fmt.Sprintf("String: %q", n.Value)
	case *ast.ArrayLiteral:
		return cyan + "Array"
	case *ast.IndexExpression:
		return white + "Index Expression"
	default:
		return red + fmt.Sprintf("%T", node)
	}
}
