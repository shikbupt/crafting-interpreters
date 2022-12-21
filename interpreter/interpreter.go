package interpreter

import (
	"craftinginterpreters/lox/parser"
	"craftinginterpreters/lox/scanner"
	"errors"
	"fmt"
)

var _ parser.ExprVisitor = &Interpreter{}
var _ parser.StmtVisitor = &Interpreter{}

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnv(nil),
	}
}

func (i *Interpreter) Interpret(statements []parser.Stmt) (err error) {
	defer func() {
		if terr := recover(); terr != nil {
			err = errors.New(terr.(string))
		}
	}()
	for _, statement := range statements {
		i.evaluateStmt(statement)
	}
	return nil
}

func (i *Interpreter) VisitLiteralExpr(l *parser.Literal) any {
	return l.Value
}

func (i *Interpreter) VisitGroupingExpr(e *parser.Grouping) any {
	return i.evaluateExpr(e.Expression)
}

func (i *Interpreter) VisitUnaryExpr(u *parser.Unary) any {
	right := i.evaluateExpr(u.Right)

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
	left := i.evaluateExpr(b.Left)
	right := i.evaluateExpr(b.Right)

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

func (i *Interpreter) VisitVariableExpr(v *parser.Variable) any {
	return i.env.get(v.Name)
}

func (i *Interpreter) isEqual(left, right any) bool {
	return left == right
}

func (i *Interpreter) evaluateExpr(expr parser.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) evaluateStmt(expr parser.Stmt) any {
	return expr.Accept(i)
}

func (i *Interpreter) checkNumberOperand(operator scanner.Token, objects ...any) {
	for _, object := range objects {
		if _, ok := object.(float64); !ok {
			panic(fmt.Sprintf("[line %d ], Operand must be a number.", operator.Line))
		}
	}
}

func (i *Interpreter) VisitExpressionStmt(e *parser.Expression) any {
	i.evaluateExpr(e.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(p *parser.Print) any {
	value := i.evaluateExpr(p.Expression)
	fmt.Println(value)
	return nil
}

func (i *Interpreter) VisitVarStmt(v *parser.Var) any {
	var value any
	if v.Initializer != nil {
		value = i.evaluateExpr(v.Initializer)
	}
	i.env.define(v.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitAssignExpr(a *parser.Assign) any {
	value := i.evaluateExpr(a.Value)
	i.env.assign(a.Name, value)
	return value
}

func (i *Interpreter) VisitBlockStmt(b *parser.Block) any {
	i.executeBlock(b.Statements, NewEnv(i.env))
	return nil
}

func (i *Interpreter) executeBlock(statements []parser.Stmt, env *Environment) {
	parentEnv := i.env
	defer func() {
		i.env = parentEnv
	}()
	i.env = env
	for _, statement := range statements {
		i.evaluateStmt(statement)
	}
}
