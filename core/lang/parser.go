package lang

import "fmt"

type Node any

type Identifier struct {
	Name string
}

type StringLiteral struct {
	Value string
}

type CallExpression struct {
	Name string
	Args []Node
}

type Parser struct {
	lexer *Lexer
	curr  Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer: lexer,
	}

	p.next()

	return p
}

func (p *Parser) Parse() (Node, error) {
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
	if p.curr.Type == TokenIdent {
		return p.parseIdentifierOrCall()
	}

	if p.curr.Type == TokenString {
		value := p.curr.Value

		if err := p.next(); err != nil {
			return nil, err
		}

		return StringLiteral{Value: value}, nil
	}

	return nil, p.error("Invalid expression")
}

func (p *Parser) parseIdentifierOrCall() (Node, error) {
	name := p.curr.Value

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

		if err := p.next(); err != nil {
			return nil, err
		}

		return CallExpression{Name: name, Args: args}, nil
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

	if err != nil {
		return err
	}

	p.curr = token

	return nil
}

func (p *Parser) error(cause string) error {
	return fmt.Errorf("%s\n%s", cause, p.lexer.buffer)
}
