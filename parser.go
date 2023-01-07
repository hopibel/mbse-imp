package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
)

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

// Lexer compiled regexes

var rWhitespace = regexp.MustCompile(`\s`)

// Lexer

type Lexer struct {
	s       string       // source code
	cursor  int          // current position in source
	tokType TokType      // current token type
	tok     bytes.Buffer // current token string
}

func newLexer(f string) *Lexer {
	code, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}
	c := string(code)
	return &Lexer{c, 0, TokEOF, bytes.Buffer{}}
}

// next token
func (l *Lexer) next() (bool, error) {
	// detect EOF
	if l.cursor == len(l.s) {
		l.tokType = TokEOF
		l.tok.Reset()
		return true, nil
	}

	// ignore whitespace
	for rWhitespace.MatchString(string(l.s[l.cursor])) {
		l.cursor++
	}

	// token matchers. cursor is only advanced on successful
	switch {
	case l.lex_int(): // integer literals
	case l.lex_bool(): // boolean literals
	case l.lex_ident(): // vars and keywords
	case l.lex_operator(): // operators
	default:
		// TODO: return informative error message https://go.dev/doc/tutorial/handle-errors
		panic(fmt.Sprintf("Lexer.next: unexpected character at position %d:\n%s", l.cursor, l.s[l.cursor:]))
	}

	return true, nil
}

func (l *Lexer) lex_int() bool {
	return false
}

func (l *Lexer) lex_bool() bool {
	return false
}

var rIdent = regexp.MustCompile(`^[a-z]\w*`)

func (l *Lexer) lex_ident() bool {
	// slurp ^[a-z]\w
	loc := rIdent.FindStringIndex(l.s[l.cursor:])
	if loc == nil {
		return false
	}
	tok := l.s[loc[0]:loc[1]]
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
	l.cursor = loc[1]

	return true
}

func (l *Lexer) lex_operator() bool {
	return false
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

	// debug
	for status, _ := p.lexer.next(); status && p.lexer.tokType != TokEOF; status, _ = p.lexer.next() {
		switch p.lexer.tokType {
		case TokSemicolon:
			fmt.Println("TokSemicolon")
		case TokBraceOpen:
			fmt.Println("TokBraceOpen")
		case TokBraceClose:
			fmt.Println("TokBraceClose")
		case TokDecl:
			fmt.Println("TokDecl")
		case TokAssign:
			fmt.Println("TokAssign")
		case TokWhile:
			fmt.Println("TokWhile")
		case TokIf:
			fmt.Println("TokIf")
		case TokElse:
			fmt.Println("TokElse")
		case TokPrint:
			fmt.Println("TokPrint")
		case TokInt:
			fmt.Printf("TokInt: %s\n", p.lexer.tok.String())
		case TokBool:
			fmt.Printf("TokBool: %s\n", p.lexer.tok.String())
		case TokPlus:
			fmt.Println("TokPlus")
		case TokMult:
			fmt.Println("TokMult")
		case TokOr:
			fmt.Println("TokOr")
		case TokAnd:
			fmt.Println("TokAnd")
		case TokNot:
			fmt.Println("TokNot")
		case TokEqual:
			fmt.Println("TokEqual")
		case TokLess:
			fmt.Println("TokLess")
		case TokParenOpen:
			fmt.Println("TokParenOpen")
		case TokParenClose:
			fmt.Println("TokParenClose")
		case TokName:
			fmt.Printf("TokName: %s\n", p.lexer.tok.String())
		case TokEOF:
			panic("lexer test should not reach EOF")
		default:
			panic("unrecognized token")
		}
		fmt.Println("cursor:", p.lexer.cursor)
	}

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
