package main

import (
	"bufio"
	"fmt"
	"strings"
)

type TokenType int

const (
	T_EOF TokenType = iota
	T_FUNC
	T_STRUCT
	T_TYPE
	T_LPAREN
	T_RPAREN
	T_LBRACE
	T_RBRACE
	T_LBRACKET
	T_RBRACKET
	T_ASSIGN
	T_IF
	T_ELSE
	T_VAR
	T_FOR
	T_STRING
	T_UNKNOWN
	T_SLASH
	T_COMMA
	T_DOT
	T_COLON
	T_SCOLON
	T_ASTERISK
	T_PLUS
	T_MINUS
	T_LT
	T_GT
	T_LE
	T_GE
	T_PERCENT
	T_COMMENT
	T_NAME
	T_NUMBER
	T_PACKAGE
	T_IMPORT
	T_CONST
	T_DEFINE
	T_INTERFACE
)

var TOKEN_NAMES = []string{
	"EOF",
	"FUNC",
	"STRUCT",
	"TYPE",
	"LPAREN",
	"RPAREN",
	"LBRACE",
	"RBRACE",
	"LBRACKET",
	"RBRACKET",
	"ASSIGN",
	"IF",
	"ELSE",
	"VAR",
	"FOR",
	"STRING",
	"UNKNOWN",
	"SLASH",
	"COMMA",
	"DOT",
	"COLON",
	"SCOLON",
	"ASTERISK",
	"PLUS",
	"MINUS",
	"LT",
	"GT",
	"LE",
	"GE",
	"PERCENT",
	"COMMENT",
	"NAME",
	"NUMBER",
	"PACKAGE",
	"IMPORT",
	"CONST",
	"DEFINE",
	"INTERFACE",
}

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input   *bufio.Reader
	current rune
	eof     bool
}

func NewLexer(input *bufio.Reader) *Lexer {
	l := &Lexer{input: input}
	l.readRune()
	return l
}

func (l *Lexer) readRune() {
	l.current = 0
	if l.eof {
		return
	}
	value, _, err := l.input.ReadRune()
	if err != nil {
		return
	}
	l.current = value
}

func (l *Lexer) NextToken() (Token, error) {

	l.skipWhitespace()

	var tok Token
	switch l.current {
	case '=':
		tok = Token{Type: T_ASSIGN, Value: "="}
	case '{':
		tok = Token{Type: T_LBRACE, Value: "{"}
	case '}':
		tok = Token{Type: T_RBRACE, Value: "}"}
	case '(':
		tok = Token{Type: T_LPAREN, Value: "("}
	case ')':
		tok = Token{Type: T_RPAREN, Value: ")"}
	case '[':
		tok = Token{Type: T_LBRACKET, Value: "["}
	case ']':
		tok = Token{Type: T_RBRACKET, Value: "]"}
	case '%':
		tok = Token{Type: T_PERCENT, Value: "]"}
	case '<':
		l.readRune()
		if l.current == '=' {
			return Token{Type: T_LE, Value: "<"}, nil
		} else {
			return Token{Type: T_LT, Value: "<="}, nil
		}
	case '>':
		l.readRune()
		if l.current == '=' {
			return Token{Type: T_GE, Value: ">"}, nil
		} else {
			return Token{Type: T_GT, Value: ">="}, nil
		}
	case '/':
		l.readRune()
		if l.current == '/' {
			literal := l.readLineComment()
			return Token{Type: T_COMMENT, Value: literal}, nil
		} else {
			return Token{Type: T_SLASH, Value: "/"}, nil
		}
	case '"':
		literal, err := l.readString()
		if err != nil {
			return Token{Type: T_UNKNOWN, Value: ""}, err
		}
		return Token{Type: T_STRING, Value: literal}, nil
	case '+':
		tok = Token{Type: T_PLUS, Value: "+"}
	case '-':
		tok = Token{Type: T_MINUS, Value: "-"}
	case ',':
		tok = Token{Type: T_COMMA, Value: ","}
	case '.':
		l.readRune()
		if isDigit(l.current) {
			literal := l.readNumber(".")
			return Token{Type: T_NUMBER, Value: literal}, nil
		}
		return Token{Type: T_DOT, Value: "."}, nil
	case ':':
		l.readRune()
		if l.current == '=' {
			tok = Token{Type: T_DEFINE, Value: ":="}
		} else {
			tok = Token{Type: T_COLON, Value: ":"}
		}
	case ';':
		tok = Token{Type: T_SCOLON, Value: ";"}
	case '*':
		tok = Token{Type: T_ASTERISK, Value: "*"}
	case 0:
		tok = Token{Type: T_EOF, Value: ""}
	default:
		if isIdentifierSymbol(l.current, true) {
			literal := l.readIdentifier()
			return Token{Type: lookupIdent(literal), Value: literal}, nil
		} else if isDigit(l.current) {
			literal := l.readNumber("")
			return Token{Type: T_NUMBER, Value: literal}, nil
		} else {
			return Token{Type: T_UNKNOWN, Value: ""}, fmt.Errorf("unknown token '%c'", l.current)
		}
	}

	l.readRune()
	return tok, nil
}

func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
		l.readRune()
	}
}

func (l *Lexer) readString() (string, error) {
	l.readRune() // skip quotes
	var builder strings.Builder
	for l.current != '"' && l.current != 0 {
		builder.WriteRune(l.current)
		l.readRune()
	}
	if l.current == 0 {
		return "", fmt.Errorf("unterminated string")
	}
	l.readRune()
	return builder.String(), nil
}

func (l *Lexer) readLineComment() string {
	l.readRune() // skip '/'
	var builder strings.Builder
	for l.current != '\n' && l.current != 0 {
		builder.WriteRune(l.current)
		l.readRune()
	}
	return builder.String()
}

func (l *Lexer) readIdentifier() string {
	var builder strings.Builder
	for isIdentifierSymbol(l.current, false) {
		builder.WriteRune(l.current)
		l.readRune()
	}
	return builder.String()
}

func (l *Lexer) readNumber(prefix string) string {
	var builder strings.Builder
	builder.WriteString(prefix)
	for isDigit(l.current) || l.current == '.' || l.current == 'e' || l.current == 'E' {
		builder.WriteRune(l.current)
		l.readRune()
	}
	return builder.String()
}

func isIdentifierSymbol(ch rune, start bool) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || (!start && (ch == '_' || isDigit(ch)))
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	keywords := map[string]TokenType{
		"var":       T_VAR,
		"const":     T_CONST,
		"func":      T_FUNC,
		"struct":    T_STRUCT,
		"type":      T_TYPE,
		"if":        T_IF,
		"else":      T_ELSE,
		"for":       T_FOR,
		"package":   T_PACKAGE,
		"import":    T_IMPORT,
		"interface": T_INTERFACE,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return T_NAME
}
