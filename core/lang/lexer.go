package lang

import "fmt"

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

			for range l.expr.Open {
				l.buffer.Advance()
			}

			l.mode = ModeExpr
			l.exprFilled = false

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
		if !l.exprFilled {
			return Token{}, l.error("Illegal empty expression")
		}

		for range l.expr.Close {
			l.buffer.Advance()
		}

		l.mode = ModeRaw

		return Token{Type: TokenCloseExpr, Value: l.expr.Close}, nil
	}

	ch := l.buffer.Peek()

	if ch == '(' {
		l.buffer.Advance()
		l.exprFilled = true
		return Token{Type: TokenOpenParen, Value: "("}, nil
	}

	if ch == ')' {
		l.buffer.Advance()
		l.exprFilled = true
		return Token{Type: TokenCloseParen, Value: ")"}, nil
	}

	if ch == ',' {
		l.buffer.Advance()
		l.exprFilled = true
		return Token{Type: TokenComma, Value: ","}, nil
	}

	if ch == '"' || ch == '\'' {
		start := l.buffer.Index

		err := l.advanceString(ch)

		if err != nil {
			return Token{}, err
		}

		l.exprFilled = true

		return Token{Type: TokenString, Value: l.buffer.SliceFrom(start)}, nil
	}

	if l.isLetter(ch) || ch == '_' {
		start := l.buffer.Index

		l.buffer.Advance()

		for {
			ch := l.buffer.Peek()

			if !(l.isLetter(ch) || l.isDigit(ch) || ch == '_') {
				break
			}

			l.buffer.Advance()
		}

		l.exprFilled = true

		return Token{Type: TokenIdent, Value: l.buffer.SliceFrom(start)}, nil
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

func (l *Lexer) isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (l *Lexer) isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == '\t'
}

func (l *Lexer) error(cause string) error {
	return fmt.Errorf("%s\n%s", cause, l.buffer)
}
