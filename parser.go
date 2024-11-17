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

type VariableGroup struct {
	Variables []*VariableDecl
	IsConst   bool
}

type VariableDecl struct {
	Names   []string
	Type    string
	Values  []interface{}
	IsConst bool
}

type Function struct {
	Name       string
	Parameters []*ParameterDecl
	Body       []interface{}
	Return     interface{}
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
	O_ADD           // +
	O_SUB           // -
	O_MUL           // *
	O_DIV           // /
	O_EQ            // ==
	O_NE            // !=
	O_LT            // <
	O_GT            // >
	O_LE            // <=
	O_GE            // >=
	O_AS            // =
	O_DREF          // *
)

const LOOKAHEAD_SIZE = 5

type TypeDecl struct {
	IsPointer bool
	Type      interface{}
	Length    interface{}
}

type QualifiedName struct {
	Names []interface{}
}

type UnaryExpr struct {
	Right    interface{}
	Operator Operator
}

type BinaryExpr struct {
	Left     interface{}
	Right    interface{}
	Operator Operator
}

type AssignStmt struct {
	BinaryExpr
}

type Literal struct {
	Value string
}

type NameLit struct {
	Value *QualifiedName
}

type StringLit struct {
	Literal
}

type NumberLit struct {
	Literal
}

type CallExpr struct {
	Callee  interface{}
	ArgList []interface{}
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
	for i := 0; i < count; {
		token, err := p.lexer.NextToken()
		if err != nil {
			return false
		}
		// ignoring comments
		if token.Type == T_COMMENT {
			continue
		}
		fmt.Printf("-- %s [%s]\n", TOKEN_NAMES[token.Type], token.Value)
		p.lookahead = append(p.lookahead, token)
		i++
	}
	return true
}

