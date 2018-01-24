package ast

import (
	"bytes"
	"monkey/token"
	"strings"
	"unicode/utf8"
)

type Node interface {
	Pos() token.Position // position of first character belonging to the node
	End() token.Position // position of first character immediately after the node

	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
	Includes   map[string]*IncludeStatement
}

func (p *Program) Pos() token.Position {
	if len(p.Statements) > 0 {
		return p.Statements[0].Pos()
	}
	return token.Position{}
}

func (p *Program) End() token.Position {
	aLen := len(p.Statements)
	if aLen > 0 {
		return p.Statements[aLen-1].End()
	}
	return token.Position{}
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) Pos() token.Position {
	return bs.Token.Pos

}

func (bs *BlockStatement) End() token.Position {
	aLen := len(bs.Statements)
	if aLen > 0 {
		return bs.Statements[aLen-1].End()
	}
	return bs.Token.Pos
}

func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

///////////////////////////////////////////////////////////
//                        FOR LOOP                       //
///////////////////////////////////////////////////////////
type ForLoop struct {
	Token  token.Token
	Init   Expression
	Cond   Expression
	Update Expression
	Block  *BlockStatement
}

func (fl *ForLoop) Pos() token.Position {
	return fl.Token.Pos
}

func (fl *ForLoop) End() token.Position {
	return fl.Block.End()
}

func (fl *ForLoop) expressionNode()      {}
func (fl *ForLoop) TokenLiteral() string { return fl.Token.Literal }

func (fl *ForLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for")
	out.WriteString(" ( ")
	out.WriteString(fl.Init.String())
	out.WriteString(" ; ")
	out.WriteString(fl.Cond.String())
	out.WriteString(" ; ")
	out.WriteString(fl.Update.String())
	out.WriteString(" ) ")
	out.WriteString(" { ")
	out.WriteString(fl.Block.String())
	out.WriteString(" }")

	return out.String()
}

type ForEachArrayLoop struct {
	Token token.Token
	Var   string
	Value Expression //value to range over
	Cond  Expression //conditional clause(nil if there is no 'WHERE' clause)
	Block *BlockStatement
}

func (fal *ForEachArrayLoop) Pos() token.Position {
	return fal.Token.Pos
}

func (fal *ForEachArrayLoop) End() token.Position {
	return fal.Block.End()
}

func (fal *ForEachArrayLoop) expressionNode()      {}
func (fal *ForEachArrayLoop) TokenLiteral() string { return fal.Token.Literal }

func (fal *ForEachArrayLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fal.Var)
	out.WriteString(" in ")
	out.WriteString(fal.Value.String())
	if fal.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(fal.Cond.String())
	}
	out.WriteString(" { ")
	out.WriteString(fal.Block.String())
	out.WriteString(" }")

	return out.String()
}

type ForEachMapLoop struct {
	Token token.Token
	Key   string
	Value string
	X     Expression //value to range over
	Cond  Expression //Conditional clause(nil if there is no 'WHERE' clause)
	Block *BlockStatement
}

func (fml *ForEachMapLoop) Pos() token.Position {
	return fml.Token.Pos
}

func (fml *ForEachMapLoop) End() token.Position {
	return fml.Block.End()
}

func (fml *ForEachMapLoop) expressionNode()      {}
func (fml *ForEachMapLoop) TokenLiteral() string { return fml.Token.Literal }

func (fml *ForEachMapLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fml.Key + ", " + fml.Value)
	out.WriteString(" in ")
	out.WriteString(fml.X.String())
	if fml.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(fml.Cond.String())
	}
	out.WriteString(" { ")
	out.WriteString(fml.Block.String())
	out.WriteString(" }")

	return out.String()
}

type ForEverLoop struct {
	Token token.Token
	Block *BlockStatement
}

func (fel *ForEverLoop) Pos() token.Position {
	return fel.Token.Pos
}

func (fel *ForEverLoop) End() token.Position {
	return fel.Block.End()
}

func (fel *ForEverLoop) expressionNode()      {}
func (fel *ForEverLoop) TokenLiteral() string { return fel.Token.Literal }

func (fel *ForEverLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(" { ")
	out.WriteString(fel.Block.String())
	out.WriteString(" }")

	return out.String()
}

//for i in start..end <where cond> { }
type ForEachDotRange struct {
	Token    token.Token
	Var      string
	StartIdx Expression
	EndIdx   Expression
	Cond     Expression //conditional clause(nil if there is no 'WHERE' clause)
	Block    *BlockStatement
}

func (fdr *ForEachDotRange) Pos() token.Position {
	return fdr.Token.Pos
}

func (fdr *ForEachDotRange) End() token.Position {
	return fdr.Block.End()
}

func (fdr *ForEachDotRange) expressionNode()      {}
func (fdr *ForEachDotRange) TokenLiteral() string { return fdr.Token.Literal }

func (fdr *ForEachDotRange) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fdr.Var)
	out.WriteString(" in ")
	out.WriteString(fdr.StartIdx.String())
	out.WriteString(" .. ")
	out.WriteString(fdr.EndIdx.String())
	if fdr.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(fdr.Cond.String())
	}
	out.WriteString(" { ")
	out.WriteString(fdr.Block.String())
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                        WHILE LOOP                     //
///////////////////////////////////////////////////////////
type WhileLoop struct {
	Token     token.Token
	Condition Expression
	Block     *BlockStatement
}

func (wl *WhileLoop) Pos() token.Position {
	return wl.Token.Pos
}

func (wl *WhileLoop) End() token.Position {
	return wl.Block.End()
}

func (wl *WhileLoop) expressionNode()      {}
func (wl *WhileLoop) TokenLiteral() string { return wl.Token.Literal }

func (wl *WhileLoop) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(wl.Condition.String())
	out.WriteString("{")
	out.WriteString(wl.Block.String())
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         DO LOOP                       //
///////////////////////////////////////////////////////////
type DoLoop struct {
	Token token.Token
	Block *BlockStatement
}

func (dl *DoLoop) Pos() token.Position {
	return dl.Token.Pos
}

func (dl *DoLoop) End() token.Position {
	return dl.Block.End()
}

func (dl *DoLoop) expressionNode()      {}
func (dl *DoLoop) TokenLiteral() string { return dl.Token.Literal }

func (dl *DoLoop) String() string {
	var out bytes.Buffer

	out.WriteString("do")
	out.WriteString(" { ")
	out.WriteString(dl.Block.String())
	out.WriteString(" }")
	return out.String()
}

///////////////////////////////////////////////////////////
//                        IDENTIFIER                     //
///////////////////////////////////////////////////////////
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) Pos() token.Position {
	return i.Token.Pos
}

