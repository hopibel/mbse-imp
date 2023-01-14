package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// IMP parser grammar

/*
prog    ::= seq
block   ::= "{" seq "}"
seq     ::= stmt ";" seq2
seq2    ::= stmt ";" seq2
          | epsilon
stmt    ::= vars ":=" exp
          | vars "=" exp
          | "while" exp block
          | "if" exp block "else" block
          | "print" exp
exp     ::= exp2 comp
comp    ::= "==" exp2 comp
          | "<" exp2 comp
          | epsilon
exp2    ::= term exp3
exp3    ::= "+" term exp3
          | "||" term exp3
          | epsilon
term    ::= factor term2
term2   ::= "*" factor term2
          | "&&" factor term2
          | epsilon
factor  ::= lit | vars
          | "!" factor
          | "(" exp ")"
lit     ::= 0 | 1 | -1 | ...
          | "true" | "false"
*/

// Tokens

type TokType int

const (
	TokSemicolon TokType = iota
	TokBraceOpen
	TokBraceClose
	TokDecl
	TokAssign
	TokWhile
	TokIf
	TokElse
	TokPrint
	TokInt
	TokBool
	TokPlus
	TokMult
	TokOr
	TokAnd
	TokNot
	TokEqual
	TokLess
	TokParenOpen
	TokParenClose
	TokName
	TokEOF
)

// Lexer

type Lexer struct {
	s       string       // source code
	cursor  int          // current position in source
	tokType TokType      // current token type
	tok     bytes.Buffer // current token string
	line    int          // current line
}

func newLexer(f string) *Lexer {
	code, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}
	c := string(code)
	return &Lexer{c, 0, TokEOF, bytes.Buffer{}, 1}
}

// Lexer compiled regexes

var rWhitespace = regexp.MustCompile(`^\s+`)
var rNewline = regexp.MustCompile(`\n`)
var rInt = regexp.MustCompile(`^-?\d+`)
var rBool = regexp.MustCompile(`^(true|false)`)
var rIdent = regexp.MustCompile(`^[a-z]\w*`)
var rOperator = regexp.MustCompile(`^(:=|=[^=]|\+|\*|\|\||&&|!|==|<)`) // TODO: test all operators

// next token
func (l *Lexer) next() (bool, error) {
	l.tok.Reset()
	if l.eol() {
		return true, nil
	}

	// ignore whitespace, count newlines
	s := l.s[l.cursor:]
	loc := rWhitespace.FindStringIndex(s)
	if loc != nil {
		ws := s[loc[0]:loc[1]]
		ns := rNewline.FindAllStringIndex(ws, -1)
		if ns != nil {
			l.line += len(ns)
		}
		l.cursor += loc[1]
	}

	// check EOL again after skipping whitespace
	if l.eol() {
		return true, nil
	}

	// token matchers. cursor is only advanced on successful
	switch {
	case l.lex_int(): // integer literals
	case l.lex_bool(): // boolean literals
	case l.lex_ident(): // vars and keywords
	case l.lex_operator(): // operators
	case l.lex_brace(): // parens and curly braces
	case l.lex_semi(): // semicolon
	default:
		// TODO: return informative error message https://go.dev/doc/tutorial/handle-errors
		panic(fmt.Sprintf("Lexer.next: unexpected character on line %d: \"%c\"", l.line, l.s[l.cursor]))
	}

	return true, nil
}

// detect EOF
func (l *Lexer) eol() bool {
	if l.cursor == len(l.s) {
		l.tokType = TokEOF
		l.tok.Reset()
		return true
	}
	return false
}

func (l *Lexer) lex_int() bool {
	s := l.s[l.cursor:]
	loc := rInt.FindStringIndex(s)
	if loc == nil {
		return false
	}
	tok := s[loc[0]:loc[1]]
	l.tok.WriteString(tok)
	l.tokType = TokInt
	l.cursor += loc[1]
	return true
}

