package main

import (
	"reflect"
	"testing"
)

func TestParserGood(t *testing.T) {
	tests := []struct {
		name string
		code string
		want Program
	}{
		// Sequences
		{"single stmt", "print 42;", seq(printStmt(Num(42)))},
		{"seq", "print 42; print 54;", seq(printStmt(Num(42)), printStmt(Num(54)))},
		{"nested seq", "print 42; print 54; print 9001;",
			seq(printStmt(Num(42)), printStmt(Num(54)), printStmt(Num(9001)))},
		{"newline", "print 42; \n\n print 54;\n", seq(printStmt(Num(42)), printStmt(Num(54)))},
		{"comment", "print 42; // this is a comment", printStmt(Num(42))},

		// Statements
		{":=", "x := 42;", seq(Decl{"x", Num(42)})},
		{"=", "x = 42;", seq(Assign{"x", Num(42)})},
		{"while body with 1 stmt",
			"while true {print 42;};",
			seq(While{Bool(true), printStmt(Num(42))})},
		{"while body with >1 stmt",
			"while true {print 42; print 54;};",
			seq(While{Bool(true), seq(printStmt(Num(42)), printStmt(Num(54)))})},
		{"if-then-else",
			"if true {print 42;} else {print 54;};",
			seq(IfThenElse{Bool(true), printStmt(Num(42)), printStmt(Num(54))})},
		{"print", "print 42;", printStmt(Num(42))},

		// Expressions
		{"==", "print x == y;", printStmt(equal(Var("x"), Var("y")))},
		{"<", "print x < y;", printStmt(less(Var("x"), Var("y")))},
		{"+", "print x + y;", printStmt(plus(Var("x"), Var("y")))},
		{"||", "print x || y;", printStmt(or(Var("x"), Var("y")))},
		{"*", "print x * y;", printStmt(mult(Var("x"), Var("y")))},
		{"&&", "print x && y;", printStmt(and(Var("x"), Var("y")))},
		{"int", "print 42;", printStmt(Num(42))},
		{"bool", "print true;", printStmt(Bool(true))},
		{"vars", "print x;", printStmt(Var("x"))},
		{"not (!)", "print !x;", printStmt(not(Var("x")))},
		{"(exp) => exp", "print (x);", printStmt(Var("x"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newParser().parse_fromstring(tt.code)
			failed := false
			if err != nil {
				t.Errorf("Parser returned error: %s", err.Error())
				failed = true
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_fromstring() = %v, want %v", got, tt.want)
				failed = true
			}
			if failed {
				t.Log("Code:", tt.code)
			}
		})
	}
}

// expected failures
func TestParserBad(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"unexpected character", "%"},
		{"unexpected token", "42"},
		{"missing semicolon between stmts", "x := 42 x = 54;"},
		{"missing semicolon at end", "x := 42; x = 54"},
		{"decl/assign", "x < 42;"},
		{"bad while cond", "while < {print 54;};"},
		{"bad if cond", "if == {print 42;} else {print 54;};"},
		{"bad then stmt", "if true {42;} else {print 54;};"},
		{"missing else", "if true {print 42;};"},
		{"missing opening brace", "if true print 42;};"},
		{"missing closing brace", "if true {print 42;;"},
		{"bad equal rhs", "x := 42 == ;"},
		{"bad less rhs", "x := 42 < ;"},
		{"bad plus rhs", "x := 42 + ;"},
		{"bad or rhs", "x := true || ;"},
		{"bad mult rhs", "x := 6 * ;"},
		{"bad and rhs", "x := true && ;"},
		{"bad paren exp", "x := (;);"},
		{"missing close paren", "x := (42;"},
		{"bad factor", "x := ; + ;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newParser().parse_fromstring(tt.code)
			if err == nil {
				t.Errorf("Parser accepted invalid code: %s", tt.code)
			}
		})
	}
}

func TestTypeChecker(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		// Sequences
		{"seq", "print 42; print 54;", true},
		{"seq2", "x := 42 < true; print 54 < false;", false},
		{"seq3", "print 42; print 54 < false;", false},

		// Statements
		{"assign", "x := 42; x = 54;", true},
		{"assign2", "x := 42; x = true;", false},
		{"assign3", "x = 42;", false},
		{"while", "while true {print 42;};", true},
		{"while2", "while 42 {print 42;};", false},
		{"if-then-else", "if false {print 42;} else {print 54;};", true},
		{"if-then-else2", "if true {print 42 < true;} else {print 54;};", false},
		{"if-then-else3", "if 42 {print 42;} else {print 54;};", false},
		{"print", "print 42;", true},

		// Expressions
		{"equal", "x := true == false;", true},
		{"equal2", "x := true == 54;", false},
		{"less", "x := 42 < 54;", true},
		{"less2", "x := 42 < true;", false},
		{"plus", "x := 42 + 54;", true},
		{"plus2", "x := 42 + false;", false},
		{"mult", "x := 6 * 9;", true},
		{"mult2", "x := 6 * false;", false},
		{"(exp) => exp", "x := (42+54);", true},
		{"(exp) => exp 2", "x := (42+false);", false},

		// note: short circuit is supported but fails type check which requires two bools
		{"or", "x := false || true;", true},
		{"or2", "x := false || 42;", false},
		{"or sc", "x := true || false;", true},
		{"or sc2", "x := true || 54;", false},

		{"and", "x := true && true;", true},
		{"and2", "x := true && 42;", false},
		{"and sc", "x := false && true;", true},
		{"and sc2", "x := false && 54;", false},

		{"not", "x := !true;", true},
		{"not2", "x := !42;", false},
		{"not3", "x := true; y := !x;", true},
		{"not4", "x := 54; y := !x;", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := newParser().parse_fromstring(tt.code)
			failed := false
			if err != nil {
				t.Errorf("Parser returned error: %s", err.Error())
				failed = true
			}
			got := prog.check(newTyState())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Program.check() = %v, want %v", got, tt.want)
				failed = true
			}
			if failed {
				t.Log("Code:", tt.code)
			}
		})
	}
}

