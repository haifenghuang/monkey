package parser

import (
	"fmt"
	"io/ioutil"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"os"
	"strconv"
	"strings"
)

const (
	_ int = iota
	LOWEST
	PIPE
	ASSIGN
	CONDOR
	CONDAND
	EQUALS
	LESSGREATER
	BITOR
	BITXOR
	BITAND
	SHIFTS
	SLICE
	TERNARY
	DOTDOT
	SUM
	PRODUCT
	PREFIX
	MATCHING
	CALL
	INDEX
	INCREMENT
)

var precedences = map[token.TokenType]int{
	token.PIPE:       PIPE,
	token.ASSIGN:     ASSIGN,
	token.CONDOR:     CONDOR,
	token.OR:         CONDOR,
	token.AND:        CONDAND,
	token.CONDAND:    CONDAND,
	token.EQ:         EQUALS,
	token.NEQ:        EQUALS,
	token.LT:         LESSGREATER,
	token.LE:         LESSGREATER,
	token.GT:         LESSGREATER,
	token.GE:         LESSGREATER,
	token.BITOR:      BITOR,
	token.BITOR_A:    BITOR,
	token.BITXOR_A:   BITXOR,
	token.BITXOR:     BITXOR,
	token.BITAND_A:   BITAND,
	token.BITAND:     BITAND,
	token.SHIFT_L:    SHIFTS,
	token.SHIFT_R:    SHIFTS,
	token.COLON:      SLICE,
	token.QUESTIONM:  TERNARY,
	token.DOTDOT:     DOTDOT,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.PLUS_A:     SUM,
	token.MINUS_A:    SUM,
	token.MOD:        PRODUCT,
	token.MOD_A:      PRODUCT,
	token.ASTERISK:   PRODUCT,
	token.ASTERISK_A: PRODUCT,
	token.SLASH:      PRODUCT,
	token.SLASH_A:    PRODUCT,
	token.POWER:      PRODUCT,
	token.MATCH:      MATCHING,
	token.NOTMATCH:   MATCHING,
	token.LPAREN:     CALL,
	token.DOT:        CALL,
	token.LBRACKET:   INDEX,
	token.INCREMENT:  INCREMENT,
	token.DECREMENT:  INCREMENT,
}

type Parser struct {
	l      *lexer.Lexer
	errors []string
	path   string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer, wd string) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
		path:   wd,
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.DO, p.parseDoLoopExpression)
	p.registerPrefix(token.WHILE, p.parseWhileLoopExpression)
	p.registerPrefix(token.FOR, p.parseForLoopExpression)
	p.registerPrefix(token.GREP, p.parseGrepExpression)
	p.registerPrefix(token.MAP, p.parseMapExpression)
	p.registerPrefix(token.CASE, p.parseCaseExpression)
	p.registerPrefix(token.TRY, p.parseTryStatement)
	p.registerPrefix(token.STRING, p.parseStringLiteralExpression)
	p.registerPrefix(token.REGEX, p.parseRegExLiteralExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayExpression)
	p.registerPrefix(token.LBRACE, p.parseHashExpression)
	p.registerPrefix(token.STRUCT, p.parseStructExpression)
	p.registerPrefix(token.ISTRING, p.parseInterpolatedString)
	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)
	p.registerPrefix(token.INCREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.DECREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.YIELD, p.parseYieldExpression)
	p.registerPrefix(token.NIL, p.parseNilExpression)
	p.registerPrefix(token.ENUM, p.parseEnumExpression)
	p.registerPrefix(token.QW, p.parseQWExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.POWER, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.MATCH, p.parseInfixExpression)
	p.registerInfix(token.NOTMATCH, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LE, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.CONDAND, p.parseInfixExpression)
	p.registerInfix(token.CONDOR, p.parseInfixExpression)
	p.registerInfix(token.SHIFT_L, p.parseInfixExpression)
	p.registerInfix(token.SHIFT_R, p.parseInfixExpression)
	p.registerInfix(token.BITAND, p.parseInfixExpression)
	p.registerInfix(token.BITOR, p.parseInfixExpression)
	p.registerInfix(token.BITXOR, p.parseInfixExpression)

	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.PLUS_A, p.parseAssignExpression)
	p.registerInfix(token.MINUS_A, p.parseAssignExpression)
	p.registerInfix(token.ASTERISK_A, p.parseAssignExpression)
	p.registerInfix(token.SLASH_A, p.parseAssignExpression)
	p.registerInfix(token.MOD_A, p.parseAssignExpression)
	p.registerInfix(token.BITOR_A, p.parseAssignExpression)
	p.registerInfix(token.BITAND_A, p.parseAssignExpression)
	p.registerInfix(token.BITXOR_A, p.parseAssignExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpressions)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseMethodCallExpression)
	p.registerInfix(token.QUESTIONM, p.parseTernaryExpression)
	p.registerInfix(token.COLON, p.parseSliceExpression)
	p.registerInfix(token.INCREMENT, p.parsePostfixExpression)
	p.registerInfix(token.DECREMENT, p.parsePostfixExpression)
	p.registerInfix(token.PIPE, p.parsePipeExpression)
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	program.Includes = make(map[string]*ast.IncludeStatement)

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			if include, ok := stmt.(*ast.IncludeStatement); ok {
				_, ok := program.Includes[include.IncludePath.String()]
				if !ok {
					program.Includes[include.IncludePath.String()] = include
				}
			} else {
				program.Statements = append(program.Statements, stmt)
			}
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.DEFER:
		return p.parseDeferStatement()
	case token.SPAWN:
		return p.parseSpawnStatement()
	case token.INCLUDE:
		return p.parseIncludeStatement()
	case token.THROW:
		return p.parseThrowStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// Check if a token is ignored in expression parsing.