func (l *Lexer) lex_bool() bool {
	s := l.s[l.cursor:]
	loc := rBool.FindStringIndex(s)
	if loc == nil {
		return false
	}
	tok := s[loc[0]:loc[1]]
	l.tok.WriteString(tok)
	l.tokType = TokBool
	l.cursor += loc[1]
	return true
}

func (l *Lexer) lex_ident() bool {
	// slurp ^[a-z]\w
	s := l.s[l.cursor:]
	loc := rIdent.FindStringIndex(s)
	if loc == nil {
		return false
	}
	tok := s[loc[0]:loc[1]]
	// test for keywords, else variable name
	switch tok {
	case "while":
		l.tokType = TokWhile
	case "if":
		l.tokType = TokIf
	case "else":
		l.tokType = TokElse
	case "print":
		l.tokType = TokPrint
	default: // variable name
		l.tokType = TokName
	}
	l.tok.WriteString(tok)
	l.cursor += loc[1]
	return true
}

func (l *Lexer) lex_operator() bool {
	s := l.s[l.cursor:]
	loc := rOperator.FindStringIndex(s)
	if loc == nil {
		return false
	}
	tok := s[loc[0]:loc[1]]
	switch tok {
	case ":=":
		l.tokType = TokDecl
	case "+":
		l.tokType = TokPlus
	case "*":
		l.tokType = TokMult
	case "||":
		l.tokType = TokOr
	case "&&":
		l.tokType = TokAnd
	case "!":
		l.tokType = TokNot
	case "==":
		l.tokType = TokEqual
	case "<":
		l.tokType = TokLess
	default: // assignment =
		l.tokType = TokAssign
		tok = string(s[loc[0]])
		loc[1]--
	}
	l.tok.WriteString(tok)
	l.cursor += loc[1]
	return true
}

func (l *Lexer) lex_brace() bool {
	tok := l.s[l.cursor]
	switch tok {
	case '{':
		l.tokType = TokBraceOpen
	case '}':
		l.tokType = TokBraceClose
	case '(':
		l.tokType = TokParenOpen
	case ')':
		l.tokType = TokParenClose
	default:
		return false
	}
	l.tok.WriteByte(tok)
	l.cursor++
	return true
}

func (l *Lexer) lex_semi() bool {
	tok := l.s[l.cursor]
	if tok == ';' {
		l.tokType = TokSemicolon
		l.tok.WriteByte(tok)
		l.cursor++
		return true
	}
	return false
}

// debug method to test tokenizer/lexer
func (l *Lexer) lex_file() {
	fmt.Println("Token stream:")
	for status, _ := l.next(); status && l.tokType != TokEOF; status, _ = l.next() {
		switch l.tokType {
		case TokSemicolon:
			fmt.Print("TokSemicolon")
		case TokBraceOpen:
			fmt.Print("TokBraceOpen")
		case TokBraceClose:
			fmt.Print("TokBraceClose")
		case TokDecl:
			fmt.Print("TokDecl")
		case TokAssign:
			fmt.Print("TokAssign")
		case TokWhile:
			fmt.Print("TokWhile")
		case TokIf:
			fmt.Print("TokIf")
		case TokElse:
			fmt.Print("TokElse")
		case TokPrint:
			fmt.Print("TokPrint")
		case TokInt:
			fmt.Print("TokInt")
		case TokBool:
			fmt.Print("TokBool")
		case TokPlus:
			fmt.Print("TokPlus")
		case TokMult:
			fmt.Print("TokMult")
		case TokOr:
			fmt.Print("TokOr")
		case TokAnd:
			fmt.Print("TokAnd")
		case TokNot:
			fmt.Print("TokNot")
		case TokEqual:
			fmt.Print("TokEqual")
		case TokLess:
			fmt.Print("TokLess")
		case TokParenOpen:
			fmt.Print("TokParenOpen")
		case TokParenClose:
			fmt.Print("TokParenClose")
		case TokName:
			fmt.Print("TokName")
		case TokEOF:
			panic("lexer test should not reach EOF")
		default:
			panic("unrecognized token")
		}
		fmt.Printf("(%s) ", l.tok.String())
	}
	fmt.Println()
}

