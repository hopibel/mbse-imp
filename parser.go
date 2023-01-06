package main

import (
	"bytes"
	"os"
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

// Lexer

type Lexer struct {
	s       *string
	tokType TokType
	tok     bytes.Buffer // WriteString, String
}

func newLexer(f string) *Lexer {
	code, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}
	c := string(code)
	return &Lexer{&c, TokEOF, bytes.Buffer{}}
}

// Parser

type Parser struct {
	s *string
}

func (p Parser) parse_file(f string) Program {
	lexer := newLexer(f)

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