func (p *Parser) isIgnoredAsExpression(tok token.TokenType) bool {
	ignored := []token.TokenType{token.EOF, token.RBRACKET, token.DO, token.COMMA}
	for _, v := range ignored {
		if v == tok {
			return true
		}
	}

	return false
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// Run the infix function until the next token has
	// a higher precedence.
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{Token: p.curToken, Left: left, Operator: p.curToken.Literal}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseTryStatement() ast.Expression {

	ts := &ast.TryStmt{Token: p.curToken}
	ts.Catches = []ast.Expression{}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	ts.Block = p.parseBlockStatement().(*ast.BlockStatement)
	p.nextToken()

	for {
		if !p.curTokenIs(token.CATCH) {
			break
		}
		p.nextToken()
		savedToken := p.curToken
		if !p.curTokenIs(token.LBRACE) {
			catchStmt := &ast.CatchStmt{Token: savedToken}

			if p.curToken.Type == token.STRING {
				catchStmt.Var = p.curToken.Literal
				catchStmt.VarType = 0
			} else if p.curToken.Type == token.IDENT {
				aVar := p.parseIdentifier()
				catchStmt.Var = aVar.(*ast.Identifier).Value
				catchStmt.VarType = 1
			} else {
				return nil
			}

			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			catchStmt.Block = p.parseBlockStatement().(*ast.BlockStatement)
			ts.Catches = append(ts.Catches, catchStmt)
		} else {
			if !p.curTokenIs(token.LBRACE) {
				return nil
			}
			catchAllStmt := &ast.CatchAllStmt{Token: savedToken}
			catchAllStmt.Block = p.parseBlockStatement().(*ast.BlockStatement)
			ts.Catches = append(ts.Catches, catchAllStmt)
		}

		if !p.peekTokenIs(token.CATCH) {
			break
		}
		p.nextToken()
	} //end for

	//finally
	if p.curTokenIs(token.FINALLY) || p.peekTokenIs(token.FINALLY) {
		if p.peekTokenIs(token.FINALLY) {
			p.nextToken()
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
		} else {
			p.nextToken()
		}

		ts.Finally = p.parseBlockStatement().(*ast.BlockStatement)
	}

	if len(ts.Catches) == 0 && ts.Finally == nil { //no catch and no finally
		return nil
	}
	return ts
}

func (p *Parser) parseThrowStatement() *ast.ThrowStmt {
	stmt := &ast.ThrowStmt{Token: p.curToken}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}
	p.nextToken()
	stmt.Expr = p.parseExpressionStatement().Expression

	return stmt

}