// Parser

type Parser struct {
	lexer *Lexer
}

func newParser() *Parser {
	return &Parser{nil}
}

func (p *Parser) err_expected(what string) error {
	return fmt.Errorf("expected %s on line %d, found \"%s\"", what, p.lexer.line, p.lexer.tok.String())
}

func (p *Parser) parse_file(f string) (Program, error) {
	p.lexer = newLexer(f)
	p.lexer.next()
	return p.parse_prog()
}

func (p *Parser) parse_prog() (Program, error) {
	prog, err := p.parse_seq()
	return (Program)(prog), err
}

func (p *Parser) parse_seq() (Stmt, error) {
	stmt, err := p.parse_stmt()
	if err != nil {
		return stmt, err
	}
	// check for ";"
	if p.lexer.tokType != TokSemicolon {
		return stmt, p.err_expected("semicolon")
	}
	p.lexer.next()
	// seq2
	return p.parse_seq2(stmt)
}

func (p *Parser) parse_seq2(stmt Stmt) (Stmt, error) {
	// epsilon if next token is close brace or EOF
	if tok := p.lexer.tokType; tok == TokBraceClose || tok == TokEOF {
		return stmt, nil
	}
	// otherwise parse stmt ; seq2
	stmt2, err := p.parse_stmt()
	if err != nil {
		return Seq{stmt, stmt2}, err
	}
	// check for ";"
	if p.lexer.tokType != TokSemicolon {
		return Seq{stmt, stmt2}, p.err_expected("semicolon")
	}
	p.lexer.next()
	// seq2
	seq2, err := p.parse_seq2(stmt2)
	return Seq{stmt, seq2}, err
}

func (p *Parser) parse_stmt() (Stmt, error) {
	// TODO: impl missing
	switch p.lexer.tokType {
	case TokName:
		lhs := p.lexer.tok.String()
		p.lexer.next()
		switch p.lexer.tokType {
		case TokDecl:
			p.lexer.next()
			rhs, err := p.parse_exp()
			return Decl{lhs, rhs}, err
		case TokAssign:
			p.lexer.next()
			rhs, err := p.parse_exp()
			return Assign{lhs, rhs}, err
		default:
			return Seq{}, p.err_expected("declaration or assignment")
		}
	case TokWhile:
		p.lexer.next()
		cond, err := p.parse_exp()
		if err != nil {
			return Seq{}, err
		}
		body, err := p.parse_block()
		return While{cond, body}, err
	case TokIf:
		p.lexer.next()
		cond, err := p.parse_exp()
		if err != nil {
			return Seq{}, err
		}
		thenStmt, err := p.parse_block()
		if err != nil {
			return Seq{}, err
		}
		if p.lexer.tokType != TokElse {
			return Seq{}, p.err_expected("keyword \"else\"")
		}
		p.lexer.next()
		elseStmt, err := p.parse_block()
		return IfThenElse{cond, thenStmt, elseStmt}, err
	case TokPrint:
		p.lexer.next()
		exp, err := p.parse_exp()
		return Print{exp}, err
	default:
		return Seq{}, p.err_expected("name or keyword")
	}
}

func (p *Parser) parse_block() (Stmt, error) {
	if p.lexer.tokType != TokBraceOpen {
		return Seq{}, p.err_expected("\"{\"")
	}
	p.lexer.next()
	block, err := p.parse_seq()
	if err != nil {
		return block, err
	} else if p.lexer.tokType != TokBraceClose {
		return block, p.err_expected("\"}\"")
	}
	p.lexer.next()
	return block, err
}

func (p *Parser) parse_exp() (Exp, error) {
	exp2, err := p.parse_exp2()
	if err != nil {
		return exp2, err
	}
	exp, err := p.parse_comp(exp2)
	return exp, err
}

