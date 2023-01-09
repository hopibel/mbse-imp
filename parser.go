package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
)

// IMP parser grammar

/*
prog    ::= stmts
block   ::= "{" stmts "}"
stmts   ::= stmt ";" stmts
          | stmt ";"
stmt    ::= restTMP
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
	default: // special case: =[^=]
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

func (p *Parser) parse_file(f string) Program {
	p.lexer = newLexer(f)

	return Program{}
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
