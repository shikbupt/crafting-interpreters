package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/samber/lo"
)

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

type Scanner struct {
	source   []rune
	start    int
	current  int
	line     int
	tokens   []Token
	hadError bool
}

func New() *Scanner {
	return &Scanner{
		line:   1,
		tokens: make([]Token, 0),
	}
}

func (s *Scanner) ScanAll(source string) []Token {
	s.source = []rune(source)
	s.scanTokens()
	return s.tokens
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{
		Type:    EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	})
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '!':
		s.addToken(lo.Ternary(s.match('='), BANG_EQUAL, BANG), nil)
	case '=':
		s.addToken(lo.Ternary(s.match('='), EQUAL_EQUAL, EQUAL), nil)
	case '<':
		s.addToken(lo.Ternary(s.match('='), LESS_EQUAL, LESS), nil)
	case '>':
		s.addToken(lo.Ternary(s.match('='), GREATER_EQUAL, GREATER), nil)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			s.Error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.Error(s.line, "Unterminated string.")
		return
	}
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, string(value))
}

func (s *Scanner) isDigit(c rune) bool {
	return unicode.IsDigit(c)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	value, _ := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	s.addToken(NUMBER, value)
}

func (s *Scanner) isAlpha(c rune) bool {
	// return (c >= 'a' && c <= 'z') ||
	// 	(c >= 'A' && c <= 'Z') ||
	// 	c == '_'
	return c == '_' || unicode.IsLetter(c)
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	if t, ok := Keywords[string(text)]; ok {
		s.addToken(t, nil)
		return
	}
	s.addToken(IDENTIFIER, nil)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	s.current++
	return rune(s.source[s.current-1])
}

func (s *Scanner) addToken(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		Type:    tokenType,
		Lexeme:  string(text),
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) match(c rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != c {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return unicode.ReplacementChar
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return unicode.ReplacementChar
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) Error(line int, message string) {
	s.report(line, "", message)
}

func (s *Scanner) report(line int, where, message string) {
	fmt.Printf("[line %d ] Error %s: %s", line, where, message)
	s.hadError = true
}
