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

func (p *Parser) Parse() (res []Stmt, err error) {
	defer func() {
		if terr := recover(); terr != nil {
			res = nil
			err = terr.(error)
		}
	}()
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.Declaration())
	}
	return statements, nil
}

func (p *Parser) Declaration() Stmt {
	if p.match(scanner.VAR) {
		return p.VarDeclaration()
	}
	return p.Statement()
}

func (p *Parser) VarDeclaration() Stmt {
	name := p.comsume(scanner.IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(scanner.EQUAL) {
		initializer = p.Expression()
	}
	p.comsume(scanner.SEMICOLON, "Expect ';' after variable declaration.")
	return &Var{
		Name:        name,
		Initializer: initializer,
	}
}

func (p *Parser) Expression() Expr {
	return p.Assignment()
}

func (p *Parser) Assignment() Expr {
	expr := p.Equality()
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value := p.Assignment()
		if v, ok := expr.(*Variable); ok {
			name := v.Name
			return &Assign{
				Name:  name,
				Value: value,
			}
		}
		panic(fmt.Sprintf("%v Invalid assignment target.", equals))
	}
	return expr
}

func (p *Parser) Statement() Stmt {
	if p.match(scanner.PRINT) {
		return p.PrintStatement()
	}
	if p.match(scanner.LEFT_BRACE) {
		return &Block{p.Block()}
	}
	return p.ExpressionStatement()
}

func (p *Parser) Block() []Stmt {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() && !p.check(scanner.RIGHT_BRACE) {
		statements = append(statements, p.Declaration())
	}
	p.comsume(scanner.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) PrintStatement() Stmt {
	value := p.Expression()
	p.comsume(scanner.SEMICOLON, "Expect ';' after value.")
	return &Print{value}
}

func (p *Parser) ExpressionStatement() Stmt {
	value := p.Expression()
	p.comsume(scanner.SEMICOLON, "Expect ';' after expression.")
	return &Expression{value}
}

func (p *Parser) Equality() Expr {
	expr := p.Comparison()

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.Comparison()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) Comparison() Expr {
	expr := p.Term()
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right := p.Term()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) Term() Expr {
	expr := p.Factor()
	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right := p.Factor()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) Factor() Expr {
	expr := p.Unary()
	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right := p.Unary()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) Unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right := p.Unary()
		return &Unary{operator, right}
	}
	return p.Primary()
}

func (p *Parser) Primary() Expr {
	switch {
	case p.match(scanner.FALSE):
		return &Literal{false}
	case p.match(scanner.TRUE):
		return &Literal{true}
	case p.match(scanner.NIL):
		return &Literal{nil}
	case p.match(scanner.NUMBER, scanner.STRING):
		return &Literal{p.previous().Literal}
	case p.match(scanner.IDENTIFIER):
		return &Variable{p.previous()}
	case p.match(scanner.LEFT_PAREN):
		expr := p.Expression()
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