func (i *Identifier) End() token.Position {
	length := utf8.RuneCountInString(i.Value)
	return token.Position{Line: i.Token.Pos.Line, Col: i.Token.Pos.Col + length - 1}
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

///////////////////////////////////////////////////////////
//                         IFELSE                        //
///////////////////////////////////////////////////////////
//type IfExpression struct {
//	Token       token.Token
//	Condition   Expression
//	Consequence *BlockStatement
//	Alternative *BlockStatement
//}
//
//func (ifex *IfExpression) Pos() token.Position {
//	return ifex.Token.Pos
//}
//
//func (ifex *IfExpression) End() token.Position {
//	if ifex.Alternative != nil {
//		return ifex.Alternative.End()
//	}
//	return ifex.Consequence.End()
//}
//
//func (ifex *IfExpression) expressionNode()      {}
//func (ifex *IfExpression) TokenLiteral() string { return ifex.Token.Literal }
//
//func (ifex *IfExpression) String() string {
//	var out bytes.Buffer
//
//	out.WriteString("if ")
//	out.WriteString("(")
//	out.WriteString(ifex.Condition.String())
//	out.WriteString(")")
//	out.WriteString(" { ")
//	out.WriteString(ifex.Consequence.String())
//	out.WriteString(" }")
//	if ifex.Alternative != nil {
//		out.WriteString(" else ")
//		out.WriteString(" { ")
//		out.WriteString(ifex.Alternative.String())
//		out.WriteString(" }")
//	}
//
//	return out.String()
//}

type IfExpression struct {
	Token       token.Token
	Conditions  []*IfConditionExpr //if or elseif part
	Alternative *BlockStatement    //else part
}

func (ifex *IfExpression) Pos() token.Position {
	return ifex.Token.Pos
}

func (ifex *IfExpression) End() token.Position {
	if ifex.Alternative != nil {
		return ifex.Alternative.End()
	}

	aLen := len(ifex.Conditions)
	return ifex.Conditions[aLen-1].End()
}

func (ifex *IfExpression) expressionNode()      {}
func (ifex *IfExpression) TokenLiteral() string { return ifex.Token.Literal }

func (ifex *IfExpression) String() string {
	var out bytes.Buffer

	for i, c := range ifex.Conditions {
		if i == 0 {
			out.WriteString("if ")
		} else {
			out.WriteString("elseif ")
		}
		out.WriteString(c.String())
	}

	if ifex.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(" { ")
		out.WriteString(ifex.Alternative.String())
		out.WriteString(" }")
	}

	return out.String()
}

//if/else-if condition
type IfConditionExpr struct {
	Token token.Token
	Cond  Expression      //condition
	Block *BlockStatement //body
}

func (ic *IfConditionExpr) Pos() token.Position {
	return ic.Token.Pos
}

func (ic *IfConditionExpr) End() token.Position {
	return ic.Block.End()
}

func (ic *IfConditionExpr) expressionNode()      {}
func (ic *IfConditionExpr) TokenLiteral() string { return ic.Token.Literal }

func (ic *IfConditionExpr) String() string {
	var out bytes.Buffer

	out.WriteString(ic.Cond.String())
	out.WriteString(" { ")
	out.WriteString(ic.Block.String())
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                    UNLESS-ELSE                        //
///////////////////////////////////////////////////////////
type UnlessExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ul *UnlessExpression) Pos() token.Position {
	return ul.Token.Pos
}

func (ul *UnlessExpression) End() token.Position {
	if ul.Alternative != nil {
		return ul.Alternative.End()
	}
	return ul.Consequence.End()
}

func (ul *UnlessExpression) expressionNode()      {}
func (ul *UnlessExpression) TokenLiteral() string { return ul.Token.Literal }

func (ul *UnlessExpression) String() string {
	var out bytes.Buffer

	out.WriteString("unless ")
	out.WriteString("(")
	out.WriteString(ul.Condition.String())
	out.WriteString(")")
	out.WriteString(" { ")
	out.WriteString(ul.Consequence.String())
	out.WriteString(" }")
	if ul.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(" { ")
		out.WriteString(ul.Alternative.String())
		out.WriteString(" }")
	}

	return out.String()
}

///////////////////////////////////////////////////////////
//                         HASH LITERAL                  //
///////////////////////////////////////////////////////////
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) Pos() token.Position {
	return h.Token.Pos
}

func (h *HashLiteral) End() token.Position {
	maxLineMap := make(map[int]Expression)

	for _, value := range h.Pairs {
		v := value.(Expression)
		maxLineMap[v.End().Line] = v
	}

	maxLine := 0
	for line, _ := range maxLineMap {
		if line > maxLine {
			maxLine = line
		}
	}

	ret := maxLineMap[maxLine].(Expression)
	return ret.End()
}

func (h *HashLiteral) expressionNode()      {}
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String()+"=>"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     NIL LITERAL                   //
///////////////////////////////////////////////////////////
type NilLiteral struct {
	Token token.Token
}

func (n *NilLiteral) Pos() token.Position {
	return n.Token.Pos
}

func (n *NilLiteral) End() token.Position {
	length := len(n.Token.Literal)
	pos := n.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length - 1}
}

func (n *NilLiteral) expressionNode()      {}
func (n *NilLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NilLiteral) String() string       { return n.Token.Literal }

///////////////////////////////////////////////////////////
//                     INTEGER LITERAL                   //
///////////////////////////////////////////////////////////
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) Pos() token.Position {
	return il.Token.Pos
}

func (il *IntegerLiteral) End() token.Position {
	length := utf8.RuneCountInString(il.Token.Literal)
	pos := il.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length - 1}
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

///////////////////////////////////////////////////////////
//               UNSIGNED INTEGER LITERAL                //
///////////////////////////////////////////////////////////
type UIntegerLiteral struct { //U: Unsigned
	Token token.Token
	Value uint64
}

func (il *UIntegerLiteral) Pos() token.Position {
	return il.Token.Pos
}

func (il *UIntegerLiteral) End() token.Position {
	length := utf8.RuneCountInString(il.Token.Literal)
	pos := il.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length - 1}
}

func (il *UIntegerLiteral) expressionNode()      {}
func (il *UIntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *UIntegerLiteral) String() string       { return il.Token.Literal }

///////////////////////////////////////////////////////////
//                     FLOAT LITERAL                     //
///////////////////////////////////////////////////////////
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) Pos() token.Position {
	return fl.Token.Pos
}

func (fl *FloatLiteral) End() token.Position {
	length := utf8.RuneCountInString(fl.Token.Literal)
	pos := fl.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length - 1}
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

