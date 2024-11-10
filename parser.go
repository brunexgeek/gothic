package main

import (
	"fmt"
)

type Module struct {
	Package string
	Imports *ImportStmt
	Body    []interface{}
}

type ImportStmt struct {
	Names []string
}

type ParameterDecl struct {
	Names []string
	Type  string
}

type VariableDecl struct {
	Names   []string
	Type    string
	Values  []string
	IsConst bool
}

type Function struct {
	Name       string
	Parameters []*ParameterDecl
	Body       []interface{}
}

type Struct struct {
	Name   string
	Fields []*VariableDecl
}

type Expression struct {
}

type Operator int

const (
	O_NONE Operator = iota
	O_ADD
	O_SUB
	O_MUL
	O_DIV
	O_EQ
	O_NE
	O_LT
	O_GT
	O_LE
	O_GE
)

const LOOKAHEAD_SIZE = 5

type UnaryExpr struct {
	Right    interface{}
	Operator Operator
}

type BinaryExpr struct {
	Left     interface{}
	Right    interface{}
	Operator Operator
}

type Literal struct {
	Value string
}

type NameLit struct {
	Literal
}

type StringLit struct {
	Literal
}

type IntegerLit struct {
	Literal
}

type IfStmt struct {
	Condition interface{}
	Body      []interface{}
}

type ForStmt struct {
	Initializer *Expression
	Condition   *Expression
	Step        *Expression
	Body        []interface{}
}

type ForRangeStmt struct {
	Variables []string
	Create    bool
	Range     *Expression
	Body      []interface{}
}

type Parser struct {
	lexer *Lexer
	//token     Token
	errors    []string
	module    *Module
	lookahead []Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{lexer: l, module: &Module{}}
	p.nextToken()
	return p
}

func (p *Parser) peek() Token {
	if len(p.lookahead) == 0 {
		var zero Token
		return zero
	}
	return p.lookahead[0]
}

func (p *Parser) peekType() TokenType {
	if len(p.lookahead) == 0 {
		return T_EOF
	}
	return p.lookahead[0].Type
}

func (p *Parser) nextToken() bool {
	var ok = true
	if len(p.lookahead) == 0 {
		ok = p.refill(1)
	} else if len(p.lookahead) == 1 {
		ok = p.refill(2)
		if ok {
			p.lookahead = p.lookahead[1:]
		}
	} else {
		p.lookahead = p.lookahead[1:]
	}
	return ok
}

// Refill the lookahead buffer so it contains at least 'c' tokens
func (p *Parser) refill(c int) bool {
	count := min(LOOKAHEAD_SIZE, c) - len(p.lookahead)
	for i := 0; i < count; i++ {
		token, err := p.lexer.NextToken()
		fmt.Printf("-- %s\n", TOKEN_NAMES[token.Type])
		if err != nil {
			return false
		}
		p.lookahead = append(p.lookahead, token)
	}
	return true
}

func (p *Parser) expected(tokens ...TokenType) bool {
	if len(tokens)+1 > LOOKAHEAD_SIZE {
		return false
	}
	// get enough tokens to satisfy the lookahead length
	if len(p.lookahead) < len(tokens)+1 {
		p.refill(len(tokens) + 1)
	}

	for i, t := range tokens {
		if p.lookahead[i].Type != t {
			return false
		}
	}
	p.lookahead = p.lookahead[len(tokens):]
	return true
}

func (p *Parser) required(t TokenType) bool {
	if len(p.lookahead) == 0 {
		p.refill(1)
	}

	if p.lookahead[0].Type == t {
		p.nextToken()
		return true
	} else {
		err := fmt.Sprintf("expected next token to be %s, got %s instead", TOKEN_NAMES[t], TOKEN_NAMES[p.lookahead[0].Type])
		p.errors = append(p.errors, err)
		fmt.Println(err)
		return false
	}
}

