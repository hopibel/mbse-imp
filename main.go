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

func interpret_file(f string, verbose bool) {
	if verbose {
		lexer := newFileLexer(f)
		lexer.lex_file()
		fmt.Println()
	}
	parser := newParser()
	prog, err := parser.parse_file(f)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to parse", f)
	}
	if verbose {
		fmt.Println("Pretty print AST:")
		fmt.Print(prog.pretty(), "\n\n")
	}
	// typecheck program
	ts := newTyState()
	if prog.check(ts) {
		if verbose {
			fmt.Printf("Successfully type-checked %s\n\n", f)
		}
		// run program
		vs := newValState()
		prog.eval(vs)
	} else {
		fmt.Printf("%s contains type errors\n", f)
	}
}

func main() {
	var verbose bool
	var fname string
	args := os.Args
	if len(args) == 2 {
		// mbse-imp <file>
		verbose = false
		fname = args[1]
	} else if len(args) == 3 && (args[1] == "-v" || args[2] == "-v") {
		// -v: verbose
		if args[1] == "-v" {
			verbose = true
			fname = args[2]
		} else {
			verbose = true
			fname = args[1]
		}
	} else {
		fmt.Printf("usage: %s [-v] <filename>\n", args[0])
		os.Exit(1)
	}

	interpret_file(fname, verbose)
}
