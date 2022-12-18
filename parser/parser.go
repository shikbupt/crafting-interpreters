package parser

import (
	"craftinginterpreters/lox/scanner"
	"errors"
	"fmt"
)

type Parser struct {
	Tokens   []scanner.Token
	current  int
	HadError bool
}

func NewParser(t []scanner.Token) *Parser {
	return &Parser{
		Tokens: t,
	}
}

func (p *Parser) Parse() (res Expr, err error) {
	defer func() {
		if terr := recover(); terr != nil {
			res = nil
			err = terr.(error)
		}
	}()
	res = p.expression()
	return res, err
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &Unary{operator, right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(scanner.FALSE):
		return &Literal{false}
	case p.match(scanner.TRUE):
		return &Literal{true}
	case p.match(scanner.NIL):
		return &Literal{nil}
	case p.match(scanner.NUMBER, scanner.STRING):
		return &Literal{p.previous().Literal}
	case p.match(scanner.LEFT_PAREN):
		expr := p.expression()
		p.comsume(scanner.RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{expr}
	}
	p.Error(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) match(tokens ...scanner.TokenType) bool {
	for _, token := range tokens {
		if p.check(token) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) comsume(tokentype scanner.TokenType, message string) scanner.Token {
	if p.check(tokentype) {
		return p.advance()
	}
	p.Error(p.peek(), message)
	return scanner.Token{}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR, scanner.IF:
			return
		case scanner.WHILE, scanner.PRINT, scanner.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) Error(token scanner.Token, message string) {
	if token.Type == scanner.EOF {
		p.report(token.Line, " at end", message)
	}
	p.report(token.Line, "at '"+token.Lexeme+"'", message)
}

func (p *Parser) report(line int, where, message string) {
	p.HadError = true
	t := fmt.Sprintf("[line %d ] Error %s: %s", line, where, message)
	panic(errors.New(t))
}
