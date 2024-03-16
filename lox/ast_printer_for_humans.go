package lox

import (
	"fmt"

	"github.com/meyegui/golox/ast"
)

// AstPrinterForHumans is a human-friendly printer for an AST.
type AstPrinterForHumans struct {
	Output string
}

func (ap *AstPrinterForHumans) Print(expr ast.Expr) string {
	expr.Accept(ap)

	return ap.Output
}

func (ap *AstPrinterForHumans) VisitUnary(u *ast.Unary) {
	ap.Output = u.Operator.Lexeme + (&AstPrinterForHumans{}).Print(u.Right)
}

func (ap *AstPrinterForHumans) VisitBinary(b *ast.Binary) {
	ap.Output = (&AstPrinterForHumans{}).Print(b.Left) +
		" " + b.Operator.Lexeme + " " +
		(&AstPrinterForHumans{}).Print(b.Right)
}

func (ap *AstPrinterForHumans) VisitGrouping(g *ast.Grouping) {
	ap.Output = "(" +
		(&AstPrinterForHumans{}).Print(g.Expr) +
		")"
}

func (ap *AstPrinterForHumans) VisitLiteral(l *ast.Literal) {
	ap.Output = fmt.Sprintf("%v", l.Value)
}