func (p *Parser) parseStringLiteralExpression() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInterpolatedString() ast.Expression {
	is := &ast.InterpolatedString{Token: p.curToken, Value: p.curToken.Literal, ExprMap: make(map[byte]ast.Expression)}

	key := "0"[0]
	for {
		if p.curTokenIs(token.LBRACE) {
			p.nextToken()
			expr := p.parseExpression(LOWEST)
			is.ExprMap[key] = expr
			key++
		}
		p.nextInterpToken()
		if p.curTokenIs(token.ISTRING) {
			break
		}
	}

	return is
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}
	p.nextToken()
	stmt.ReturnValue = p.parseExpressionStatement().Expression

	return stmt
}

func (p *Parser) parseDeferStatement() *ast.DeferStmt {
	stmt := &ast.DeferStmt{Token: p.curToken}

	p.nextToken()
	stmt.Call = p.parseExpressionStatement().Expression

	return stmt
}

func (p *Parser) parseBreakWithoutLoopContext() ast.Expression {
	msg := fmt.Sprintf("Syntax Error: %v - 'break' outside of loop context", p.curToken.Pos)

	p.errors = append(p.errors, msg)
	return p.parseBreakExpression()
}

func (p *Parser) parseBreakExpression() ast.Expression {
	return &ast.BreakExpression{Token: p.curToken}
}

func (p *Parser) parseContinueWithoutLoopContext() ast.Expression {
	msg := fmt.Sprintf("Syntax Error: %v - 'continue' outside of loop context", p.curToken.Pos)
	p.errors = append(p.errors, msg)
	return p.parseContinueExpression()
}

func (p *Parser) parseContinueExpression() ast.Expression {
	return &ast.ContinueExpression{Token: p.curToken}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	for {
		p.nextToken()
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		stmt.Names = append(stmt.Names, name)

		if p.expectPeek(token.ASSIGN) {
			p.nextToken()
			v := p.parseExpressionStatement().Expression
			stmt.Values = append(stmt.Values, v)
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		} else {
			break
		}

	} //end for

	return stmt
}

func (p *Parser) parseBlockStatement() ast.Expression {
	expression := &ast.BlockStatement{Token: p.curToken}
	expression.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			expression.Statements = append(expression.Statements, stmt)
		}
		if p.peekTokenIs(token.EOF) {
			break
		}
		p.nextToken()
	}

	return expression
}

func (p *Parser) parseAssignExpression(name ast.Expression) ast.Expression {
	e := &ast.AssignExpression{Token: p.curToken}

	if n, ok := name.(*ast.Identifier); ok {
		e.Name = n
	} else if call, ok := name.(*ast.MethodCallExpression); ok {
		e.Name = &ast.Identifier{Token: p.curToken, Value: call.String()}
		p.nextToken()
		e.Value = p.parseExpression(LOWEST)
		return e
	} else if indexExp, ok := name.(*ast.IndexExpression); ok {
		// IndexExpression(Subscript)'s left expression should be an identifier.
		switch indexExp.Left.(type) {
		case *ast.Identifier:
			e.Name = indexExp
		default:
			msg := fmt.Sprintf("Syntax Error: Assignment operator expects an identifier")
			p.errors = append(p.errors, msg)
			return nil
		}
	} else {
		msg := fmt.Sprintf("Syntax Error: %v - expected assign token to be IDENT, got %s instead", p.curToken.Pos, name.TokenLiteral())
		p.errors = append(p.errors, msg)
		return nil
	}

	p.nextToken()
	e.Value = p.parseExpression(LOWEST)

	return e
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := &ast.IncludeStatement{Token: p.curToken}

	if p.expectPeek(token.IDENT) {
		stmt.IncludePath = p.parseExpressionStatement().Expression
	}
	program, module, err := p.getIncludedStatements(stmt.IncludePath.String())
	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
	stmt.Program = program
	stmt.IsModule = module
	return stmt
}