func TestEvaluator(t *testing.T) {
	tests := []struct {
		name string
		code string
		want Val // convention: output stored in "x"
	}{
		// Sequences
		{"seq", "x := 42; y := 12; x = x + y;", mkInt(54)},
		{"seq2", "x := 42; x = false;", mkUndefined()},
		{"seq3", "x := 42; x = x + true;", mkUndefined()},

		// Statements
		{"declare", "x := 42;", mkInt(42)},
		{"declare2", "x := 42; x := true;", mkBool(true)},
		{"assign", "x := 42; x = 54;", mkInt(54)},
		{"assign2", "x = 54;", mkUndefined()},
		{"print", "x:=42; print x;", mkInt(42)},

		{"while", "n := 1; x := 0; while n<11 {x=x+n; n=n+1;};", mkInt(55)},
		// general case: declaration updates global variable if types match
		{"while2", "n := 1; x := 0; while n<11 {x:=x+n; n=n+1;};", mkInt(55)},
		{"while3", "b := true; x := 42; while b {x:=true; b=false;};", mkInt(42)},

		{"while cond bad type", "while 42 {x := 42;};", mkUndefined()},

		{"if-then-else", "if true {x := 42;} else {x := 54;};", mkUndefined()},
		{"if-then-else2", "x:=0; if true {x = 42;} else {x = 54;};", mkInt(42)},
		{"if-then-else3", "x:=0; if false {x = 42;} else {x = 54;};", mkInt(54)},
		// general case: decl updates global if same type
		{"if-then-else4", "x:=0; if true {x := 42;} else {x := 54;};", mkInt(42)},
		{"if-then-else5", "x:=42; if false {x := true;} else {x := false;};", mkInt(42)},

		{"if cond bad type", "if 42 {x := 42;} else {x := 54;};", mkUndefined()},

		// Expressions
		{"equal", "x := true == false;", mkBool(false)},
		{"equal2", "x := 42 == 42;", mkBool(true)},
		{"equal3", "x := true == 54;", mkUndefined()},
		{"less", "x := 42 < 54;", mkBool(true)},
		{"less2", "x := 42 < true;", mkUndefined()},
		{"plus", "x := 42 + 54;", mkInt(96)},
		{"plus2", "x := 42 + false;", mkUndefined()},
		{"mult", "x := 6 * 9;", mkInt(54)},
		{"mult2", "x := 6 * false;", mkUndefined()},
		{"(exp) => exp", "x := (42+54);", mkInt(96)},
		{"(exp) => exp 2", "x := (42+false);", mkUndefined()},

		// note: short circuit is supported but fails type check which requires two bools
		{"or", "x := false || true;", mkBool(true)},
		{"or2", "x := false || 42;", mkUndefined()},
		{"or sc", "x := true || false;", mkBool(true)},
		{"or sc2", "x := true || 54;", mkBool(true)},

		{"and", "x := true && true;", mkBool(true)},
		{"and2", "x := true && 42;", mkUndefined()},
		{"and sc", "x := false && true;", mkBool(false)},
		{"and sc2", "x := false && 54;", mkBool(false)},

		{"not", "x := !true;", mkBool(false)},
		{"not2", "x := !42;", mkUndefined()},
		{"not3", "y := true; x := !y;", mkBool(false)},
		{"not4", "y := 54; x := !y;", mkUndefined()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := newParser().parse_fromstring(tt.code)
			failed := false
			if err != nil {
				t.Errorf("Parser returned error: %s", err.Error())
				failed = true
			}
			env := newValState()
			prog.eval(env)
			got := env.lookup("x")   // convention: test value stored in "x"
			if !got.equal(tt.want) { // custom equality check. vars can contain unused data after reassignment
				t.Errorf("x = %v, want %v", got, tt.want)
				failed = true
			}
			if failed {
				t.Log("Code:", tt.code)
			}
		})
	}
}
