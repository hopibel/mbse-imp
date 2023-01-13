package main

import (
	"strconv"
	"strings"
)

// Values

type Kind int

const (
	ValueInt  Kind = 0
	ValueBool Kind = 1
	Undefined Kind = 2
)

type Val struct {
	flag Kind
	valI int
	valB bool
}

func mkInt(x int) Val {
	return Val{flag: ValueInt, valI: x}
}
func mkBool(x bool) Val {
	return Val{flag: ValueBool, valB: x}
}
func mkUndefined() Val {
	return Val{flag: Undefined}
}

func showVal(v Val) string {
	var s string
	switch {
	case v.flag == ValueInt:
		s = Num(v.valI).pretty()
	case v.flag == ValueBool:
		s = Bool(v.valB).pretty()
	case v.flag == Undefined:
		s = "Undefined"
	}
	return s
}

// Types

type Type int

const (
	TyIllTyped Type = 0
	TyInt      Type = 1
	TyBool     Type = 2
)

func showType(t Type) string {
	var s string
	switch {
	case t == TyInt:
		s = "Int"
	case t == TyBool:
		s = "Bool"
	case t == TyIllTyped:
		s = "Illtyped"
	}
	return s
}

// Scopes and Environment (ValState)
// Scope is a mapping from variable names to values
// ValState is a stack of multiple Scopes
type Scope map[string]Val
type ValState []Scope

func newValState() ValState {
	return ValState{make(Scope)}
}

// lookup() traverses the stack and returns the first mapping found or undefined
func (env ValState) lookup(name string) Val {
	for i := len(env) - 1; i >= 0; i-- {
		if val, ok := env[i][name]; ok {
			return val
		}
	}
	return mkUndefined()
}

// declare new value mapping in current scope
// masks previous declarations in outer scopes until current scope ends
// overwrites previous declarations in same scope
func (env ValState) declare(name string, val Val) {
	env[len(env)-1][name] = val
}

// assign new value to existing mapping
// returns false if types don't match
func (env ValState) assign(name string, new_val Val) bool {
	for i := len(env) - 1; i >= 0; i-- {
		val, ok := env[i][name]
		if ok && val.flag == new_val.flag {
			env[i][name] = new_val
			return true
		}
	}
	return false
}

// push a new scope onto the environment stack (ValState)
// TODO: confirm this does what i think it does
func (env *ValState) startBlock() {
	*env = append(*env, make(Scope))
}

// pop top-most scope from the environment stack unless it is the global (bottom) scope
func (env *ValState) endBlock() {
	*env = (*env)[:len(*env)-1]
}

// Value State is a mapping from variable names to values
// type ValState map[string]Val

// TyScope is a mapping from variable names to types
// TyState is a stack of multiple TyScopes
type TyScope map[string]Type
type TyState []TyScope

func newTyState() TyState {
	return TyState{make(TyScope)}
}

// lookup() traverses the stack and returns the first mapping found or TyIllTyped
func (t TyState) lookup(name string) Type {
	for i := len(t) - 1; i >= 0; i-- {
		if ty, ok := t[i][name]; ok {
			return ty
		}
	}
	return TyIllTyped
}

// declare new type mapping in current scope
// masks previous declarations in outer scopes until current scope ends
// overwrites previous declarations in same scope
func (t TyState) declare(name string, ty Type) {
	t[len(t)-1][name] = ty
}

// push a new scope onto the type environment stack (TyState)
// TODO: confirm this does what i think it does
func (ts *TyState) startBlock() {
	*ts = append(*ts, make(TyScope))
}

// pop top-most scope from the type environment stack unless it is the global (bottom) scope
func (ts *TyState) endBlock() {
	*ts = (*ts)[:len(*ts)-1]
}

// Interface

type Exp interface {
	pretty() string
	eval(s ValState) Val
	infer(t TyState) Type
}

type Stmt interface {
	pretty() string
	eval(s ValState)
	check(t TyState) bool
}

// Statement cases

type Seq [2]Stmt
type Program Stmt
type Decl struct {
	lhs string
	rhs Exp
}
type Assign struct {
	lhs string
	rhs Exp
}
type While struct {
	cond Exp
	body Stmt
}
type IfThenElse struct {
	cond     Exp
	thenStmt Stmt
	elseStmt Stmt
}
type Print struct {
	exp Exp
}

// Expression cases

type Num int
type Bool bool
type Plus [2]Exp
type Mult [2]Exp
type Or [2]Exp
type And [2]Exp
type Not Exp
type Equal [2]Exp
type Less [2]Exp
type Var string

/////////////////////////
// Stmt instances

// pretty print

func (stmt Seq) pretty() string {
	ret := stmt[0].pretty() + ";\n" + stmt[1].pretty()
	if ret[len(ret)-1] != ';' {
		ret += ";"
	}
	return ret
}

func (decl Decl) pretty() string {
	return decl.lhs + " := " + decl.rhs.pretty()
}

func (assign Assign) pretty() string {
	return assign.lhs + " = " + assign.rhs.pretty()
}

func (while While) pretty() string {
	ret := "while " + while.cond.pretty() + " {\n" +
		"\t" + strings.ReplaceAll(while.body.pretty(), "\n", "\n\t")
	if ret[len(ret)-1] != ';' {
		ret += ";"
	}
	return ret + "\n}"
}

func (ite IfThenElse) pretty() string {
	ret := "if " + ite.cond.pretty() + " {\n" +
		"\t" + strings.ReplaceAll(ite.thenStmt.pretty(), "\n", "\n\t")
	if ret[len(ret)-1] != ';' {
		ret += ";"
	}
	ret += "\n} " + "else {\n" +
		"\t" + strings.ReplaceAll(ite.elseStmt.pretty(), "\n", "\n\t")
	if ret[len(ret)-1] != ';' {
		ret += ";"
	}
	return ret + "\n}"
}

func (print Print) pretty() string {
	return "print " + print.exp.pretty()
}

/////////////////////////
// Exp instances

// pretty print

func (x Var) pretty() string {
	return (string)(x)
}

func (x Bool) pretty() string {
	if x {
		return "true"
	} else {
		return "false"
	}

}

func (x Num) pretty() string {
	return strconv.Itoa(int(x))
}

func (e Equal) pretty() string {
	return "(" + e[0].pretty() + "==" + e[1].pretty() + ")"
}

func (e Less) pretty() string {
	return "(" + e[0].pretty() + "<" + e[1].pretty() + ")"
}

func (e Mult) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "*"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Plus) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "+"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e And) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "&&"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Or) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "||"
	x += e[1].pretty()
	x += ")"

	return x
}