func (p *Parser) getIncludedStatements(importpath string) (*ast.Program, bool, error) {
	module := false
	path := p.path

	fn := path + "/" + importpath + ".my"
	f, err := ioutil.ReadFile(fn)
	if err != nil {
		path = path + "/" + importpath
		_, err := os.Stat(path)
		if err != nil {
			return nil, module, fmt.Errorf("no file or directory: %s.my, %s", importpath, path)
		}
		m, err := ioutil.ReadFile(path + "/module.my")
		if err != nil {
			return nil, module, err
		}
		module = true
		f = m
	}

	l := lexer.New(fn, string(f))
	ps := New(l, path)
	parsed := ps.ParseProgram()
	if len(ps.errors) != 0 {
		p.errors = append(p.errors, ps.errors...)
	}
	return parsed, module, nil
}

func (p *Parser) parseDoLoopExpression() ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.DoLoop{Token: p.curToken}

	p.expectPeek(token.LBRACE)
	loop.Block = p.parseBlockStatement().(*ast.BlockStatement)

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

func (p *Parser) parseWhileLoopExpression() ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.WhileLoop{Token: p.curToken}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
	}

	p.nextToken()
	loop.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	loop.Block = p.parseBlockStatement().(*ast.BlockStatement)

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

func (p *Parser) parseForLoopExpression() ast.Expression {
	curToken := p.curToken //save current token

	if p.peekTokenIs(token.LBRACE) {
		return p.parseForEverLoopExpression(curToken)
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.parseCForLoopExpression(curToken)
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	variable := p.curToken.Literal //save current identifier

	if p.peekTokenIs(token.COMMA) {
		return p.parseForEachMapExpression(curToken, variable)
	}

	ret := p.parseForEachArrayOrRangeExpression(curToken, variable)
	return ret
}

//for (init; condition; update) {}
func (p *Parser) parseCForLoopExpression(curToken token.Token) ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.ForLoop{Token: curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	loop.Init = p.parseExpression(LOWEST)

	p.nextToken()
	p.nextToken()
	loop.Cond = p.parseExpression(LOWEST)

	p.nextToken()
	p.nextToken()
	loop.Update = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	loop.Block = p.parseBlockStatement().(*ast.BlockStatement)

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

//for item in array <where cond> {}
//for item in start..end <where cond> {}
func (p *Parser) parseForEachArrayOrRangeExpression(curToken token.Token, variable string) ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	var isRange bool = false
	//loop := &ast.ForEachArrayLoop{Token: curToken, Var:variable}

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()
	aValue1 := p.parseExpression(LOWEST)

	var aValue2 ast.Expression
	if p.peekTokenIs(token.DOTDOT) {
		isRange = true
		p.nextToken()
		p.nextToken()
		aValue2 = p.parseExpression(DOTDOT)
	}

	var aCond ast.Expression
	if p.peekTokenIs(token.WHERE) {
		p.nextToken()
		p.nextToken()
		aCond = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	aBlock := p.parseBlockStatement().(*ast.BlockStatement)

	var result ast.Expression
	if !isRange {
		result = &ast.ForEachArrayLoop{Token: curToken, Var: variable, Value: aValue1, Cond: aCond, Block: aBlock}
	} else {
		result = &ast.ForEachDotRange{Token: curToken, Var: variable, StartIdx: aValue1, EndIdx: aValue2, Cond: aCond, Block: aBlock}
	}

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return result
}

//for key, value in hash {}
func (p *Parser) parseForEachMapExpression(curToken token.Token, variable string) ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.ForEachMapLoop{Token: curToken}
	loop.Key = variable

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	loop.Value = p.curToken.Literal

	if !p.expectPeek(token.IN) {
		return nil
	}

	p.nextToken()
	loop.X = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.WHERE) {
		p.nextToken()
		p.nextToken()
		loop.Cond = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	loop.Block = p.parseBlockStatement().(*ast.BlockStatement)

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

//Almost same with parseDoLoopExpression()
func (p *Parser) parseForEverLoopExpression(curToken token.Token) ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.ForEverLoop{Token: curToken}

	p.expectPeek(token.LBRACE)
	loop.Block = p.parseBlockStatement().(*ast.BlockStatement)

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	var value int64
	var err error

	if strings.HasPrefix(p.curToken.Literal, "0b") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 2, 64)
	} else if strings.HasPrefix(p.curToken.Literal, "0x") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 16, 64)
	} else if strings.HasPrefix(p.curToken.Literal, "0c") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 8, 64)
	} else {
		value, err = strconv.ParseInt(p.curToken.Literal, 10, 64)
	}

	if err != nil {
		msg := fmt.Sprintf("Syntax Error: %v - could not parse %q as integer", p.curToken.Pos, p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("Syntax Error: %v - could not parse %q as float", p.curToken.Pos, p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseRegExLiteralExpression() ast.Expression {
	return &ast.RegExLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseGrepExpression() ast.Expression {
	if p.peekTokenIs(token.LBRACE) {
		return p.parseGrepBlockExpression(p.curToken)
	}
	return p.parseGrepExprExpression(p.curToken)
}

//grep { BLOCK } LIST
func (p *Parser) parseGrepBlockExpression(tok token.Token) ast.Expression {
	gl := &ast.GrepExpr{Token: p.curToken, Var: "$_"}

	p.nextToken()
	gl.Block = p.parseBlockStatement().(*ast.BlockStatement)

	p.nextToken()
	gl.Value = p.parseExpression(LOWEST)

	return gl
}

//grep EXPR, LIST
func (p *Parser) parseGrepExprExpression(tok token.Token) ast.Expression {
	ge := &ast.GrepExpr{Token: p.curToken, Var: "$_"}

	p.nextToken()
	ge.Expr = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	p.nextToken()
	ge.Value = p.parseExpression(LOWEST)

	return ge
}

func (p *Parser) parseMapExpression() ast.Expression {
	if p.peekTokenIs(token.LBRACE) {
		return p.parseMapBlockExpression(p.curToken)
	}
	return p.parseMapExprExpression(p.curToken)
}

//map { BLOCK } LIST
func (p *Parser) parseMapBlockExpression(tok token.Token) ast.Expression {
	me := &ast.MapExpr{Token: p.curToken, Var: "$_"}

	p.nextToken()
	me.Block = p.parseBlockStatement().(*ast.BlockStatement)

	p.nextToken()
	me.Value = p.parseExpression(LOWEST)

	return me
}

//map EXPR, LIST
func (p *Parser) parseMapExprExpression(tok token.Token) ast.Expression {
	me := &ast.MapExpr{Token: p.curToken, Var: "$_"}

	p.nextToken()
	me.Expr = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	p.nextToken()
	me.Value = p.parseExpression(LOWEST)

	return me
}

//func (p *Parser) parseIfExpression() ast.Expression {
//	expression := &ast.IfExpression{Token: p.curToken}
//
//	if p.peekTokenIs(token.LPAREN) {
//		p.nextToken()
//	}
//	p.nextToken()
//	expression.Condition = p.parseExpression(LOWEST)
//
//	if p.peekTokenIs(token.RPAREN) {
//		p.nextToken()
//	}
//
//	if !p.expectPeek(token.LBRACE) {
//		return nil
//	}
//
//	expression.Consequence = p.parseBlockStatement().(*ast.BlockStatement)
//	if p.peekTokenIs(token.ELSE) {
//		p.nextToken()
//		if p.expectPeek(token.LBRACE) {
//			expression.Alternative = p.parseBlockStatement().(*ast.BlockStatement)
//		}
//	}
//
//	return expression
//}

func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{Token: p.curToken}
	// parse if/else-if expressions
	ie.Conditions = p.parseConditionalExpressions()

	// ELSE or RBRACE
	if p.peekTokenIs(token.ELSE) {
		p.nextToken() //skip "}"
		p.nextToken() //skip "else"
		ie.Alternative = p.parseBlockStatement().(*ast.BlockStatement)
	}

	return ie
}

func (p *Parser) parseConditionalExpressions() []*ast.IfConditionExpr {
	// if part
	ic := []*ast.IfConditionExpr{p.parseConditionalExpression()}

	//else-if
	for p.peekTokenIs(token.ELSEIF) || p.peekTokenIs(token.ELSIF) { //could be 'elseif' or 'elsif'
		p.nextToken()
		ic = append(ic, p.parseConditionalExpression())
	}

	return ic
}

func (p *Parser) parseConditionalExpression() *ast.IfConditionExpr {
	ic := &ast.IfConditionExpr{Token: p.curToken}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken() //skip current token
	}
	p.nextToken() //skip "{"

	ic.Cond = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken() //skip current token
	}
	p.nextToken() //skip "}"

	ic.Block = p.parseBlockStatement().(*ast.BlockStatement)

	return ic
}

func (p *Parser) parseSliceExpression(start ast.Expression) ast.Expression {
	slice := &ast.SliceExpression{Token: p.curToken, StartIndex: start}
	if p.peekTokenIs(token.RBRACKET) {
		slice.EndIndex = nil
	} else {
		p.nextToken()
		slice.EndIndex = p.parseExpression(LOWEST)
	}

	return slice
}

func (p *Parser) parseIndexExpression(arr ast.Expression) ast.Expression {
	var index ast.Expression
	indexExp := &ast.IndexExpression{Token: p.curToken, Left: arr}
	if p.peekTokenIs(token.COLON) {
		indexTok := token.Token{Type: token.INT, Literal: "0"}
		prefix := &ast.IntegerLiteral{Token: indexTok, Value: int64(0)}
		p.nextToken()
		index = p.parseSliceExpression(prefix)
	} else {
		p.nextToken()
		index = p.parseExpression(LOWEST)
	}
	indexExp.Index = index
	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
	}

	return indexExp
}