///////////////////////////////////////////////////////////
//                     BOOLEAN LITERAL                   //
///////////////////////////////////////////////////////////
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) Pos() token.Position {
	return b.Token.Pos
}

func (b *Boolean) End() token.Position {
	length := utf8.RuneCountInString(b.Token.Literal)
	pos := b.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length - 1}
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

///////////////////////////////////////////////////////////
//                      REGEX LITERAL                    //
///////////////////////////////////////////////////////////
type RegExLiteral struct {
	Token token.Token
	Value string
}

func (rel *RegExLiteral) Pos() token.Position {
	return rel.Token.Pos
}

func (rel *RegExLiteral) End() token.Position {
	return rel.Token.Pos
}

func (rel *RegExLiteral) expressionNode()      {}
func (rel *RegExLiteral) TokenLiteral() string { return rel.Token.Literal }
func (rel *RegExLiteral) String() string       { return rel.Value }

///////////////////////////////////////////////////////////
//                      ARRAY LITERAL                    //
///////////////////////////////////////////////////////////
type ArrayLiteral struct {
	Token   token.Token
	Members []Expression
}

func (a *ArrayLiteral) Pos() token.Position {
	return a.Token.Pos
}

func (a *ArrayLiteral) End() token.Position {
	aLen := len(a.Members)
	if aLen > 0 {
		return a.Members[aLen-1].End()
	}
	return a.Token.Pos
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	members := []string{}
	for _, m := range a.Members {
		members = append(members, m.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(members, ", "))
	out.WriteString("]")
	return out.String()
}

///////////////////////////////////////////////////////////
//                      RANGE LITERAL(..)                //
///////////////////////////////////////////////////////////
//type RangeLiteral struct {
//	Token    token.Token
//	StartIdx Expression
//	EndIdx   Expression
//}
//
//func (r *RangeLiteral) expressionNode()      {}
//func (r *RangeLiteral) TokenLiteral() string { return r.Token.Literal }
//func (r *RangeLiteral) String() string {
//	var out bytes.Buffer
//
//	out.WriteString("(")
//	out.WriteString(r.StartIdx.String())
//	out.WriteString(" .. ")
//	out.WriteString(r.EndIdx.String())
//	out.WriteString(")")
//
//	return out.String()
//}

///////////////////////////////////////////////////////////
//                     FUNCTION LITERAL                  //
///////////////////////////////////////////////////////////
type FunctionLiteral struct {
	Token      token.Token
	Parameters []Expression
	Body       *BlockStatement

	//Default values
	Values map[string]Expression

	Variadic bool

	ModifierLevel ModifierLevel //for 'class' use
}

func (fl *FunctionLiteral) Pos() token.Position {
	return fl.Token.Pos
}

func (fl *FunctionLiteral) End() token.Position {
	return fl.Body.End()
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(fl.TokenLiteral())
	params := []string{}
	for i, p := range fl.Parameters {
		param := p.String()
		if fl.Variadic && i == len(fl.Parameters)-1 {
			param = "..." + param
		}

		params = append(params, p.String())

	}
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString("{ ")
	out.WriteString(fl.Body.String())
	out.WriteString(" }")
	return out.String()
}

///////////////////////////////////////////////////////////
//                  FUNCTION STATEMENT                   //
///////////////////////////////////////////////////////////
type FunctionStatement struct {
	Token           token.Token
	Name            *Identifier
	FunctionLiteral *FunctionLiteral
}

func (f *FunctionStatement) Pos() token.Position {
	return f.Token.Pos
}

func (f *FunctionStatement) End() token.Position {
	return f.FunctionLiteral.Body.End()
}

func (f *FunctionStatement) statementNode() {}
func (f *FunctionStatement) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionStatement) String() string {
	var out bytes.Buffer

	out.WriteString(f.FunctionLiteral.ModifierLevel.String())

	out.WriteString("fn ")
	out.WriteString(f.Name.String())

	params := []string{}
	for i, p := range f.FunctionLiteral.Parameters {
		param := p.String()
		if f.FunctionLiteral.Variadic && i == len(f.FunctionLiteral.Parameters)-1 {
			param = "..." + param
		}

		params = append(params, p.String())

	}
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString("{ ")
	out.WriteString(f.FunctionLiteral.Body.String())
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                      STRING LITERAL                   //
///////////////////////////////////////////////////////////
type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) Pos() token.Position {
	return s.Token.Pos
}

func (s *StringLiteral) End() token.Position {
	return s.Token.Pos
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return s.Token.Literal }

///////////////////////////////////////////////////////////
//                  INTERPOLATED STRING                  //
///////////////////////////////////////////////////////////
type InterpolatedString struct {
	Token   token.Token
	Value   string
	ExprMap map[byte]Expression
}

func (is *InterpolatedString) Pos() token.Position {
	return is.Token.Pos
}

func (is *InterpolatedString) End() token.Position {
	return is.Token.Pos
}

func (is *InterpolatedString) expressionNode()      {}
func (is *InterpolatedString) TokenLiteral() string { return is.Token.Literal }
func (is *InterpolatedString) String() string       { return is.Token.Literal }

///////////////////////////////////////////////////////////
//                    TRY/CATCH/FINALLY                  //
///////////////////////////////////////////////////////////
//TryStmt provide "try/catch/finally" statement.
type TryStmt struct {
	Token   token.Token
	Block   *BlockStatement
	Catches []Expression //catch
	Finally *BlockStatement
}

func (t *TryStmt) Pos() token.Position {
	return t.Token.Pos
}

func (t *TryStmt) End() token.Position {
	if t.Finally != nil {
		t.Finally.End()
	}

	aLen := len(t.Catches)
	if aLen > 0 {
		return t.Catches[aLen-1].End()
	}
	return t.Token.Pos
}

func (t *TryStmt) expressionNode()      {}
func (t *TryStmt) TokenLiteral() string { return t.Token.Literal }

func (t *TryStmt) String() string {
	var out bytes.Buffer

	out.WriteString("try")
	out.WriteString(" { ")
	out.WriteString(t.Block.String())
	out.WriteString(" } ")

	for _, o := range t.Catches {
		out.WriteString(o.String())
	}

	if t.Finally != nil {
		out.WriteString("finally")
		out.WriteString(" { ")
		out.WriteString(t.Finally.String())
		out.WriteString(" }")
	}

	return out.String()
}

type CatchStmt struct {
	Token   token.Token
	Var     string //maybe nil
	VarType int    // 0:STRING, 1:IDENTIFIER
	Block   *BlockStatement
}

