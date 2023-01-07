package main

import (
	"bytes"
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
	TokName
	TokEOF
)

// Lexer compiled regexes

var whitespace *regexp.Regexp = regexp.MustCompile("\\s")

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
func (l Lexer) next() (bool, error) {
	// detect EOF
	if l.cursor == len(l.s) {
		l.tokType = TokEOF
		l.tok.Reset()
		return true, nil
	}

	// ignore whitespace
	for whitespace.MatchString(l.s[l.cursor : l.cursor+1]) {
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
		panic("Lexer.next invalid token")
	}

	return true, nil
}

func (l Lexer) lex_int() bool {
	return false
}

func (l Lexer) lex_bool() bool {
	return false
}

func (l Lexer) lex_ident() bool {
	return false
}

func (l Lexer) lex_operator() bool {
	return false
}

// Parser

type Parser struct {
	s *string
}

func (p Parser) parse_file(f string) Program {
	lexer := newLexer(f)

	// debug
	while

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