func (p *Parser) parseHashExpression() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return hash
	}
	for !p.curTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.FATARROW) {
			return nil
		}
		p.nextToken()
		hash.Pairs[key] = p.parseExpression(LOWEST)
		p.nextToken()
	}
	return hash
}

func (p *Parser) parseStructExpression() ast.Expression {
	s := &ast.StructLiteral{Token: p.curToken}
	s.Pairs = make(map[ast.Expression]ast.Expression)
	p.expectPeek(token.LBRACE)
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return s
	}
	for !p.curTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.FATARROW) {
			return nil
		}
		p.nextToken()
		s.Pairs[key] = p.parseExpression(LOWEST)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
			continue
		}
		p.nextToken()
	}
	return s
}

func (p *Parser) parseArrayExpression() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Members = p.parseExpressionArray(array.Members, token.RBRACKET)
	return array
}

// case expr in {
//    expr,expr { expr }
//    expr { expr }
//    else { expr }
// }
func (p *Parser) parseCaseExpression() ast.Expression {

	ce := &ast.CaseExpr{Token: p.curToken}
	ce.Matches = []ast.Expression{}

	p.nextToken()

	ce.Expr = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.IN) {
		ce.IsWholeMatch = false
	} else if p.peekTokenIs(token.IS) {
		ce.IsWholeMatch = true
	} else {
		return nil
	}
	p.nextToken()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.ELSE) {
			aElse := &ast.CaseElseExpr{Token: p.curToken}
			ce.Matches = append(ce.Matches, aElse)
			p.nextToken()
			aElse.Block = p.parseBlockStatement().(*ast.BlockStatement)
			p.nextToken() //skip the '}'
		} else {
			var aMatches []*ast.CaseMatchExpr
			for !p.curTokenIs(token.LBRACE) {
				aMatch := &ast.CaseMatchExpr{Token: p.curToken}
				aMatch.Expr = p.parseExpression(LOWEST)
				aMatches = append(aMatches, aMatch)

				if !p.peekTokenIs(token.COMMA) {
					p.nextToken()
					break
				}
				p.nextToken()
				p.nextToken()

			} //end for

			if !p.curTokenIs(token.LBRACE) {
				msg := fmt.Sprintf("Syntax Error: %v - expected next token to be '{', got %s instead", p.peekToken.Pos, p.curToken.Type)
				p.errors = append(p.errors, msg)
			}

			aMatchBlock := p.parseBlockStatement().(*ast.BlockStatement)
			for i := 0; i < len(aMatches); i++ {
				aMatches[i].Block = aMatchBlock
			}

			for i := 0; i < len(aMatches); i++ {
				ce.Matches = append(ce.Matches, aMatches[i])
			}

			p.nextToken() //skip the '}'
		}
	} //end for

	if !p.curTokenIs(token.RBRACE) {
		return nil
	}

	return ce
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken, Variadic: false}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.parseFuncExpressionArray(fn, token.RPAREN)

	if p.expectPeek(token.LBRACE) {
		fn.Body = p.parseBlockStatement().(*ast.BlockStatement)
	}
	return fn
}

