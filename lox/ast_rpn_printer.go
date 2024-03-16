package lox

import (
	"fmt"
	"strings"

	"github.com/meyegui/golox/ast"
)

type AstRpnPrinter struct {
	output []string
}

// AstRpnPrinter is a Reverse Polish Notation printer for an AST.
func (arp *AstRpnPrinter) Print(expr ast.Expr) string {
	expr.Accept(arp)

	return strings.Join(arp.output, " ")
}

func (arp *AstRpnPrinter) VisitUnary(u *ast.Unary) {
	arp.output = append(arp.output, new(AstRpnPrinter).Print(u.Right))
	arp.output = append(arp.output, u.Operator.Lexeme)
}

func (arp *AstRpnPrinter) VisitBinary(b *ast.Binary) {
	arp.output = append(arp.output, new(AstRpnPrinter).Print(b.Left))
	arp.output = append(arp.output, new(AstRpnPrinter).Print(b.Right))
	arp.output = append(arp.output, b.Operator.Lexeme)
}

func (arp *AstRpnPrinter) VisitGrouping(g *ast.Grouping) {
	arp.output = append(arp.output, new(AstRpnPrinter).Print(g.Expr))
}

func (arp *AstRpnPrinter) VisitLiteral(l *ast.Literal) {
	arp.output = append(arp.output, fmt.Sprintf("%v", l.Value))
}
