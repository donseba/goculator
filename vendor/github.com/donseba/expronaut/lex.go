package expronaut

import "fmt"

type TokenType string

const (
	TokenTypeEOF TokenType = "EOF"

	TokenTypeInt       TokenType = "INT"
	TokenTypeFloat     TokenType = "FLOAT"
	TokenTypeParenLeft TokenType = "PAREN_LEFT"

	TokenTypeParenRight TokenType = "PAREN_RIGHT"
	TokenTypeIllegal    TokenType = "ILLEGAL"
	TokenTypeString     TokenType = "STRING"
	TokenTypeVariable   TokenType = "VARIABLE"
	TokenTypeFunction   TokenType = "FUNCTION"
	TokenTypeArray      TokenType = "ARRAY"
	TokenTypeArrayStart TokenType = "ARRAY_START"
	TokenTypeArrayEnd   TokenType = "ARRAY_END"
	TokenTypeComma      TokenType = "COMMA"

	TokenTypeAnd                TokenType = "AND"
	TokenTypeOr                 TokenType = "OR"
	TokenTypePlus               TokenType = "PLUS"
	TokenTypeMinus              TokenType = "MINUS"
	TokenTypeMultiply           TokenType = "MULTIPLY"
	TokenTypeDivide             TokenType = "DIVIDE"
	TokenTypeDivideInteger      TokenType = "DIVIDE_INTEGER"
	TokenTypeModulo             TokenType = "MODULO"
	TokenTypeEqual              TokenType = "EQUAL"
	TokenTypeNotEqual           TokenType = "NOT_EQUAL"
	TokenTypeLessThan           TokenType = "LESS_THAN"
	TokenTypeGreaterThan        TokenType = "GREATER_THAN"
	TokenTypeLessThanOrEqual    TokenType = "LESS_THAN_OR_EQUAL"
	TokenTypeGreaterThanOrEqual TokenType = "GREATER_THAN_OR_EQUAL"
	TokenTypeBool               TokenType = "BOOL"
	TokenTypeExponent           TokenType = "EXPONENT"
	TokenTypeLeftShift          TokenType = "LEFT_SHIFT"
	TokenTypeRightShift         TokenType = "RIGHT_SHIFT"
)

func TokenGoTemplate(tok TokenType) string {
	switch tok {
	case TokenTypeAnd:
		return "and"
	case TokenTypeOr:
		return "or"
	case TokenTypeEqual:
		return "eq"
	case TokenTypeNotEqual:
		return "ne"
	case TokenTypeLessThan:
		return "lt"
	case TokenTypeGreaterThan:
		return "gt"
	case TokenTypeLessThanOrEqual:
		return "le"
	case TokenTypeGreaterThanOrEqual:
		return "ge"
	case TokenTypePlus:
		return "add"
	case TokenTypeMinus:
		return "sub"
	case TokenTypeMultiply:
		return "mul"
	case TokenTypeDivide:
		return "div"
	case TokenTypeModulo:
		return "mod"
	default:
		return string(tok)
	}
}

type Token struct {
	Type    TokenType // The type of token, indicating its role (e.g., operator, number, parenthesis)
	Literal string    // The actual text that the token represents (e.g., "123", "+", "(")
}

