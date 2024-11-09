package main

import (
	"fmt"
)

type Module struct {
	Package string
	Body    []interface{}
}

type ConstantDecl struct {
	ParameterDecl
	Values []string
}

type VariableDecl struct {
	ParameterDecl
	Values []string
}

type ParameterDecl struct {
	Names []string
	Type  string
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
	lexer  *Lexer
	token  Token
	errors []string
	module *Module
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{lexer: l, module: &Module{}}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.token, _ = p.lexer.NextToken()
	fmt.Printf("-- %s\n", p.token)
}

func (p *Parser) expected(tokens ...TokenType) bool {
	for _, t := range tokens {
		if p.token.Type == t {
			p.nextToken()
			return true
		}
	}
	return false
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.token.Type == t {
		p.nextToken()
		return true
	} else {
		err := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.token.Type)
		p.errors = append(p.errors, err)
		fmt.Println(err)
		return false
	}
}

func (p *Parser) Parse() *Module {
	for p.token.Type != T_EOF {
		switch p.token.Type {
		case T_FOR:
			p.module.Body = append(p.module.Body, p.parseFor())
		case T_IF:
			p.module.Body = append(p.module.Body, p.parseIf())
		case T_FUNC:
			p.module.Body = append(p.module.Body, p.parseFunction())
		default:
			p.errors = append(p.errors, fmt.Sprintf("unrecognized token '%s'", p.token))
			return nil
		}
	}
	return p.module
}

func (p *Parser) parseVarDeclaration() *VariableDecl {
	return nil
}

func (p *Parser) parseFunction() *Function {
	return nil
}

func (p *Parser) parseStruct() *Struct {
	return nil
}

func (p *Parser) parseIf() *IfStmt {
	p.nextToken() // Skip 'if'
	if !p.expectPeek(T_LPAREN) {
		return nil
	}
	condition := p.parseExpression()
	if !p.expectPeek(T_RPAREN) {
		return nil
	}
	if !p.expectPeek(T_LBRACE) {
		return nil
	}
	body := p.parseBlock()
	if !p.expectPeek(T_RBRACE) {
		return nil
	}
	return &IfStmt{Condition: condition, Body: body}
}

func (p *Parser) parseExpression() interface{} {
	return p.parseUnaryExpr()
}

// ParseUnary = unary_op UnaryExpr | PrimaryExpr .
func (p *Parser) parseUnaryExpr() interface{} {
	// has a 'unary_op'?
	operator := p.token.Type
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
	switch p.token.Type {
	case T_NAME:
		return &NameLit{Literal{Value: p.token.Value}}
	case T_STRING:
		return &StringLit{Literal{Value: p.token.Value}}
	case T_INTEGER:
		return &StringLit{Literal{Value: p.token.Value}}
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
