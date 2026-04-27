package lang

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/util"
)

type TokenType int
type ModeType int

const (
	TokenRaw TokenType = iota
	TokenIdent
	TokenString
	TokenOpenParen
	TokenCloseParen
	TokenComma
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
	case TokenIdent:
		return "TokenIdent"
	case TokenString:
		return "TokenString"
	case TokenOpenParen:
		return "TokenOpenParen"
	case TokenCloseParen:
		return "TokenCloseParen"
	case TokenComma:
		return "TokenComma"
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
	mode       ModeType
	buffer     *Buffer
	expr       *Expression
	exprFilled bool
}

func NewLexer(buffer *Buffer, expr *Expression) *Lexer {
	return &Lexer{
		mode:       ModeRaw,
		buffer:     buffer,
		expr:       expr,
		exprFilled: false,
	}
}

func (l *Lexer) abort(cause string) {
	cli.LogAndExit(fmt.Sprintf("%s\n%s", cause, l.buffer), cli.LogSeverityError)
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

func (l *Lexer) advanceOpenExpr() {
	for range l.expr.Open {
		l.buffer.Advance()
	}
}

func (l *Lexer) advanceCloseExpr() {
	for range l.expr.Close {
		l.buffer.Advance()
	}
}

func (l *Lexer) Next() Token {
	// EOF
	if l.buffer.Index >= len(l.buffer.Data) {
		return Token{
			Type: TokenEOF,
		}
	}

	if l.mode == ModeRaw {
		return l.readRaw()
	}

	if l.mode == ModeExpr {
		return l.readExpr()
	}

	// EOF
	return Token{
		Type: TokenEOF,
	}
}

func (l *Lexer) readRaw() Token {
	start := l.buffer.Index

	for {
		if l.buffer.Index >= len(l.buffer.Data) {
			break
		}

		if l.match(l.expr.Open) {
			if l.buffer.Index > start {
				return Token{
					Type:  TokenRaw,
					Value: string(l.buffer.Data[start:l.buffer.Index]),
				}
			}

			l.advanceOpenExpr()
			l.mode = ModeExpr
			l.exprFilled = false

			return Token{
				Type:  TokenOpenExpr,
				Value: l.expr.Open,
			}
		}

		l.buffer.Advance()
	}

	return Token{
		Type:  TokenRaw,
		Value: string(l.buffer.Data[start:l.buffer.Index]),
	}
}

func (l *Lexer) readExpr() Token {
	if l.buffer.Peek() == 0 {
		l.abort("Illegal unterminated expression")
	}

	if l.match(l.expr.Open) {
		l.abort("Illegal nested expression")
	}

	if l.match(l.expr.Close) {
		if !l.exprFilled {
			l.abort("Illegal empty expression")
		}

		l.advanceCloseExpr()
		l.mode = ModeRaw

		return Token{
			Type:  TokenCloseExpr,
			Value: l.expr.Close,
		}
	}

	ch := l.buffer.Peek()

	if ch == '(' {
		l.buffer.Advance()
		l.exprFilled = true
		return Token{
			Type:  TokenOpenParen,
			Value: "(",
		}
	}

	if ch == ')' {
		l.buffer.Advance()
		l.exprFilled = true
		return Token{
			Type:  TokenCloseParen,
			Value: ")",
		}
	}

	if ch == ',' {
		l.buffer.Advance()
		l.exprFilled = true
		return Token{
			Type:  TokenComma,
			Value: ",",
		}
	}

	if ch == '"' || ch == '\'' {
		l.exprFilled = true
		return Token{
			Type:  TokenString,
			Value: l.readString(ch),
		}
	}

	if util.IsLetter(ch) {
		start := l.buffer.Index

		for {
			ch := l.buffer.Peek()

			if !(util.IsLetter(ch) || util.IsDigit(ch) || ch == '_') {
				break
			}

			l.buffer.Advance()
		}

		l.exprFilled = true

		return Token{
			Type:  TokenIdent,
			Value: string(l.buffer.Data[start:l.buffer.Index]),
		}
	}

	if util.IsWhitespace(ch) {
		l.buffer.Advance()

		return l.Next()
	}

	l.abort("Illegal character in expression")

	// unreachable
	return Token{}
}

func (l *Lexer) readString(quote rune) string {
	// skip opening quote
	l.buffer.Advance()

	start := l.buffer.Index

	for {
		ch := l.buffer.Peek()

		if ch == 0 {
			l.abort("Illegal unterminated string")
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

	value := string(l.buffer.Data[start:l.buffer.Index])

	// skip closing quote
	if l.buffer.Peek() == quote {
		l.buffer.Advance()
	}

	return value
}