func (p *Parser) Parse() *Module {
loop:
	for {
		ttype := p.peekType()
		switch ttype {
		case T_PACKAGE:
			p.module.Package = p.parsePackage()
		case T_IMPORT:
			p.module.Imports = p.parseImport()
		case T_VAR:
			p.module.Body = append(p.module.Body, p.parseVarDeclaration())
		case T_FOR:
			p.module.Body = append(p.module.Body, p.parseFor())
		case T_IF:
			p.module.Body = append(p.module.Body, p.parseIf())
		case T_FUNC:
			p.module.Body = append(p.module.Body, p.parseFunction())
		case T_COMMENT:
			p.nextToken()
			// ignore comments
		case T_EOF:
			break loop
		default:
			p.parseVarDeclaration()
			p.errors = append(p.errors, fmt.Sprintf("unrecognized token '%s'", TOKEN_NAMES[ttype]))
			return nil
		}
	}
	return p.module
}

func (p *Parser) parseImport() *ImportStmt {
	result := &ImportStmt{}
	if p.expected(T_IMPORT) {
		// have we got multiple imports?
		if p.expected(T_LPAREN) {
			for p.peekType() == T_STRING {
				result.Names = append(result.Names, p.peek().Value)
				p.nextToken()
			}
			if !p.expected(T_RPAREN) {
				p.errors = append(p.errors, fmt.Sprintf("expected ')' but found '%s'", TOKEN_NAMES[p.peekType()]))
				return nil
			}
			return result
		}

		// single import
		if p.peekType() == T_STRING {
			result.Names = append(result.Names, p.peek().Value)
			p.nextToken()
			return result
		}

		p.errors = append(p.errors, fmt.Sprintf("expected ')' but found '%s'", TOKEN_NAMES[p.peekType()]))
		return nil
	}
	return nil
}

func (p *Parser) parsePackage() string {
	if p.expected(T_PACKAGE) && p.peekType() == T_NAME {
		literal := p.peek().Value
		p.nextToken()
		return literal
	}
	return ""
}

func (p *Parser) parseParameterDecl() interface{} {
	return &ParameterDecl{}
}

func (p *Parser) parseVarDeclaration() interface{} {
	return nil
}

func (p *Parser) parseFunction() *Function {
	return nil
}

func (p *Parser) parseStruct() *Struct {
	return nil
}

func (p *Parser) parseIf() *IfStmt {
	return nil
}

func (p *Parser) parseExpression() interface{} {
	return p.parseUnaryExpr()
}

// ParseUnary = unary_op UnaryExpr | PrimaryExpr .
func (p *Parser) parseUnaryExpr() interface{} {
	// has a 'unary_op'?
	operator := p.peekType()
	if !p.expected(T_PLUS, T_MINUS) {
		return p.parsePrimaryExpr()
	}

	right := p.parseUnaryExpr()
	if right == nil {
		return nil
	}
	return &UnaryExpr{Operator: deduceOperator(operator), Right: right}
}

func deduceOperator(token TokenType) Operator {
	switch token {
	case T_PLUS:
		return O_ADD
	case T_MINUS:
		return O_SUB
	default:
		return O_NONE
	}
}

// PrimaryExpr = Operand
func (p *Parser) parsePrimaryExpr() interface{} {
	return p.parseOperand()
}

// Operand = Literal | OperandName .
func (p *Parser) parseOperand() interface{} {
	token := p.peek()
	switch token.Type {
	case T_NAME:
		return &NameLit{Literal{Value: token.Value}}
	case T_STRING:
		return &StringLit{Literal{Value: token.Value}}
	case T_INTEGER:
		return &StringLit{Literal{Value: token.Value}}
	default:
		return nil
	}
}

func (p *Parser) parseBinaryExpr() interface{} {
	return nil
}

func (p *Parser) parseBlock() []interface{} {
	return nil
}

func (p *Parser) parseFor() interface{} {
	return nil
}