type Lexer struct {
	input        string
	position     int  // Current position in input (points to current char)
	readPosition int  // Current reading position in input (after current char)
	ch           byte // Current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Initialize the first character
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" character signifies end of input
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = newToken(TokenTypeParenLeft, l.ch)
	case ')':
		tok = newToken(TokenTypeParenRight, l.ch)
	case '+':
		tok = newToken(TokenTypePlus, l.ch)
	case '-':
		tok = newToken(TokenTypeMinus, l.ch)

		if isDigit(l.peekChar()) {
			return l.readNumber(true)
		} else {
			tok = newToken(TokenTypeMinus, l.ch)
		}
	case '*':
		if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeExponent, Literal: literal}
		} else {
			tok = newToken(TokenTypeMultiply, l.ch)
		}
	case '^':
		tok = newToken(TokenTypeExponent, l.ch)
	case '/':
		if l.peekChar() == '/' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeDivideInteger, Literal: literal}
		} else {
			tok = newToken(TokenTypeDivide, l.ch)
		}
	case '%':
		tok = newToken(TokenTypeModulo, l.ch)
	case ',':
		tok = newToken(TokenTypeComma, l.ch)
	case '[':
		tok = newToken(TokenTypeArrayStart, l.ch)
	case ']':
		tok = newToken(TokenTypeArrayEnd, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeLessThanOrEqual, Literal: literal}
		} else if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeLeftShift, Literal: literal}
		} else {
			tok = newToken(TokenTypeLessThan, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeGreaterThanOrEqual, Literal: literal}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeRightShift, Literal: literal}
		} else {
			tok = newToken(TokenTypeGreaterThan, l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeEqual, Literal: literal}
		} else {
			tok = newToken(TokenTypeIllegal, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeNotEqual, Literal: literal}
		} else {
			tok = newToken(TokenTypeIllegal, l.ch)
		}
	case '"', '\'':
		tok.Literal = l.readString()
		tok.Type = TokenTypeString
	case '`':
		tok.Literal = l.readUntilBacktick()
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeAnd, Literal: literal}
		} else {
			tok = newToken(TokenTypeIllegal, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TokenTypeOr, Literal: literal}
		} else {
			tok = newToken(TokenTypeIllegal, l.ch)
		}
	case 0:
		tok.Type = TokenTypeEOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			identifier, iType := l.readIdentifier()

			switch iType {
			case identifierTypeFunction:
				tok.Type = TokenTypeFunction
				tok.Literal = identifier
				// The '(' character will be processed in the next iteration of NextToken
				return tok
			case identifierTypeArray:
				tok.Type = TokenTypeArray
				tok.Literal = identifier
				// The '[' character will be processed in the next iteration of NextToken
				return tok
			case identifierTypeVariable:
				tok.Type = TokenTypeVariable
				tok.Literal = identifier
				return tok
			case identifierTypeBool:
				tok.Type = TokenTypeBool
				tok.Literal = identifier
			default:
				fmt.Println("Unknown identifier type")
			}
		} else if isDigit(l.ch) {
			return l.readNumber(false)
		} else {
			tok = newToken(TokenTypeIllegal, l.ch)
		}
	}

	l.readChar()
	return tok
}

// newToken is a helper function to create a new Token.
func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

// peekChar looks at the next character without moving the current position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func (l *Lexer) readNumber(isNegative bool) Token {
	position := l.position

	hasDecimal := false
	for isDigit(l.ch) || (!hasDecimal && l.ch == '.') || (isNegative && l.ch == '-') {
		if l.ch == '.' {
			hasDecimal = true
		}
		l.readChar()
	}

	if hasDecimal {
		return Token{Type: TokenTypeFloat, Literal: l.input[position:l.position]}
	}

	return Token{Type: TokenTypeInt, Literal: l.input[position:l.position]}
}

type identifierType string

const (
	identifierTypeVariable identifierType = "variable"
	identifierTypeFunction identifierType = "function"
	identifierTypeArray    identifierType = "array"
	identifierTypeBool     identifierType = "bool"
)

func (l *Lexer) readIdentifier() (string, identifierType) {
	startPosition := l.position
	returnType := identifierTypeVariable

	for isLetter(l.peekChar()) || isDigit(l.peekChar()) || l.peekChar() == '_' || l.peekChar() == '.' {
		l.readChar()
	}

	if l.peekChar() == '[' {
		returnType = identifierTypeArray
	} else if l.peekChar() == '(' {
		returnType = identifierTypeFunction
	}

	l.readChar()

	ident := l.input[startPosition:l.position]
	if ident == "true" || ident == "false" {
		return ident, identifierTypeBool
	}

	return ident, returnType
}

// Helper function to read string literals surrounded by double quotes
func (l *Lexer) readString() string {
	position := l.position + 1 // Start after the initial double quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readUntilBacktick() string {
	position := l.position + 1 // Start after the initial backtick
	for {
		if l.peekChar() == '`' {
			l.readChar()
			break
		}
		l.readChar()
	}
	l.readChar()
	return l.input[position:l.position]
}