func (p *Parser) parse_comp(lhs Exp) (Exp, error) {
	switch p.lexer.tokType {
	case TokEqual:
		p.lexer.next()
		rhs, err := p.parse_exp2()
		if err != nil {
			return rhs, err
		}
		ret := Equal{lhs, rhs}
		return p.parse_comp(ret)
	case TokLess:
		p.lexer.next()
		rhs, err := p.parse_exp2()
		if err != nil {
			return rhs, err
		}
		ret := Less{lhs, rhs}
		return p.parse_comp(ret)
	default:
		return lhs, nil
	}
}

func (p *Parser) parse_exp2() (Exp, error) {
	term, err := p.parse_term()
	if err != nil {
		return term, err
	}
	exp, err := p.parse_exp3(term)
	return exp, err
}

func (p *Parser) parse_exp3(lhs Exp) (Exp, error) {
	tok := p.lexer.tokType
	if tok != TokPlus && tok != TokOr {
		return lhs, nil
	}
	p.lexer.next()
	rhs, err := p.parse_term()
	if err != nil {
		return rhs, err
	}
	switch tok {
	case TokPlus:
		return p.parse_exp3(Plus([2]Exp{lhs, rhs}))
	case TokOr:
		return p.parse_exp3(Or([2]Exp{lhs, rhs}))
	default:
		panic("should not reach")
	}
}

func (p *Parser) parse_term() (Exp, error) {
	factor, err := p.parse_factor()
	if err != nil {
		return factor, err
	}
	return p.parse_term2(factor)
}

func (p *Parser) parse_term2(lhs Exp) (Exp, error) {
	tok := p.lexer.tokType
	if tok != TokMult && tok != TokAnd {
		return lhs, nil
	}
	p.lexer.next()
	rhs, err := p.parse_factor()
	if err != nil {
		return rhs, err
	}
	switch tok {
	case TokMult:
		return p.parse_term2(Mult([2]Exp{lhs, rhs}))
	case TokAnd:
		return p.parse_term2(And([2]Exp{lhs, rhs}))
	default:
		panic("should not reach")
	}
}

func (p *Parser) parse_factor() (Exp, error) {
	switch p.lexer.tokType {
	case TokInt:
		num, err := strconv.Atoi(p.lexer.tok.String())
		if err != nil {
			return Num(0), err
		}
		p.lexer.next()
		return Num(num), nil
	case TokBool:
		val := p.lexer.tok.String()
		p.lexer.next()
		switch val {
		case "true":
			return Bool(true), nil
		case "false":
			return Bool(false), nil
		default:
			panic("should not reach") // lexer already checked true/false
		}
	case TokName:
		name := Var(p.lexer.tok.String())
		p.lexer.next()
		return name, nil
	case TokNot:
		p.lexer.next()
		factor, err := p.parse_factor()
		return Not(factor), err
	case TokParenOpen:
		p.lexer.next()
		exp, err := p.parse_exp()
		if err != nil {
			return exp, err
		} else if p.lexer.tokType != TokParenClose {
			return exp, p.err_expected("\")\"")
		}
		p.lexer.next()
		return exp, err
	default:
		return Plus{}, p.err_expected("value or expression")
	}
}

// Helper functions to build ASTs by hand

func number(x int) Exp {
	return Num(x)
}

func boolean(x bool) Exp {
	return Bool(x)
}

func plus(x, y Exp) Exp {
	return (Plus)([2]Exp{x, y})

	// The type Plus is defined as the two element array consisting of Exp elements.
	// Plus and [2]Exp are isomorphic but different types.
	// We first build the AST value [2]Exp{x,y}.
	// Then cast this value (of type [2]Exp) into a value of type Plus.

}

func mult(x, y Exp) Exp {
	return (Mult)([2]Exp{x, y})
}

func and(x, y Exp) Exp {
	return (And)([2]Exp{x, y})
}

func or(x, y Exp) Exp {
	return (Or)([2]Exp{x, y})
}
