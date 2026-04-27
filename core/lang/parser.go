package lang

import (
	"fmt"
	"mimic/core/cli"
)

type Node interface{}

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

func (p *Parser) abort(cause string) {
	cli.LogAndExit(fmt.Sprintf("%s\n%s", cause, p.lexer.buffer), cli.LogSeverityError)
}

func (p *Parser) next() {
	p.curr = p.lexer.Next()
}

func (p *Parser) Parse() Node {
	node := p.parseExpression()

	if p.curr.Type != TokenEOF && p.curr.Type != TokenCloseExpr {
		p.abort("Invalid token")
	}

	return node
}

func (p *Parser) parseExpression() Node {
	if p.curr.Type == TokenIdent {
		return p.parseIdentifierOrCall()
	}

	if p.curr.Type == TokenString {
		value := p.curr.Value

		p.next()

		return StringLiteral{
			Value: value,
		}
	}

	p.abort("Invalid expression")

	// unreachable
	return nil
}

func (p *Parser) parseIdentifierOrCall() Node {
	name := p.curr.Value

	p.next()

	// function call
	if p.curr.Type == TokenOpenParen {
		p.next()

		args := p.parseArguments()

		if p.curr.Type != TokenCloseParen {
			p.abort("Invalid call signature")
		}

		p.next()

		return CallExpression{
			Name: name,
			Args: args,
		}
	}

	// plain identifier
	return Identifier{
		Name: name,
	}
}

func (p *Parser) parseArguments() []Node {
	var args []Node

	if p.curr.Type == TokenCloseParen {
		return args
	}

	for {
		arg := p.parseExpression()

		if arg == nil {
			p.abort("Invalid argument")
		}

		args = append(args, arg)

		if p.curr.Type == TokenComma {
			p.next()

			if p.curr.Type == TokenCloseParen {
				p.abort("Invalid trailing comma")
			}

			continue
		}

		break
	}

	return args
}
