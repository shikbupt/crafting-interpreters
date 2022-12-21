package parser

import (
	"fmt"
	"strings"
)

var _ ExprVisitor = AstPrinter{}

type AstPrinter struct {
}

func (a AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) VisitBinaryExpr(b *Binary) any {
	return a.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (a AstPrinter) VisitGroupingExpr(g *Grouping) any {
	return a.parenthesize("group", g.Expression)
}

func (a AstPrinter) VisitLiteralExpr(l *Literal) any {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprint(l.Value)
}

func (a AstPrinter) VisitUnaryExpr(u *Unary) any {
	return a.parenthesize(u.Operator.Lexeme, u.Right)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	builder := strings.Builder{}

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(a).(string))
	}
	builder.WriteString(")")

	return builder.String()
}

func (a AstPrinter) VisitVariableExpr(v *Variable) any {
	return nil
}

func (a AstPrinter) VisitAssignExpr(v *Assign) any {
	return nil
}
