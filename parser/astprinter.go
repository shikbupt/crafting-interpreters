package parser

import (
	"fmt"
	"strings"
)

var _ Visitor = AstPrinter{}

type AstPrinter struct {
}

func (a AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) VisitorBinary(b *Binary) any {
	return a.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (a AstPrinter) VisitorGrouping(g *Grouping) any {
	return a.parenthesize("group", g.Expression)
}

func (a AstPrinter) VisitorLiteral(l *Literal) any {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprint(l.Value)
}

func (a AstPrinter) VisitorUnary(u *Unary) any {
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
