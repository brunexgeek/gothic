package main

import (
	"fmt"
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
	T_INTEGER
	T_PACKAGE
	T_IMPORT
	T_CONST
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
	"INTEGER",
	"PACKAGE",
	"IMPORT",
	"CONST",
}

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() (Token, error) {

	l.skipWhitespace()

	var tok Token
	switch l.ch {
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
		l.readChar()
		if l.ch == '=' {
			return Token{Type: T_LE, Value: "<"}, nil
		} else {
			return Token{Type: T_LT, Value: "<="}, nil
		}
	case '>':
		l.readChar()
		if l.ch == '=' {
			return Token{Type: T_GE, Value: ">"}, nil
		} else {
			return Token{Type: T_GT, Value: ">="}, nil
		}
	case '/':
		l.readChar()
		if l.ch == '/' {
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
		tok = Token{Type: T_DOT, Value: "."}
	case ':':
		tok = Token{Type: T_COLON, Value: ":"}
	case ';':
		tok = Token{Type: T_SCOLON, Value: ";"}
	case '*':
		tok = Token{Type: T_ASTERISK, Value: "*"}
	case 0:
		tok = Token{Type: T_EOF, Value: ""}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			return Token{Type: lookupIdent(literal), Value: literal}, nil
		} else if isDigit(l.ch) {
			literal := l.readNumber()
			return Token{Type: T_INTEGER, Value: literal}, nil
		} else {
			return Token{Type: T_UNKNOWN, Value: ""}, fmt.Errorf("unknown token '%c'", l.ch)
		}
	}

	l.readChar()
	return tok, nil
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() (string, error) {
	l.readChar() // skip quotes
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	if l.ch == 0 {
		return "", fmt.Errorf("unterminated string")
	}
	end := l.position
	l.readChar()
	return l.input[start:end], nil
}

func (l *Lexer) readLineComment() string {
	l.readChar() // skip '/'
	start := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	keywords := map[string]TokenType{
		"var":     T_VAR,
		"func":    T_FUNC,
		"struct":  T_STRUCT,
		"type":    T_TYPE,
		"if":      T_IF,
		"else":    T_ELSE,
		"for":     T_FOR,
		"package": T_PACKAGE,
		"import":  T_IMPORT,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return T_NAME
}