func (c *CatchStmt) Pos() token.Position {
	return c.Token.Pos
}

func (c *CatchStmt) End() token.Position {
	return c.Block.End()
}

func (c *CatchStmt) expressionNode()      {}
func (c *CatchStmt) TokenLiteral() string { return c.Token.Literal }

func (c *CatchStmt) String() string {
	var out bytes.Buffer

	out.WriteString("catch ")

	if len(c.Var) > 0 {
		out.WriteString(c.Var)
	}

	out.WriteString(" { ")
	out.WriteString(c.Block.String())
	out.WriteString(" }")

	return out.String()
}

type CatchAllStmt struct {
	Token token.Token
	Block *BlockStatement
}

func (ca *CatchAllStmt) Pos() token.Position {
	return ca.Token.Pos
}

func (ca *CatchAllStmt) End() token.Position {
	return ca.Block.End()
}

func (ca *CatchAllStmt) expressionNode()      {}
func (ca *CatchAllStmt) TokenLiteral() string { return ca.Token.Literal }

func (ca *CatchAllStmt) String() string {
	var out bytes.Buffer

	out.WriteString("catch ")
	out.WriteString(" { ")
	out.WriteString(ca.Block.String())
	out.WriteString(" }")

	return out.String()
}

//ThrowStmt provide "throw" expression statement.
//Note: only support throwing "String" object
type ThrowStmt struct {
	Token token.Token
	Expr  Expression
}

func (ts *ThrowStmt) Pos() token.Position {
	return ts.Token.Pos
}

func (ts *ThrowStmt) End() token.Position {
	return ts.Expr.End()
}

func (ts *ThrowStmt) statementNode()       {}
func (ts *ThrowStmt) TokenLiteral() string { return ts.Token.Literal }

func (ts *ThrowStmt) String() string {
	var out bytes.Buffer

	out.WriteString("throw ")
	out.WriteString(ts.Expr.String())
	out.WriteString(";")

	return out.String()
}

///////////////////////////////////////////////////////////
//                      STRUCT LITERAL                   //
///////////////////////////////////////////////////////////
type StructLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (s *StructLiteral) Pos() token.Position {
	return s.Token.Pos
}

func (s *StructLiteral) End() token.Position {
	maxLineMap := make(map[int]Expression)

	for _, value := range s.Pairs {
		v := value.(Expression)
		maxLineMap[v.End().Line] = v
	}

	maxLine := 0
	for line, _ := range maxLineMap {
		if line > maxLine {
			maxLine = line
		}
	}

	ret := maxLineMap[maxLine].(Expression)
	return ret.End()
}