func (p *Parser) expectedOneOf(tokens ...TokenType) bool {
	for _, t := range tokens {
		if p.expected(t) {
			return true
		}
	}
	return false
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

func (p *Parser) required(t TokenType) (bool, Token) {
	if len(p.lookahead) == 0 {
		p.refill(1)
	}

	if p.lookahead[0].Type == t {
		token := p.lookahead[0]
		p.nextToken()
		return true, token
	} else {
		err := fmt.Sprintf("expected next token to be %s, got %s instead", TOKEN_NAMES[t], TOKEN_NAMES[p.lookahead[0].Type])
		p.errors = append(p.errors, err)
		fmt.Println(err)
		return false, Token{}
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
		case T_CONST:
			fallthrough
		case T_VAR:
			stmt := p.parseVarDeclaration()
			if stmt == nil {
				break loop
			}
			p.module.Body = append(p.module.Body, stmt)
		case T_FOR:
			stmt := p.parseFor()
			if stmt == nil {
				break loop
			}
			p.module.Body = append(p.module.Body, stmt)
		case T_IF:
			stmt := p.parseIf()
			if stmt == nil {
				break loop
			}
			p.module.Body = append(p.module.Body, stmt)
		case T_FUNC:
			stmt := p.parseFunction()
			if stmt == nil {
				break loop
			}
			p.module.Body = append(p.module.Body, stmt)
		case T_COMMENT:
			p.nextToken()
			// ignore comments
		case T_EOF:
			break loop
		default:
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

func (p *Parser) parseParameterDecl() *ParameterDecl {
	return &ParameterDecl{}
}

func (p *Parser) parseVarDeclaration() interface{} {
	is_const := p.peekType() == T_CONST
	p.nextToken()

	if p.expected(T_LPAREN) {
		result := &VariableGroup{IsConst: is_const}
		for p.peekType() != T_RPAREN {
			stmt := p.parseSingleVariable()
			if stmt == nil {
				break
			}
			result.Variables = append(result.Variables, stmt)
		}
		if ok, _ := p.required(T_RPAREN); !ok {
			return nil
		}
		return result
	} else {
		result := p.parseSingleVariable()
		if result == nil {
			return nil
		}
		result.IsConst = is_const
		return result
	}
}

func (p *Parser) parseSingleVariable() *VariableDecl {
	result := &VariableDecl{}
	// variable names
	for {
		literal := p.peek().Value
		if ok, _ := p.required(T_NAME); !ok {
			return nil
		}
		result.Names = append(result.Names, literal)
		if !p.expected(T_COMMA) {
			break
		}
	}
	// variable type
	if p.peekType() != T_ASSIGN {
		result.Type = p.peek().Value
		if ok, _ := p.required(T_NAME); !ok {
			return nil
		}
	}
	// initialization
	if p.expected(T_ASSIGN) {
		for {
			value := p.parseExpression()
			if value == nil {
				return nil
			}
			result.Values = append(result.Values, value)
			if !p.expected(T_COMMA) {
				break
			}
		}
	}
	return result
}

func (p *Parser) parseFunction() *Function {
	function := &Function{}
	p.nextToken() // skip 'func' keyword

	// function name
	if ok, token := p.required(T_NAME); ok {
		function.Name = token.Value
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected name of function but found '%s'", TOKEN_NAMES[p.peekType()]))
		return nil
	}
	// function parameters
	function.Parameters = p.parseParameters()
	// optional return type
	if p.peekType() != T_LBRACE {
		if p.peekType() == T_LPAREN {
			function.Return = p.parseParameters()
		} else {
			function.Return = p.parseType()
		}
	}
	// function body
	function.Body = p.parseBlock()

	return function
}

func (p *Parser) parseParameters() []*ParameterDecl {
	result := make([]*ParameterDecl, 0)
	if ok, _ := p.required(T_LPAREN); !ok {
		return nil
	}
	for p.peekType() != T_RPAREN {
		result = append(result, p.parseParameterDecl())
		if !p.expected(T_COMMA) {
			break
		}
	}
	if ok, _ := p.required(T_RPAREN); !ok {
		return nil
	}
	return result
}

func (p *Parser) parseType() *TypeDecl {
	result := &TypeDecl{}
	switch p.peekType() {
	case T_LBRACKET:
		if !p.expected(T_RBRACKET) {
			result.Length = p.parseExpression()
		}
		if ok, _ := p.required(T_RBRACKET); !ok {
			return nil
		}
		result.Type = p.parseType()
	case T_NAME:
		result.Type = p.parseQualifiedName()
	case T_INTERFACE:
		p.nextToken()
		if !p.expected(T_LBRACE, T_RBRACE) {
			p.errors = append(p.errors, fmt.Sprintf("expected '{}' but found '%s'", TOKEN_NAMES[p.peekType()]))
			return nil
		}
	case T_ASTERISK:
		p.nextToken()
		result.IsPointer = true
		result.Type = p.parseType()
	}

	return result
}

func (p *Parser) parseStruct() *Struct {
	p.errors = append(p.errors, "struct parsing not implemented")
	return nil
}

func (p *Parser) parseIf() *IfStmt {
	p.errors = append(p.errors, "conditional parsing not implemented")
	return nil
}

func (p *Parser) parseExpression() interface{} {
	return p.parseUnaryExpr()
}

// ParseUnary = unary_op UnaryExpr | PrimaryExpr .
func (p *Parser) parseUnaryExpr() interface{} {
	// has a 'unary_op'?
	operator := p.peekType()
	if !p.expectedOneOf(T_PLUS, T_MINUS, T_ASTERISK) {
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
	case T_ASTERISK:
		return O_DREF
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
		name := p.parseQualifiedName()
		if p.peekType() == T_LPAREN {
			return p.parseCallExpr(name)
		}
		return &NameLit{Value: name}
	case T_STRING:
		p.nextToken()
		return &StringLit{Literal{Value: token.Value}}
	case T_NUMBER:
		p.nextToken()
		return &NumberLit{Literal{Value: token.Value}}
	default:
		return nil
	}
}

func (p *Parser) parseBinaryExpr() interface{} {
	return nil
}

func (p *Parser) parseBlock() []interface{} {
	body := []interface{}{}
	if ok, _ := p.required(T_LBRACE); !ok {
		return nil
	}
loop:
	for {
		ttype := p.peekType()

		switch ttype {
		case T_FOR:
			stmt := p.parseFor()
			if stmt == nil {
				break loop
			}
			body = append(body, stmt)
		case T_RBRACE:
			p.nextToken()
			break loop
		case T_VAR:
			fallthrough
		case T_CONST:
			stmt := p.parseVarDeclaration()
			if stmt == nil {
				break loop
			}
			body = append(body, stmt)
		//case T_RETURN:
		//	stmt := p.parseReturn()
		//	if stmt == nil {
		//		break loop
		//	}
		//	body = append(body, stmt)
		case T_IF:
			stmt := p.parseIf()
			if stmt == nil {
				break loop
			}
			body = append(body, stmt)
		case T_NAME:
			body = append(body, p.parseNameStatement())
		default:
			p.errors = append(p.errors, fmt.Sprintf("unexpected token %s", TOKEN_NAMES[ttype]))
			return nil
		}
	}
	return body
}

func (p *Parser) parseAssignement(name *QualifiedName) interface{} {
	result := &AssignStmt{}
	result.Left = name
	result.Operator = deduceOperator(p.peekType())
	if !p.expectedOneOf(T_ASSIGN, T_DEFINE) {
		p.errors = append(p.errors, fmt.Sprintf("expected assignment operator but found '%s'", TOKEN_NAMES[p.peekType()]))
		return nil
	}
	result.Right = p.parseExpression()
	if result.Right == nil {
		return nil
	}
	return result
}

func (p *Parser) ahead(index int) Token {
	if index < 0 || index >= LOOKAHEAD_SIZE {
		return Token{Type: T_EOF}
	}
	if index+1 > len(p.lookahead) {
		p.refill(index + 1)
	}
	return p.lookahead[index]
}

func (p *Parser) parseNameStatement() interface{} {
	name := p.parseQualifiedName()
	// is it a function call?
	if p.peekType() == T_LPAREN {
		return p.parseCallExpr(name)
	}
	if p.peekType() == T_ASSIGN || p.peekType() == T_DEFINE {
		return p.parseAssignement(name)
	}

	// is it a return statement?
	//if tahead.Type == T_RETURN {
	//	return p.parseReturn()
	//}
	p.errors = append(p.errors, fmt.Sprintf("unexpected token %s", TOKEN_NAMES[p.peekType()]))
	return nil
}

func (p *Parser) parseQualifiedName() *QualifiedName {
	if p.peekType() != T_NAME {
		p.errors = append(p.errors, fmt.Sprintf("expected name but found '%s'", TOKEN_NAMES[p.peekType()]))
		return nil
	}
	qname := &QualifiedName{}
	for p.peekType() == T_NAME {
		qname.Names = append(qname.Names, p.peek().Value)
		p.nextToken()
		if !p.expected(T_DOT) {
			break
		}
	}
	return qname
}

func (p *Parser) parseCallExpr(name *QualifiedName) interface{} {
	result := &CallExpr{}
	result.Callee = name
	if result.Callee == nil {
		return nil
	}
	if ok, _ := p.required(T_LPAREN); !ok {
		return nil
	}
	for p.peekType() != T_RPAREN {
		result.ArgList = append(result.ArgList, p.parseExpression())
		if !p.expected(T_COMMA) {
			break
		}
	}
	if ok, _ := p.required(T_RPAREN); !ok {
		return nil
	}
	return result
}

func (p *Parser) parseFor() interface{} {
	p.errors = append(p.errors, "for loop parsing not implemented")
	return nil
}
