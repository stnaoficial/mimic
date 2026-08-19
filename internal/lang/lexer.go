package lang

import (
	"fmt"
)

type TokenType int
type ModeType int

const (
	TokenRaw TokenType = iota
	TokenString
	TokenNumber
	TokenIdentifier
	TokenOpenParen
	TokenCloseParen
	TokenComma
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenOpenExpr
	TokenCloseExpr
	TokenEOF
)

const (
	ModeRaw ModeType = iota
	ModeExpr
)

func (t TokenType) String() string {
	switch t {
	case TokenRaw:
		return "TokenRaw"
	case TokenString:
		return "TokenString"
	case TokenNumber:
		return "TokenNumber"
	case TokenIdentifier:
		return "TokenIdentifier"
	case TokenOpenParen:
		return "TokenOpenParen"
	case TokenCloseParen:
		return "TokenCloseParen"
	case TokenComma:
		return "TokenComma"
	case TokenPlus:
		return "TokenPlus"
	case TokenMinus:
		return "TokenMinus"
	case TokenMultiply:
		return "TokenMultiply"
	case TokenDivide:
		return "TokenDivide"
	case TokenOpenExpr:
		return "TokenOpenExpr"
	case TokenCloseExpr:
		return "TokenCloseExpr"
	case TokenEOF:
		return "TokenEOF"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	mode   ModeType
	buffer *Buffer
	expr   *Expression
}

func NewLexer(buffer *Buffer, expr *Expression) *Lexer {
	return &Lexer{
		mode:   ModeRaw,
		buffer: buffer,
		expr:   expr,
	}
}

func (l *Lexer) Next() (Token, error) {
	// EOF
	if l.buffer.Index >= len(l.buffer.Data) {
		return Token{Type: TokenEOF}, nil
	}

	if l.mode == ModeRaw {
		return l.readRaw()
	}

	if l.mode == ModeExpr {
		return l.readExpr()
	}

	// EOF
	return Token{Type: TokenEOF}, nil
}

func (l *Lexer) readRaw() (Token, error) {
	start := l.buffer.Index

	for {
		if l.buffer.Index >= len(l.buffer.Data) {
			break
		}

		if l.match(l.expr.Open) {
			if l.buffer.Index > start {
				return Token{Type: TokenRaw, Value: l.buffer.SliceFrom(start)}, nil
			}

			// skip open expression syntax
			for range l.expr.Open {
				l.buffer.Advance()
			}

			l.mode = ModeExpr

			return Token{Type: TokenOpenExpr, Value: l.expr.Open}, nil
		}

		l.buffer.Advance()
	}

	return Token{Type: TokenRaw, Value: l.buffer.SliceFrom(start)}, nil
}

func (l *Lexer) readExpr() (Token, error) {
	if l.buffer.Peek() == 0 {
		return Token{}, l.error("Illegal unterminated expression")
	}

	if l.match(l.expr.Open) {
		return Token{}, l.error("Illegal nested expression")
	}

	if l.match(l.expr.Close) {
		// skip close expression syntax
		for range l.expr.Close {
			l.buffer.Advance()
		}

		l.mode = ModeRaw

		return Token{Type: TokenCloseExpr, Value: l.expr.Close}, nil
	}

	ch := l.buffer.Peek()

	if ch == '(' {
		l.buffer.Advance()
		return Token{Type: TokenOpenParen, Value: "("}, nil
	}

	if ch == ')' {
		l.buffer.Advance()
		return Token{Type: TokenCloseParen, Value: ")"}, nil
	}

	if ch == ',' {
		l.buffer.Advance()
		return Token{Type: TokenComma, Value: ","}, nil
	}

	if ch == '+' {
		l.buffer.Advance()
		return Token{Type: TokenPlus, Value: "+"}, nil
	}

	if ch == '-' {
		l.buffer.Advance()
		return Token{Type: TokenMinus, Value: "-"}, nil
	}

	if ch == '*' {
		l.buffer.Advance()
		return Token{Type: TokenMultiply, Value: "*"}, nil
	}

	if ch == '/' {
		l.buffer.Advance()
		return Token{Type: TokenDivide, Value: "/"}, nil
	}

	if ch == '"' || ch == '\'' {
		start := l.buffer.Index

		if err := l.advanceString(ch); err != nil {
			return Token{}, err
		}

		return Token{Type: TokenString, Value: l.buffer.SliceFrom(start)}, nil
	}

	if l.isDigit(ch) {
		start := l.buffer.Index

		if err := l.advanceNumber(); err != nil {
			return Token{}, err
		}

		return Token{Type: TokenNumber, Value: l.buffer.SliceFrom(start)}, nil
	}

	if l.isLetter(ch) || ch == '_' {
		start := l.buffer.Index

		if err := l.advanceIdentifier(); err != nil {
			return Token{}, err
		}

		return Token{Type: TokenIdentifier, Value: l.buffer.SliceFrom(start)}, nil
	}

	if l.isWhitespace(ch) {
		l.buffer.Advance()
		return l.Next()
	}

	return Token{}, l.error("Illegal character in expression")
}

func (l *Lexer) advanceString(quote rune) error {
	// skip opening quote
	l.buffer.Advance()

	for {
		ch := l.buffer.Peek()

		if ch == 0 {
			return l.error("Illegal unterminated string")
		}

		if ch == '\\' {
			// skip escape + next char
			l.buffer.Advance()
			l.buffer.Advance()
			continue
		}

		if ch == quote {
			break
		}

		l.buffer.Advance()
	}

	// skip closing quote
	if l.buffer.Peek() == quote {
		l.buffer.Advance()
	}

	return nil
}

func (l *Lexer) advanceNumber() error {
	// skip first digit
	l.buffer.Advance()

	for l.isNumberPart(l.buffer.Peek()) {
		l.buffer.Advance()
	}

	return nil
}

func (l *Lexer) advanceIdentifier() error {
	// skip first letter
	l.buffer.Advance()

	for l.isIdentifierPart(l.buffer.Peek()) {
		l.buffer.Advance()
	}

	return nil
}

func (l *Lexer) match(str string) bool {
	if l.buffer.Index+len(str) > len(l.buffer.Data) {
		return false
	}

	for i, r := range str {
		if l.buffer.Data[l.buffer.Index+i] != r {
			return false
		}
	}

	return true
}

func (l *Lexer) isLowerCaseLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z')
}

func (l *Lexer) isUpperCaseLetter(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z')
}

func (l *Lexer) isLetter(ch rune) bool {
	return l.isLowerCaseLetter(ch) || l.isUpperCaseLetter(ch)
}

func (l *Lexer) isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == '\t'
}

func (l *Lexer) isNumberPart(ch rune) bool {
	return l.isDigit(ch) || ch == '.'
}

func (l *Lexer) isIdentifierPart(ch rune) bool {
	return l.isLetter(ch) || l.isDigit(ch) || ch == '_'
}

func (l *Lexer) error(cause string) error {
	return fmt.Errorf("%s\n%s", cause, l.buffer)
}