func (s *StructLiteral) expressionNode()      {}
func (s *StructLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StructLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range s.Pairs {
		pairs = append(pairs, key.String()+"=>"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     DEFER STATEMENT                   //
///////////////////////////////////////////////////////////
type DeferStmt struct {
	Token token.Token
	Call  Expression
}

func (ds *DeferStmt) Pos() token.Position {
	return ds.Token.Pos
}

func (ds *DeferStmt) End() token.Position {
	return ds.Call.End()
}

func (ds *DeferStmt) statementNode()       {}
func (ds *DeferStmt) TokenLiteral() string { return ds.Token.Literal }

func (ds *DeferStmt) String() string {
	var out bytes.Buffer

	out.WriteString(ds.TokenLiteral() + " ")
	out.WriteString(ds.Call.String())
	out.WriteString("; ")

	return out.String()
}

///////////////////////////////////////////////////////////
//                    RETURN STATEMENT                   //
///////////////////////////////////////////////////////////
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) Pos() token.Position {
	return rs.Token.Pos
}

func (rs *ReturnStatement) End() token.Position {
	return rs.ReturnValue.End()
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString("; ")

	return out.String()
}

///////////////////////////////////////////////////////////
//                  EXPRESSION STATEMENT                 //
///////////////////////////////////////////////////////////
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) Pos() token.Position {
	return es.Token.Pos
}

func (es *ExpressionStatement) End() token.Position {
	return es.Expression.End()
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

///////////////////////////////////////////////////////////
//                      LET STATEMENT                    //
///////////////////////////////////////////////////////////
type LetStatement struct {
	Token  token.Token
	Names  []*Identifier
	Values []Expression

	ModifierLevel ModifierLevel //used in 'class'
}

func (ls *LetStatement) Pos() token.Position {
	return ls.Token.Pos
}

func (ls *LetStatement) End() token.Position {
	aLen := len(ls.Values)
	if aLen > 0 {
		return ls.Values[aLen-1].End()
	}
	return ls.Token.Pos
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.ModifierLevel.String())

	valuesLen := len(ls.Values)

	out.WriteString(ls.TokenLiteral() + " ")
	for idx, item := range ls.Names {
		out.WriteString(item.String())
		out.WriteString(" = ")

		if idx >= valuesLen {
			out.WriteString("")
		} else {
			if ls.Values[idx] != nil {
				out.WriteString(ls.Values[idx].String())
			}
		}
		if idx != len(ls.Names)-1 {
			out.WriteString(", ")
		}
	}

	if len(ls.Names) == 1 {
		out.WriteString("; ")
	}

	return out.String()
}

///////////////////////////////////////////////////////////
//                      INCLUDE STATEMENT                //
///////////////////////////////////////////////////////////
type IncludeStatement struct {
	Token       token.Token
	IncludePath Expression
	Program     *Program
}

func (is *IncludeStatement) Pos() token.Position {
	return is.Token.Pos
}

func (is *IncludeStatement) End() token.Position {
	return is.IncludePath.End()
}

func (is *IncludeStatement) statementNode()       {}
func (is *IncludeStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncludeStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(is.IncludePath.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                         BREAK                         //
///////////////////////////////////////////////////////////
type BreakExpression struct {
	Token token.Token
}

func (be *BreakExpression) Pos() token.Position {
	return be.Token.Pos
}

func (be *BreakExpression) End() token.Position {
	length := utf8.RuneCountInString(be.Token.Literal)
	pos := be.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length - 1}
}

func (be *BreakExpression) expressionNode()      {}
func (be *BreakExpression) TokenLiteral() string { return be.Token.Literal }

func (be *BreakExpression) String() string { return be.Token.Literal }

///////////////////////////////////////////////////////////
//                         CONTINUE                      //
///////////////////////////////////////////////////////////
type ContinueExpression struct {
	Token token.Token
}

func (ce *ContinueExpression) Pos() token.Position {
	return ce.Token.Pos
}

func (ce *ContinueExpression) End() token.Position {
	length := utf8.RuneCountInString(ce.Token.Literal)
	pos := ce.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length - 1}
}

func (ce *ContinueExpression) expressionNode()      {}
func (ce *ContinueExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *ContinueExpression) String() string { return ce.Token.Literal }

///////////////////////////////////////////////////////////
//                         ASSIGN                        //
///////////////////////////////////////////////////////////
type AssignExpression struct {
	Token token.Token
	Name  Expression
	Value Expression
}

func (ae *AssignExpression) Pos() token.Position {
	return ae.Token.Pos
}

func (ae *AssignExpression) End() token.Position {
	return ae.Value.End()
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }

func (ae *AssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Name.String())
	//out.WriteString(" = ")
	out.WriteString(ae.Token.Literal)
	out.WriteString(ae.Value.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                         GREP                          //
///////////////////////////////////////////////////////////
type GrepExpr struct {
	Token token.Token
	Var   string          //Name is "$_"
	Value Expression      //value to range over
	Block *BlockStatement //Grep Block, may be nil
	Expr  Expression      //Grep Expr, may be nil
}

func (ge *GrepExpr) Pos() token.Position {
	return ge.Token.Pos
}

func (ge *GrepExpr) End() token.Position {
	if ge.Block == nil {
		return ge.Expr.End()
	}
	if ge.Expr == nil {
		return ge.Block.End()
	}
	return ge.Token.Pos
}

func (ge *GrepExpr) expressionNode()      {}
func (ge *GrepExpr) TokenLiteral() string { return ge.Token.Literal }

func (ge *GrepExpr) String() string {
	var out bytes.Buffer

	out.WriteString("grep ")
	if ge.Block != nil {
		out.WriteString(" { ")
		out.WriteString(ge.Block.String())
		out.WriteString(" } ")
	} else {
		out.WriteString(ge.Expr.String())
		out.WriteString(" , ")
	}

	out.WriteString(ge.Value.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                         MAP                           //
///////////////////////////////////////////////////////////
type MapExpr struct {
	Token token.Token
	Var   string          //Name is "$_"
	Value Expression      //value to range over
	Block *BlockStatement //Grep Block, may be nil
	Expr  Expression      //Grep Expr, may be nil
}

func (me *MapExpr) Pos() token.Position {
	return me.Token.Pos
}

func (me *MapExpr) End() token.Position {
	if me.Block == nil {
		return me.Expr.End()
	}
	if me.Expr == nil {
		return me.Block.End()
	}
	return me.Token.Pos
}

func (me *MapExpr) expressionNode()      {}
func (me *MapExpr) TokenLiteral() string { return me.Token.Literal }

func (me *MapExpr) String() string {
	var out bytes.Buffer

	out.WriteString("map ")
	if me.Block != nil {
		out.WriteString(" { ")
		out.WriteString(me.Block.String())
		out.WriteString(" } ")
	} else {
		out.WriteString(me.Expr.String())
		out.WriteString(" , ")
	}

	out.WriteString(me.Value.String())
	return out.String()

}

///////////////////////////////////////////////////////////
//                         INFIX                         //
///////////////////////////////////////////////////////////
type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) Pos() token.Position {
	return ie.Token.Pos
}

func (ie *InfixExpression) End() token.Position {
	return ie.Right.End()
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         PREFIX                        //
///////////////////////////////////////////////////////////
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) Pos() token.Position {
	return pe.Token.Pos
}

func (pe *PrefixExpression) End() token.Position {
	return pe.Right.End()
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         POSTFIX                       //
///////////////////////////////////////////////////////////
type PostfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) Pos() token.Position {
	return pe.Token.Pos
}

func (pe *PostfixExpression) End() token.Position {
	return pe.Left.End()
}

func (pe *PostfixExpression) expressionNode() {}

func (pe *PostfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         TERNARY                         //
///////////////////////////////////////////////////////////
type TernaryExpression struct {
	Token     token.Token
	Condition Expression
	IfTrue    Expression
	IfFalse   Expression
}

func (te *TernaryExpression) Pos() token.Position {
	return te.Token.Pos
}

func (te *TernaryExpression) End() token.Position {
	return te.IfFalse.End()
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(te.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(te.IfTrue.String())
	out.WriteString(" : ")
	out.WriteString(te.IfFalse.String())
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                          CALL                         //
///////////////////////////////////////////////////////////
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) Pos() token.Position {
	return ce.Token.Pos
}

func (ce *CallExpression) End() token.Position {
	aLen := len(ce.Arguments)
	if aLen > 0 {
		return ce.Arguments[aLen-1].End()
	}
	return ce.Function.End()
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     METHOD  CALL                      //
///////////////////////////////////////////////////////////
type MethodCallExpression struct {
	Token  token.Token
	Object Expression
	Call   Expression
}

func (mc *MethodCallExpression) Pos() token.Position {
	return mc.Token.Pos
}

func (mc *MethodCallExpression) End() token.Position {
	return mc.Call.End()
}

func (mc *MethodCallExpression) expressionNode()      {}
func (mc *MethodCallExpression) TokenLiteral() string { return mc.Token.Literal }
func (mc *MethodCallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(mc.Object.String())
	out.WriteString(".")
	out.WriteString(mc.Call.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                       CASE/ESLE                       //
///////////////////////////////////////////////////////////
type CaseExpr struct {
	Token        token.Token
	IsWholeMatch bool
	Expr         Expression
	Matches      []Expression
}

func (c *CaseExpr) Pos() token.Position {
	return c.Token.Pos
}

func (c *CaseExpr) End() token.Position {
	aLen := len(c.Matches)
	if aLen > 0 {
		return c.Matches[aLen-1].End()
	}
	return c.Expr.End()
}

func (c *CaseExpr) expressionNode()      {}
func (c *CaseExpr) TokenLiteral() string { return c.Token.Literal }

func (c *CaseExpr) String() string {
	var out bytes.Buffer

	out.WriteString("case ")
	out.WriteString(c.Expr.String())
	if c.IsWholeMatch {
		out.WriteString(" is ")
	} else {
		out.WriteString(" in ")
	}
	out.WriteString(" { ")

	matches := []string{}
	for _, m := range c.Matches {
		matches = append(matches, m.String())
	}

	out.WriteString(strings.Join(matches, " "))
	out.WriteString(" }")
	return out.String()
}

type CaseMatchExpr struct {
	Token token.Token
	Expr  Expression
	Block *BlockStatement
}

func (cm *CaseMatchExpr) Pos() token.Position {
	return cm.Token.Pos
}

func (cm *CaseMatchExpr) End() token.Position {
	return cm.Block.End()
}

func (cm *CaseMatchExpr) expressionNode()      {}
func (cm *CaseMatchExpr) TokenLiteral() string { return cm.Token.Literal }

func (cm *CaseMatchExpr) String() string {
	var out bytes.Buffer

	out.WriteString(cm.Expr.String())
	out.WriteString(" { ")
	out.WriteString(cm.Block.String())
	out.WriteString(" }")

	return out.String()
}

type CaseElseExpr struct {
	Token token.Token
	Block *BlockStatement
}

func (ce *CaseElseExpr) Pos() token.Position {
	return ce.Token.Pos
}

func (ce *CaseElseExpr) End() token.Position {
	return ce.Block.End()
}

func (ce *CaseElseExpr) expressionNode()      {}
func (ce *CaseElseExpr) TokenLiteral() string { return ce.Token.Literal }

func (ce *CaseElseExpr) String() string {
	var out bytes.Buffer

	out.WriteString("else ")
	out.WriteString(" { ")
	out.WriteString(ce.Block.String())
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                       SLICE/INDEX                     //
///////////////////////////////////////////////////////////
type SliceExpression struct {
	Token      token.Token
	StartIndex Expression
	EndIndex   Expression
}

func (se *SliceExpression) Pos() token.Position {
	return se.Token.Pos
}

func (se *SliceExpression) End() token.Position {
	return se.EndIndex.End()
}

func (se *SliceExpression) expressionNode()      {}
func (se *SliceExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SliceExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	if se.StartIndex != nil {
		out.WriteString(se.StartIndex.String())
	}
	out.WriteString(":")
	if se.EndIndex != nil {
		out.WriteString(se.EndIndex.String())
	}
	out.WriteString(")")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) Pos() token.Position {
	return ie.Token.Pos
}

func (ie *IndexExpression) End() token.Position {
	return ie.Index.End()
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}

///////////////////////////////////////////////////////////
//                     SPAWN STATEMENT                   //
///////////////////////////////////////////////////////////
type SpawnStmt struct {
	Token token.Token
	Call  Expression
}

func (ss *SpawnStmt) Pos() token.Position {
	return ss.Token.Pos
}

func (ss *SpawnStmt) End() token.Position {
	return ss.Call.End()
}

func (ss *SpawnStmt) statementNode()       {}
func (ss *SpawnStmt) TokenLiteral() string { return ss.Token.Literal }

func (ss *SpawnStmt) String() string {
	var out bytes.Buffer

	out.WriteString(ss.TokenLiteral() + " ")
	out.WriteString(ss.Call.String())
	out.WriteString("; ")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     YIELD EXPRESSION                  //
///////////////////////////////////////////////////////////
type YieldExpression struct {
	Token     token.Token
	Arguments []Expression // The arguments to yield
}

func (y *YieldExpression) Pos() token.Position {
	return y.Token.Pos
}

func (y *YieldExpression) End() token.Position {
	aLen := len(y.Arguments)
	if aLen > 0 {
		return y.Arguments[aLen-1].End()
	}
	return y.Token.Pos
}

func (y *YieldExpression) expressionNode()      {}
func (y *YieldExpression) TokenLiteral() string { return y.Token.Literal }

func (y *YieldExpression) String() string {
	var out bytes.Buffer

	out.WriteString(y.Token.Literal)
	if len(y.Arguments) != 0 {
		args := []string{}
		for _, a := range y.Arguments {
			args = append(args, a.String())
		}
		out.WriteString(" ")
		out.WriteString(strings.Join(args, ", "))
	}
	return out.String()
}

///////////////////////////////////////////////////////////
//                      FIELD LITERAL                   //
///////////////////////////////////////////////////////////
type FieldLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (f *FieldLiteral) Pos() token.Position {
	return f.Token.Pos
}

func (f *FieldLiteral) End() token.Position {
	maxLineMap := make(map[int]Expression)

	for _, value := range f.Pairs {
		v := value.(Expression)
		maxLineMap[v.End().Line] = v
	}

	maxLine := 0
	for line, _ := range maxLineMap {
		if line > maxLine {
			maxLine = line
		}
	}

	ret := maxLineMap[maxLine].(Expression)
	return ret.End()
}

func (f *FieldLiteral) expressionNode()      {}
func (f *FieldLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FieldLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range f.Pairs {
		pairs = append(pairs, key.String()+"=>"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                  PIPE OPERATOR                        //
///////////////////////////////////////////////////////////
// Pipe operator.
type Pipe struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (p *Pipe) Pos() token.Position {
	return p.Token.Pos
}

func (p *Pipe) End() token.Position {
	return p.Right.End()
}

func (p *Pipe) expressionNode()      {}
func (p *Pipe) TokenLiteral() string { return p.Token.Literal }
func (p *Pipe) String() string {
	var out bytes.Buffer

	out.WriteString(p.Left.String())
	out.WriteString(" |> ")
	out.WriteString(p.Right.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                   ENUM Literal                        //
///////////////////////////////////////////////////////////
type EnumLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (e *EnumLiteral) Pos() token.Position {
	return e.Token.Pos
}

func (e *EnumLiteral) End() token.Position {
	maxLineMap := make(map[int]Expression)

	for _, value := range e.Pairs {
		v := value.(Expression)
		maxLineMap[v.End().Line] = v
	}

	maxLine := 0
	for line, _ := range maxLineMap {
		if line > maxLine {
			maxLine = line
		}
	}

	ret := maxLineMap[maxLine].(Expression)
	return ret.End()
}

func (e *EnumLiteral) expressionNode()      {}
func (e *EnumLiteral) TokenLiteral() string { return e.Token.Literal }

func (e *EnumLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("enum ")
	out.WriteString("{")

	pairs := []string{}
	for k, v := range e.Pairs {
		pairs = append(pairs, k.String()+" = "+v.String())
	}
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//        List Comprehension(for array & string)         //
///////////////////////////////////////////////////////////
// [ Expr for Var in Value <where Cond> ] ---> Value could be array or string
type ListComprehension struct {
	Token token.Token
	Var   string
	Value Expression //value(array or string) to range over
	Cond  Expression //conditional clause(nil if there is no 'WHERE' clause)
	Expr  Expression //the result expression
}

func (lc *ListComprehension) Pos() token.Position {
	return lc.Token.Pos
}

func (lc *ListComprehension) End() token.Position {
	if lc.Cond != nil {
		return lc.Cond.End()
	}
	return lc.Value.End()
}

func (lc *ListComprehension) expressionNode()      {}
func (lc *ListComprehension) TokenLiteral() string { return lc.Token.Literal }

func (lc *ListComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("[ ")
	out.WriteString(lc.Expr.String())
	out.WriteString(" for ")
	out.WriteString(lc.Var)
	out.WriteString(" in ")
	out.WriteString(lc.Value.String())
	if lc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(lc.Cond.String())
	}
	out.WriteString(" ]")

	return out.String()
}

///////////////////////////////////////////////////////////
//             List Comprehension(for range)             //
///////////////////////////////////////////////////////////
//[Expr for Var in StartIdx..EndIdx <where Cond>]
type ListRangeComprehension struct {
	Token    token.Token
	Var      string
	StartIdx Expression
	EndIdx   Expression
	Cond     Expression //conditional clause(nil if there is no 'WHERE' clause)
	Expr     Expression //the result expression
}

func (lc *ListRangeComprehension) Pos() token.Position {
	return lc.Token.Pos
}

func (lc *ListRangeComprehension) End() token.Position {
	if lc.Cond != nil {
		return lc.Cond.End()
	}
	return lc.EndIdx.End()
}

func (lc *ListRangeComprehension) expressionNode()      {}
func (lc *ListRangeComprehension) TokenLiteral() string { return lc.Token.Literal }

func (lc *ListRangeComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("[ ")
	out.WriteString(lc.Expr.String())
	out.WriteString(" for ")
	out.WriteString(lc.Var)
	out.WriteString(" in ")
	out.WriteString(lc.StartIdx.String())
	out.WriteString("..")
	out.WriteString(lc.EndIdx.String())
	if lc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(lc.Cond.String())
	}
	out.WriteString(" ]")

	return out.String()
}

///////////////////////////////////////////////////////////
//                LIST Map Comprehension                 //
///////////////////////////////////////////////////////////
//[ Expr for Key,Value in X <where Cond>]
type ListMapComprehension struct {
	Token token.Token
	Key   string
	Value string
	X     Expression //value(hash) to range over
	Cond  Expression //Conditional clause(nil if there is no 'WHERE' clause)
	Expr Expression  //the result expression
}

func (mc *ListMapComprehension) Pos() token.Position {
	return mc.Token.Pos
}

func (mc *ListMapComprehension) End() token.Position {
	if mc.Cond != nil {
		return mc.Cond.End()
	}
	return mc.Expr.End()
}

func (mc *ListMapComprehension) expressionNode()      {}
func (mc *ListMapComprehension) TokenLiteral() string { return mc.Token.Literal }

func (mc *ListMapComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("[ ")
	out.WriteString(mc.Expr.String())
	out.WriteString(" for ")
	out.WriteString(mc.Key + ", " + mc.Value)
	out.WriteString(" in ")
	out.WriteString(mc.X.String())
	if mc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(mc.Cond.String())
	}
	out.WriteString(" ]")

	return out.String()
}

///////////////////////////////////////////////////////////
//        Hash Comprehension(for array & string)         //
///////////////////////////////////////////////////////////
//{ KeyExpr:ValExpr for Var in Value <where Cond> }  -->Value could be array or string
type HashComprehension struct {
	Token   token.Token
	Var     string
	Value   Expression //value(array or string) to range over
	Cond    Expression //conditional clause(nil if there is no 'WHERE' clause)
	KeyExpr Expression //the result Key expression
	ValExpr Expression //the result Value expression
}

func (hc *HashComprehension) Pos() token.Position {
	return hc.Token.Pos
}

func (hc *HashComprehension) End() token.Position {
	if hc.Cond != nil {
		return hc.Cond.End()
	}
	return hc.Value.End()
}

func (hc *HashComprehension) expressionNode()      {}
func (hc *HashComprehension) TokenLiteral() string { return hc.Token.Literal }

func (hc *HashComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	out.WriteString(hc.KeyExpr.String())
	out.WriteString(" : ")
	out.WriteString(hc.ValExpr.String())
	out.WriteString(" for ")
	out.WriteString(hc.Var)
	out.WriteString(" in ")
	out.WriteString(hc.Value.String())
	if hc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(hc.Cond.String())
	}
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//             Hash Comprehension(for range)             //
///////////////////////////////////////////////////////////
//{ KeyExp:ValExp for Var in StartIdx..EndIdx <where Cond> }
type HashRangeComprehension struct {
	Token    token.Token
	Var      string
	StartIdx Expression
	EndIdx   Expression
	Cond     Expression //conditional clause(nil if there is no 'WHERE' clause)
	KeyExpr  Expression //the result Key expression
	ValExpr  Expression //the result Value expression
}

func (hc *HashRangeComprehension) Pos() token.Position {
	return hc.Token.Pos
}

func (hc *HashRangeComprehension) End() token.Position {
	if hc.Cond != nil {
		return hc.Cond.End()
	}
	return hc.EndIdx.End()
}

func (hc *HashRangeComprehension) expressionNode()      {}
func (hc *HashRangeComprehension) TokenLiteral() string { return hc.Token.Literal }

func (hc *HashRangeComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	out.WriteString(hc.KeyExpr.String())
	out.WriteString(" : ")
	out.WriteString(hc.ValExpr.String())
	out.WriteString(" for ")
	out.WriteString(hc.Var)
	out.WriteString(" in ")
	out.WriteString(hc.StartIdx.String())
	out.WriteString("..")
	out.WriteString(hc.EndIdx.String())
	if hc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(hc.Cond.String())
	}
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                Hash Map Comprehension                 //
///////////////////////////////////////////////////////////
//{ KeyExpr:ValExpr for Key,Value in X <where Cond> }
type HashMapComprehension struct {
	Token   token.Token
	Key     string
	Value   string
	X       Expression //value(hash) to range over
	Cond    Expression //Conditional clause(nil if there is no 'WHERE' clause)
	KeyExpr Expression  //the result Key expression
	ValExpr Expression  //the result Value expression
}

func (mc *HashMapComprehension) Pos() token.Position {
	return mc.Token.Pos
}

func (mc *HashMapComprehension) End() token.Position {
	if mc.Cond != nil {
		return mc.Cond.End()
	}
	return mc.X.End()
}

func (mc *HashMapComprehension) expressionNode()      {}
func (mc *HashMapComprehension) TokenLiteral() string { return mc.Token.Literal }

func (mc *HashMapComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	out.WriteString(mc.KeyExpr.String())
	out.WriteString(" : ")
	out.WriteString(mc.ValExpr.String())
	out.WriteString(" for ")
	out.WriteString(mc.Key + ", " + mc.Value)
	out.WriteString(" in ")
	out.WriteString(mc.X.String())
	if mc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(mc.Cond.String())
	}
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                      Tuple LITERAL                    //
///////////////////////////////////////////////////////////
type TupleLiteral struct {
	Token   token.Token
	Members []Expression
}

func (t *TupleLiteral) Pos() token.Position {
	return t.Token.Pos
}

func (t *TupleLiteral) End() token.Position {
	aLen := len(t.Members)
	if aLen > 0 {
		return t.Members[aLen-1].End()
	}
	return t.Token.Pos
}

func (t *TupleLiteral) expressionNode()      {}
func (t *TupleLiteral) TokenLiteral() string { return t.Token.Literal }
func (t *TupleLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("(")

	members := []string{}
	for _, m := range t.Members {
		members = append(members, m.String())
	}

	out.WriteString(strings.Join(members, ", "))
	out.WriteString(")")

	return out.String()
}


//class's method modifier
type ModifierLevel int8
const (
	ModifierDefault ModifierLevel = iota
	ModifierPrivate
	ModifierProtected
	ModifierPublic
)

//for debug purpose
func (m ModifierLevel) String() string {
	switch {
	//note the last space.
	case m == ModifierPrivate:
		return "private "
	case m == ModifierProtected:
		return "protected "
	case m == ModifierPublic:
		return "public "
	}

	return ""
}

///////////////////////////////////////////////////////////
//                      CLASS LITERAL                    //
///////////////////////////////////////////////////////////
// class : parentClass { block }
type ClassLiteral struct {
	Token      token.Token
	Name       string
	Parent     string
	Members    []*LetStatement  //class's fields
	Properties map[string]*PropertyDeclStmt //class's properties
	Methods    map[string]*FunctionStatement //class's methods
	Block      *BlockStatement //mainly used for debugging purpose
	Modifier   ModifierLevel  //NOT IMPLEMENTED
}

func (c *ClassLiteral) Pos() token.Position {
	return c.Token.Pos
}

func (c *ClassLiteral) End() token.Position {
	return c.Block.End()
}

func (c *ClassLiteral) expressionNode()       {}
func (c *ClassLiteral) TokenLiteral() string { return c.Token.Literal }

func (c *ClassLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(c.TokenLiteral() + " ")
	out.WriteString(c.Name)
	if len(c.Parent) != 0 {
		out.WriteString(" : " + c.Parent + " ")
	}

	out.WriteString("{ ")
	out.WriteString(c.Block.String())
	out.WriteString("} ")

	return out.String()
}

//class classname : parentClass { block }
///////////////////////////////////////////////////////////
//                     CLASS STATEMENT                   //
///////////////////////////////////////////////////////////
type ClassStatement struct {
	Token           token.Token
	Name            *Identifier //Class name
	ClassLiteral    *ClassLiteral
}

func (c *ClassStatement) Pos() token.Position {
	return c.Token.Pos
}

func (c *ClassStatement) End() token.Position {
	return c.ClassLiteral.Block.End()
}

func (c *ClassStatement) statementNode() {}
func (c *ClassStatement) TokenLiteral() string { return c.Token.Literal }
func (c *ClassStatement) String() string {
	var out bytes.Buffer

	out.WriteString(c.Token.Literal + " ")
	out.WriteString(c.Name.String())

	if len(c.ClassLiteral.Parent) > 0 {
		out.WriteString(" : " + c.ClassLiteral.Parent)
	}

	out.WriteString("{ ")
	out.WriteString(c.ClassLiteral.Block.String())
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                   NEW EXPRESSION                      //
///////////////////////////////////////////////////////////

type NewExpression struct {
	Token     token.Token
	Class     Expression
	Arguments []Expression
}

func (c *NewExpression) Pos() token.Position {
	return c.Token.Pos
}

func (c *NewExpression) End() token.Position {
	aLen := len(c.Arguments)
	if aLen > 0 {
		return c.Arguments[aLen-1].End()
	}
	return c.Class.End()
}

func (n *NewExpression) expressionNode()      {}
func (n *NewExpression) TokenLiteral() string { return n.Token.Literal }
func (n *NewExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range n.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(n.TokenLiteral() + " ")
	out.WriteString(n.Class.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(") ")

	return out.String()
}

//class's property declaration
type PropertyDeclStmt struct {
	Token         token.Token
	Name          *Identifier      //property name
	Getter        *GetterStmt      //getter
	Setter        *SetterStmt      //setter
	ModifierLevel ModifierLevel   //property's modifier
}

func (p *PropertyDeclStmt) Pos() token.Position {
	return p.Token.Pos
}

func (p *PropertyDeclStmt) End() token.Position {
	if p.Getter == nil {
		return p.Setter.End()
	}
	return p.Getter.End()
}

func (p *PropertyDeclStmt) statementNode()       {}
func (p *PropertyDeclStmt) TokenLiteral() string { return p.Token.Literal }

func (p *PropertyDeclStmt) String() string {
	var out bytes.Buffer

	out.WriteString(p.ModifierLevel.String())

	out.WriteString("property ")
	out.WriteString(p.Name.String() +" ")
	out.WriteString("{ ")

	if p.Getter != nil {
		out.WriteString(p.Getter.String())
	}

	if p.Setter != nil {
		out.WriteString(p.Setter.String())
	}

	out.WriteString("} ")
	return out.String()
}

//property's getter statement
type GetterStmt struct {
	Token token.Token
	Body *BlockStatement
}

func (g *GetterStmt) Pos() token.Position {
	return g.Token.Pos
}

func (g *GetterStmt) End() token.Position {
	return g.Body.End()
}

func (g *GetterStmt) statementNode()       {}
func (g *GetterStmt) TokenLiteral() string { return g.Token.Literal }

func (g *GetterStmt) String() string {
	var out bytes.Buffer

	out.WriteString("get")
	if len(g.Body.Statements) == 0 {
		out.WriteString("; ")
	} else {
		out.WriteString("{")
		out.WriteString(g.Body.String())
		out.WriteString("} ")
	}

	return out.String()
}

//property's setter statement
//setter variable is always 'value' like c#
type SetterStmt struct {
	Token token.Token
	Body *BlockStatement
}

func (s *SetterStmt) Pos() token.Position {
	return s.Token.Pos
}

func (s *SetterStmt) End() token.Position {
	return s.Body.End()
}

func (s *SetterStmt) statementNode()       {}
func (s *SetterStmt) TokenLiteral() string { return s.Token.Literal }

func (s *SetterStmt) String() string {
	var out bytes.Buffer

	out.WriteString("set")
	if len(s.Body.Statements) == 0 {
		out.WriteString("; ")
	} else {
		out.WriteString("{")
		out.WriteString(s.Body.String())
		out.WriteString("} ")
	}

	return out.String()
}
