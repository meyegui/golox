package lox

import (
	"fmt"
	"strings"

	"github.com/meyegui/golox/ast"
)

// AstPrinter is a Lisp-style printer for an AST.
type AstPrinter struct {
	Output string
}

func (ap *AstPrinter) Print(expr ast.Expr) string {
	expr.Accept(ap)

	return ap.Output
}

func (ap *AstPrinter) VisitUnary(u *ast.Unary) {
	ap.Output = parenthesize(u.Operator.Lexeme, u.Right)
}

func (ap *AstPrinter) VisitBinary(b *ast.Binary) {
	ap.Output = parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (ap *AstPrinter) VisitGrouping(g *ast.Grouping) {
	ap.Output = parenthesize("group", g.Expr)
}

func (ap *AstPrinter) VisitLiteral(l *ast.Literal) {
	ap.Output = fmt.Sprintf("%v", l.Value)
}

func parenthesize(name string, exprs ...ast.Expr) string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteByte(' ')

		strExpr := (&AstPrinter{}).Print(expr)
		sb.WriteString(strExpr)
	}
	sb.WriteByte(')')

	return sb.String()
}
