package interpreter

import (
	"craftinginterpreters/lox/parser"
	"craftinginterpreters/lox/scanner"
	"errors"
	"fmt"
)

var _ parser.Visitor = &Interpreter{}

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expression parser.Expr) (err error) {
	defer func() {
		if terr := recover(); terr != nil {
			err = errors.New(terr.(string))
		}
	}()
	value := i.evaluate(expression)
	fmt.Println(value)
	return nil
}

func (i *Interpreter) VisitLiteralExpr(l *parser.Literal) any {
	return l.Value
}

func (i *Interpreter) VisitGroupingExpr(e *parser.Grouping) any {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitUnaryExpr(u *parser.Unary) any {
	right := i.evaluate(u.Right)

	switch u.Operator.Type {
	case scanner.BANG:
		return !i.isTruthy(right)
	case scanner.MINUS:
		return -right.(float64)
	}

	return nil
}

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if v, ok := object.(bool); ok {
		return v
	}
	return true
}

func (i *Interpreter) VisitBinaryExpr(b *parser.Binary) any {
	left := i.evaluate(b.Left)
	right := i.evaluate(b.Right)

	switch b.Operator.Type {
	case scanner.GREATER:
		i.checkNumberOperand(b.Operator, left, right)
		return left.(float64) > right.(float64)
	case scanner.GREATER_EQUAL:
		i.checkNumberOperand(b.Operator, left, right)
		return left.(float64) >= right.(float64)
	case scanner.LESS:
		i.checkNumberOperand(b.Operator, left, right)
		return left.(float64) < right.(float64)
	case scanner.LESS_EQUAL:
		i.checkNumberOperand(b.Operator, left, right)
		return left.(float64) <= right.(float64)
	case scanner.MINUS:
		i.checkNumberOperand(b.Operator, left, right)
		return left.(float64) - right.(float64)
	case scanner.PLUS:
		_, ok1 := left.(float64)
		_, ok2 := right.(float64)
		if ok1 && ok2 {
			return left.(float64) + right.(float64)
		}
		_, ok1 = left.(string)
		_, ok2 = right.(string)
		if ok1 && ok2 {
			return left.(string) + right.(string)
		}

		panic(fmt.Sprintf("[line %d ], Operands must be two number or strings.", b.Operator.Line))
	case scanner.SLASH:
		i.checkNumberOperand(b.Operator, left, right)
		return left.(float64) / right.(float64)
	case scanner.STAR:
		i.checkNumberOperand(b.Operator, left, right)
		return left.(float64) * right.(float64)
	case scanner.BANG_EQUAL:
		return !i.isEqual(left, right)
	case scanner.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}
	return nil
}

func (i *Interpreter) isEqual(left, right any) bool {
	return left == right
}

func (i *Interpreter) evaluate(expr parser.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) checkNumberOperand(operator scanner.Token, objects ...any) {
	for _, object := range objects {
		if _, ok := object.(float64); !ok {
			panic(fmt.Sprintf("[line %d ], Operand must be a number.", operator.Line))
		}
	}
}
