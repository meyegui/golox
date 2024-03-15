package scanner

import (
	"fmt"
	"strconv"

	"github.com/meyegui/golox/contracts"
)

// ------------------------------------------------------------ TYPE DEFINITION

type Scanner struct {
	source  []rune
	tokens  []Token
	start   int
	current int
	line    int

	errorReporter contracts.ErrorReporter
}

// ------------------------------------------------------------ EXPORTED DATA

var Keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// ------------------------------------------------------------ EXPORTED FUNCTIONS

func NewScanner(source string, errorReporter contracts.ErrorReporter) *Scanner {
	return &Scanner{
		source:        []rune(source),
		tokens:        []Token{},
		start:         0,
		current:       0,
		line:          1,
		errorReporter: errorReporter,
	}
}

// ------------------------------------------------------------ EXPORTED METHODS

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, *NewToken(EOF, "", nil, s.line))

	return s.tokens
}

// ------------------------------------------------------------ UNEXPORTED METHODS

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(LEFT_PAREN)

	case ')':
		s.addToken(RIGHT_PAREN)

	case '{':
		s.addToken(LEFT_BRACE)

	case '}':
		s.addToken(RIGHT_BRACE)

	case ',':
		s.addToken(COMMA)

	case '.':
		s.addToken(DOT)

	case '+':
		s.addToken(PLUS)

	case '-':
		s.addToken(MINUS)

	case ';':
		s.addToken(SEMICOLON)

	case '*':
		s.addToken(STAR)

	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}

	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}

	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}

	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}

	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace
		break

	case '\n':
		s.line++

	case '"':
		s.scanString()

	default:
		if isDigit(char) {
			s.scanNumber()
		} else if isAlpha(char) {
			s.scanIdentifier()
		} else {
			s.errorReporter.Error(s.line, fmt.Sprintf("Unexpected character '%c'.", char))
		}
	}
}

// ---------------------------------------- HELPER METHODS

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	char := s.source[s.current]
	s.current++

	return char
}

func (s *Scanner) match(char rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != char {
		return false
	}

	s.current++

	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	index := s.current + 1
	if index > len(s.source) {
		return 0
	}

	return s.source[index]
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		s.errorReporter.Error(s.line, "Unterminated string.")
		return
	}

	// Consume the closing quotes
	s.advance()

	// Trim the surrounding quotes
	value := string(s.source[s.start+1 : s.current-1])
	s.addLiteralToken(STRING, value)
}

func (s *Scanner) scanNumber() {
	// Consume the integer part
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the dot
		s.advance()

		// Consume the fractional part
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	strValue := string(s.source[s.start:s.current])
	floatValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		s.errorReporter.Error(s.line, fmt.Sprintf("Invalid number literal: %s", strValue))
	} else {
		s.addLiteralToken(NUMBER, floatValue)
	}
}

func (s *Scanner) scanIdentifier() {
	// Consume the whole identifier
	for isAlphaNumeric(s.peek()) && !s.isAtEnd() {
		s.advance()
	}

	// Check if the identifier is a reserved keyword
	identifier := string(s.source[s.start:s.current])
	tokenType, isKeyword := Keywords[identifier]
	if !isKeyword {
		tokenType = IDENTIFIER
	}

	s.addToken(tokenType)
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isAlpha(char rune) bool {
	return (char >= 'A' && char <= 'Z') ||
		(char >= 'a' && char <= 'z') ||
		char == '_'
}

func isAlphaNumeric(char rune) bool {
	return isAlpha(char) || isDigit(char)
}

func (s *Scanner) addLiteralToken(tt TokenType, literal any) {
	lexeme := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, *NewToken(tt, lexeme, literal, s.line))
}

func (s *Scanner) addToken(tt TokenType) {
	s.addLiteralToken(tt, nil)
}