func (p *Parser) parseFuncExpressionArray(fn *ast.FunctionLiteral, closure token.TokenType) {
	if p.peekTokenIs(closure) {
		p.nextToken()
		return
	}

	var hasDefParamValue bool = false
	for {
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			msg := fmt.Sprintf("Syntax Error: Function parameter not identifier, GOT(%s)!", p.curToken.Literal)
			p.errors = append(p.errors, msg)
			return
		}
		key := p.curToken.Literal
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		fn.Parameters = append(fn.Parameters, name)

		if p.peekTokenIs(token.ASSIGN) {
			hasDefParamValue = true
			p.nextToken()
			p.nextToken()
			v := p.parseExpressionStatement().Expression

			if fn.Values == nil {
				fn.Values = make(map[string]ast.Expression)
			}
			fn.Values[key] = v
		} else if !p.peekTokenIs(token.ELLIPSIS) {
			if hasDefParamValue && !fn.Variadic {
				msg := fmt.Sprintf("Syntax Error: Function's default parameter order not correct!")
				p.errors = append(p.errors, msg)
				return
			}
		}

		if p.peekTokenIs(token.COMMA) {
			if fn.Variadic {
				msg := fmt.Sprintf("Syntax Error: Variadic argument in function should be last!")
				p.errors = append(p.errors, msg)
				return
			}
			p.nextToken()
		}

		if p.peekTokenIs(token.ELLIPSIS) { //Variadic function
			if fn.Variadic {
				msg := fmt.Sprintf("Syntax Error: Only 1 variadic argument is allowed in function!")
				p.errors = append(p.errors, msg)
				return
			}
			fn.Variadic = true

			p.nextToken()
			if !p.peekTokenIs(closure) {
				msg := fmt.Sprintf("Syntax Error: Variadic argument in function should be last!")
				p.errors = append(p.errors, msg)
				return
			}
		}

		if p.peekTokenIs(closure) {
			p.nextToken()
			break
		}
	}

	return
}

