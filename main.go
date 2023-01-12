package main

import (
	"fmt"
	"os"
)

// Simple imperative language

/*
vars       Variable names, start with lower-case letter

prog      ::= statement
block     ::= "{" statement "}"
statement ::=  statement ";" statement           -- Command sequence
            |  vars ":=" exp                     -- Variable declaration
            |  vars "=" exp                      -- Variable assignment
            |  "while" exp block                 -- While
            |  "if" exp block "else" block       -- If-then-else
            |  "print" exp                       -- Print

exp ::= 0 | 1 | -1 | ...     -- Integers
     | "true" | "false"      -- Booleans
     | exp "+" exp           -- Addition
     | exp "*" exp           -- Multiplication
     | exp "||" exp          -- Disjunction
     | exp "&&" exp          -- Conjunction
     | "!" exp               -- Negation
     | exp "==" exp          -- Equality test
     | exp "<" exp           -- Lesser test
     | "(" exp ")"           -- Grouping of expressions
     | vars                  -- Variables
*/

// Interpreter

func interpret_file(f string) {
	lexer := newLexer(f)
	lexer.lex_file()
	// parser := newParser()
	// prog, err := parser.parse_file(f)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println("Failed to parse", f)
	// }
	// fmt.Println(prog)
}

// Examples

func run(e Exp) {
	s := make(map[string]Val)
	t := make(map[string]Type)
	fmt.Printf("\n ******* ")
	fmt.Printf("\n %s", e.pretty())
	fmt.Printf("\n %s", showVal(e.eval(s)))
	fmt.Printf("\n %s", showType(e.infer(t)))
}

func ex1() {
	ast := plus(mult(number(1), number(2)), number(0))

	run(ast)
}

func ex2() {
	ast := and(boolean(false), number(0))
	run(ast)
}

func ex3() {
	ast := or(boolean(false), number(0))
	run(ast)
}

func main() {
	// ex1()
	// ex2()
	// ex3()
	// fmt.Print("\n\n")

	args := os.Args
	if len(args) != 2 {
		fmt.Printf("usage: %s <filename>\n", args[0])
		os.Exit(1)
	}

	interpret_file(args[1])
}
