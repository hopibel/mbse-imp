package main

import (
	"strconv"
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

// Scopes and Environment
// Scope is a mapping from variable names to values
// Env is a stack of multiple Scopes
type Scope map[string]Val
type Env []Scope

func newEnv() Env {
	return Env{make(Scope)}
}

// lookup() traverses the stack and returns the first mapping found or undefined
func (env Env) lookup(name string) Val {
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
func (env Env) declare(name string, val Val) {
	env[len(env)-1][name] = val
}

// assign new value to existing mapping
// returns false if types don't match
func (env Env) assign(name string, new_val Val) bool {
	for i := len(env) - 1; i >= 0; i-- {
		val, ok := env[i][name]
		if ok && val.flag == new_val.flag {
			env[i][name] = new_val
			return true
		}
	}
	return false
}

// TODO: needs to be stack for local scopes
// Value State is a mapping from variable names to values
type ValState map[string]Val

// Value State is a mapping from variable names to types
type TyState map[string]Type

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

// Statement cases (incomplete)

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
	return stmt[0].pretty() + "; " + stmt[1].pretty()
}

func (decl Decl) pretty() string {
	return decl.lhs + " := " + decl.rhs.pretty()
}

func (assign Assign) pretty() string {
	return assign.lhs + " = " + assign.rhs.pretty()
}

func (while While) pretty() string {
	panic("not yet implemented")
	return ""
}

func (ite IfThenElse) pretty() string {
	panic("not yet implemented")
	return ""
}

func (print Print) pretty() string {
	panic("not yet implemented")
	return ""
}

// type check

func (stmt Seq) check(t TyState) bool {
	if !stmt[0].check(t) {
		return false
	}
	return stmt[1].check(t)
}

func (decl Decl) check(t TyState) bool {
	ty := decl.rhs.infer(t)
	if ty == TyIllTyped {
		return false
	}

	x := (string)(decl.lhs)
	t[x] = ty
	return true
}

func (a Assign) check(t TyState) bool {
	x := (string)(a.lhs)
	return t[x] == a.rhs.infer(t)
}

func (w While) check(t TyState) bool {
	panic("not yet implemented")
	return false
}

func (ite IfThenElse) check(t TyState) bool {
	panic("not yet implemented")
	return false
}

func (print Print) check(t TyState) bool {
	panic("not yet implemented")
	return false
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
	panic("not yet implemented")
	return ""
}

func (e Less) pretty() string {
	panic("not yet implemented")
	return ""
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
