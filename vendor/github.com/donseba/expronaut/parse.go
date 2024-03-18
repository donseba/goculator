package expronaut

import (
	"strconv"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(lexer *Lexer) *Parser {
	var tokens []Token
	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenTypeEOF {
			break
		}
	}

	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// expression parses an expression.
func (p *Parser) expression() ASTNode {
	return p.logicalOr()
}

// logicalOr handles ||.
func (p *Parser) logicalOr() ASTNode {
	node := p.logicalAnd()

	for p.match(TokenTypeOr) {
		operator := p.previous()
		right := p.logicalAnd()
		node = &LogicalOperationNode{Left: node, Operator: operator.Type, Right: right}
	}

	return node
}

// logicalAnd handles &&.
func (p *Parser) logicalAnd() ASTNode {
	node := p.equality()

	for p.match(TokenTypeAnd) {
		operator := p.previous()
		right := p.equality()
		node = &LogicalOperationNode{Left: node, Operator: operator.Type, Right: right}
	}

	return node
}

// equality handles == and !=.
func (p *Parser) equality() ASTNode {
	node := p.comparison()

	for p.match(TokenTypeEqual, TokenTypeNotEqual) {
		operator := p.previous()
		right := p.comparison()
		node = &BinaryOperationNode{Left: node, Operator: operator.Type, Right: right}
	}

	return node
}

// comparison handles <, <=, >, and >=.
func (p *Parser) comparison() ASTNode {
	node := p.shift()

	for p.match(TokenTypeGreaterThan, TokenTypeGreaterThanOrEqual, TokenTypeLessThan, TokenTypeLessThanOrEqual) {
		operator := p.previous()
		right := p.shift()
		node = &BinaryOperationNode{Left: node, Operator: operator.Type, Right: right}
	}

	return node
}

// shift handles << and >>.
func (p *Parser) shift() ASTNode {
	node := p.addition()

	for p.match(TokenTypeLeftShift, TokenTypeRightShift) {
		operator := p.previous()
		right := p.addition()
		node = &BinaryOperationNode{Left: node, Operator: operator.Type, Right: right}
	}

	return node
}

// addition handles + and -.
func (p *Parser) addition() ASTNode {
	node := p.multiplication()

	for p.match(TokenTypePlus, TokenTypeMinus) {
		operator := p.previous()
		right := p.multiplication()
		node = &BinaryOperationNode{Left: node, Operator: operator.Type, Right: right}
	}

	return node
}

// multiplication handles *, /
func (p *Parser) multiplication() ASTNode {
	node := p.functions()

	for p.match(TokenTypeMultiply, TokenTypeDivide, TokenTypeModulo, TokenTypeDivideInteger) {
		operator := p.previous()
		right := p.functions()
		node = &BinaryOperationNode{Left: node, Operator: operator.Type, Right: right}
	}

	return node
}

func (p *Parser) functions() ASTNode {
	node := p.primary()

	for p.match(TokenTypeExponent, TokenTypeFunction) {
		operator := p.previous()
		right := p.functions()
		node = &BinaryOperationNode{Left: node, Operator: operator.Type, Right: right}
	}

	return node
}

// primary handles the base case of the recursive descent parser.
func (p *Parser) primary() ASTNode {
	switch {
	case p.match(TokenTypeInt):
		return &IntLiteralNode{Value: parseInt(p.previous().Literal)}
	case p.match(TokenTypeFloat):
		return &FloatLiteralNode{Value: parseFloat(p.previous().Literal)}
	case p.match(TokenTypeString):
		return &StringLiteralNode{Value: p.previous().Literal}
	case p.match(TokenTypeBool):
		return &BooleanLiteralNode{Value: parseBool(p.previous().Literal)}
	case p.match(TokenTypeVariable):
		return &VariableNode{Name: p.previous().Literal}
	case p.match(TokenTypeParenLeft):
		expr := p.expression()
		p.consume(TokenTypeParenRight, "Expect ')' after expression.")
		return expr
	case p.match(TokenTypeFunction):
		funcName := p.previous().Literal
		var arguments []ASTNode

		p.consume(TokenTypeParenLeft, "Expect '(' after function.")

		if !p.check(TokenTypeParenRight) {
			for {
				// Parse an argument expression.
				argument := p.expression()
				arguments = append(arguments, argument)

				// If there's no comma, stop parsing arguments.
				if !p.match(TokenTypeComma) {
					break
				}
			}
		}

		// After parsing all arguments, expect the closing parenthesis.
		p.consume(TokenTypeParenRight, "Expect ')' after arguments to function.")

		return &FunctionCallNode{FunctionName: funcName, Arguments: arguments}
	case p.match(TokenTypeArray):
		var elements []ASTNode
		arrayType := arrayType(p.previous().Literal)

		p.consume(TokenTypeArrayStart, "Expect '[' after array.")
		if !p.check(TokenTypeArrayEnd) {
			for {
				// Parse an element expression.
				element := p.expression()
				elements = append(elements, element)

				// If there's no comma, stop parsing elements.
				if !p.match(TokenTypeComma) {
					break
				}
			}
		}

		// After parsing all elements, expect the closing bracket.
		p.consume(TokenTypeArrayEnd, "Expect ']' after elements to array.")

		return &ArrayNode{Type: arrayType, Elements: elements}
	}

	return &IntLiteralNode{Value: 0}
}

// previous returns the previous token.
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// consume expects the next token to be of a given type and consumes it, or throws an error.
func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic(message) // Or handle the error more gracefully
}

// check looks at the current token and returns true if it matches the given type.
func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == typ
}

// isAtEnd checks if we've consumed all tokens.
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TokenTypeEOF
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// advance consumes the current token and returns it.
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.tokens[p.current-1]
}

// match checks if the current token matches any of the given types.
func (p *Parser) match(types ...TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

// Parse starts the parsing process.
func (p *Parser) Parse() ASTNode {
	return p.expression() // Start parsing from the highest level of precedence.
}

// parseNumber converts a string literal to a float64.
func parseFloat(lit string) float64 {
	value, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		panic(err) // or handle the error as appropriate
	}
	return value
}

func parseInt(lit string) int {
	value, err := strconv.Atoi(lit)
	if err != nil {
		panic(err) // or handle the error as appropriate
	}
	return value
}

func parseBool(lit string) bool {
	return lit == "true" // Returns true for "true", false otherwise
}