func (p *Parser) parseCallExpressions(f ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.curToken, Function: f}
	call.Arguments = p.parseExpressionArray(call.Arguments, token.RPAREN)
	return call
}

func (p *Parser) parseExpressionArray(a []ast.Expression, closure token.TokenType) []ast.Expression {
	if p.peekTokenIs(closure) {
		p.nextToken()
		return a
	}
	p.nextToken()
	a = append(a, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		a = append(a, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(closure) {
		return nil
	}
	return a
}

func (p *Parser) parseMethodCallExpression(obj ast.Expression) ast.Expression {
	methodCall := &ast.MethodCallExpression{Token: p.curToken, Object: obj}
	p.nextToken()

	name := p.parseIdentifier()
	if !p.peekTokenIs(token.LPAREN) {
		//methodCall.Call = p.parseExpression(LOWEST)
		//Note: here the precedence should not be `LOWEST`, or else when parsing below line:
		//     logger.LDATE + 1 ==> logger.(LDATE + 1)
		methodCall.Call = p.parseExpression(CALL)
	} else {
		p.nextToken()
		methodCall.Call = p.parseCallExpressions(name)
	}

	return methodCall
}

func (p *Parser) parseTernaryExpression(condition ast.Expression) ast.Expression {
	expression := &ast.TernaryExpression{
		Token:     p.curToken,
		Condition: condition,
	}
	precedence := p.curPrecedence()
	p.nextToken() //skip the '?'
	expression.IfTrue = p.parseExpression(precedence)

	if !p.expectPeek(token.COLON) { //skip the ":"
		return nil
	}

	// Get to next token, then parse the else part
	p.nextToken()
	expression.IfFalse = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseSpawnStatement() *ast.SpawnStmt {
	stmt := &ast.SpawnStmt{Token: p.curToken}

	p.nextToken()
	stmt.Call = p.parseExpressionStatement().Expression

	return stmt
}

//NOT IMPLEMENTED CORRECTLY
func (p *Parser) parseYieldExpression() ast.Expression {
	yield := &ast.YieldExpression{Token: p.curToken}

	p.nextToken()
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		//yield.Arguments =
		p.nextToken()
		return yield
	}
	//yield.Arguments =
	return yield
}

func (p *Parser) parseNilExpression() ast.Expression {
	return &ast.NilLiteral{Token: p.curToken}
}

func (p *Parser) parseEnumExpression() ast.Expression {
	var autoInt int64 = 0 //autoIncrement

	e := &ast.EnumLiteral{Token: p.curToken}
	e.Pairs = make(map[ast.Expression]ast.Expression)

	if !p.expectPeek(token.LBRACE) {
		return e
	}

	for {
		//check for empty `enum`
		if p.peekTokenIs(token.RBRACE) {
			p.nextToken()
			return e
		}

		// identifier is mandatory here
		if !p.expectPeek(token.IDENT) {
			return e
		}
		enum_id := p.parseIdentifier()

		// peek next that can be only '=' or ',' or '}'
		if !p.peekTokenIs(token.ASSIGN) && !p.peekTokenIs(token.COMMA) && !p.peekTokenIs(token.RBRACE) {
			msg := fmt.Sprintf("Syntax Error: %v - Token %s not allowed here.", p.peekToken.Pos, p.peekToken.Type)
			p.errors = append(p.errors, msg)
		}

		// check for optional default value (optional only in `INT` case)
		var enum_value ast.Expression
		if p.peekTokenIs(token.ASSIGN) {
			p.nextToken()
			p.nextToken()

			enum_value = p.parseExpressionStatement().Expression
		}

		if enum_value != nil {
			if _, ok := enum_value.(*ast.IntegerLiteral); ok {
				intLiteral := enum_value.(*ast.IntegerLiteral)
				autoInt = intLiteral.Value + 1
			}
		} else {
			//create a new INT token with 'autoInt' as it's value
			tok := token.Token{Type: token.INT, Literal: strconv.Itoa(int(autoInt))}
			enum_value = &ast.IntegerLiteral{Token: tok, Value: autoInt}
			autoInt++
		}

		if _, ok := e.Pairs[enum_id]; ok { //is identifier redeclared
			msg := fmt.Sprintf("Syntax Error: %v - Identifier %s redeclared.", p.peekToken.Pos, p.curToken.Type)
			p.errors = append(p.errors, msg)
		} else {
			e.Pairs[enum_id] = enum_value
		}

		if !p.peekTokenIs(token.COMMA) {
			p.nextToken()
			break
		}
		p.nextToken()
	}

	return e
}

//qw(xx, xx, xx, xx)
func (p *Parser) parseQWExpression() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Members = p.parseStrExpressionArray(array.Members)
	return array
}

func (p *Parser) parseStrExpressionArray(a []ast.Expression) []ast.Expression {
	var allowedPair = map[string]token.TokenType{
		"(": token.RPAREN,
		"<": token.GT,
		"{": token.RBRACE,
	}

	p.nextToken() //skip the 'qw'
	openPair := p.curToken.Literal
	if p.curTokenIs(allowedPair[openPair]) {
		p.nextToken()
		return a
	}

	p.nextToken()
	if !p.curTokenIs(token.IDENT) && !p.curTokenIs(token.INT) && !p.curTokenIs(token.FLOAT) {
		msg := fmt.Sprintf("Syntax Error: %v - expected next token to be 'IDENT|INT|FLOAT', got %s instead", p.peekToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	a = append(a, &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		if !p.curTokenIs(token.IDENT) && !p.curTokenIs(token.INT) && !p.curTokenIs(token.FLOAT) {
			msg := fmt.Sprintf("Syntax Error: %v - expected next token to be 'IDENT|INT|FLOAT', got %s instead", p.peekToken.Pos, p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
		a = append(a, &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal})
	}
	if !p.expectPeek(allowedPair[openPair]) {
		return nil
	}
	return a
}

// IDENT() -> IDENT()
func (p *Parser) parsePipeExpression(left ast.Expression) ast.Expression {
	expression := &ast.Pipe{
		Token: p.curToken,
		Left:  left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("Syntax Error: %v - no prefix parse functions for '%s' found", p.curToken.Pos, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) nextInterpToken() {
	p.curToken = p.l.NextInterpToken()
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Syntax Error: %v - expected next token to be %s, got %s instead", p.peekToken.Pos, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}
