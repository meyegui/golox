package scanner

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(
	Type TokenType,
	Lexeme string,
	Literal any,
	Line int,
) *Token {
	return &Token{
		Type:    Type,
		Lexeme:  Lexeme,
		Literal: Literal,
		Line:    Line,
	}
}

func (t *Token) String() string {
	var literalFormatter string
	switch t.Type {
	case STRING:
		literalFormatter = " %s"

	case NUMBER:
		literalFormatter = " %f"

	default:
		literalFormatter = ""
	}

	formatString := "%-13s %-20s" + literalFormatter
	if literalFormatter == "" {
		return fmt.Sprintf(formatString, t.Type, t.Lexeme)
	}

	return fmt.Sprintf(formatString, t.Type, t.Lexeme, t.Literal)
}
