package lang

import (
	"fmt"
	"strconv"
)

type Node any

type Null struct{}

type String struct {
	Value string
}

type Identifier struct {
	Name string
}

type Callable struct {
	Name string
	Args []Node
}

type Parser struct {
	lexer *Lexer
	curr  Token
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func (p *Parser) Parse() (Node, error) {
	p.next()

	if p.curr.Type == TokenCloseExpr {
		return Null{}, nil
	}

	node, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	if p.curr.Type != TokenEOF && p.curr.Type != TokenCloseExpr {
		return nil, p.error("Invalid token")
	}

	return node, nil
}

func (p *Parser) parseExpression() (Node, error) {
	if p.curr.Type == TokenString {
		value := p.curr.Value

		// next token
		if err := p.next(); err != nil {
			return nil, err
		}

		str, err := strconv.Unquote(value)

		if err != nil {
			return nil, p.error("Invalid string syntax")
		}

		return String{Value: str}, nil
	}

	if p.curr.Type == TokenIdentifier {
		return p.parseIdentifierOrCall()
	}

	return nil, p.error("Invalid expression")
}

func (p *Parser) parseIdentifierOrCall() (Node, error) {
	name := p.curr.Value

	// next token
	if err := p.next(); err != nil {
		return nil, err
	}

	// function call
	if p.curr.Type == TokenOpenParen {
		if err := p.next(); err != nil {
			return nil, err
		}

		args, err := p.parseArguments()

		if err != nil {
			return nil, err
		}

		if p.curr.Type != TokenCloseParen {
			return nil, p.error("Invalid call signature")
		}

		// next token
		if err := p.next(); err != nil {
			return nil, err
		}

		return Callable{Name: name, Args: args}, nil
	}

	// plain identifier
	return Identifier{Name: name}, nil
}

func (p *Parser) parseArguments() ([]Node, error) {
	var args []Node

	if p.curr.Type == TokenCloseParen {
		return args, nil
	}

	for {
		arg, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		if arg == nil {
			return nil, p.error("Invalid argument")
		}

		args = append(args, arg)

		if p.curr.Type == TokenComma {
			// next token
			if err := p.next(); err != nil {
				return nil, err
			}

			if p.curr.Type == TokenCloseParen {
				return nil, p.error("Invalid trailing comma")
			}

			continue
		}

		break
	}

	return args, nil
}

func (p *Parser) next() error {
	token, err := p.lexer.Next()
	p.curr = token
	return err
}

func (p *Parser) error(cause string) error {
	return fmt.Errorf("%s\n%s", cause, p.lexer.buffer)
}
