package main

import "fmt"

// Evaluator

// Statements

// Maps are represented via pointers.
// Hence, maps are passed by "reference" and the update is visible for the caller as well.
func (decl Decl) eval(s ValState) {
	x := (string)(decl.lhs)
	v := decl.rhs.eval(s)
	s.declare(x, v)
}

func (assign Assign) eval(s ValState) {
	x := (string)(assign.lhs)
	v := assign.rhs.eval(s)
	if !s.assign(x, v) {
		fmt.Printf("assign eval fail: type mismatch")
	}
}

func (stmt Seq) eval(s ValState) {
	stmt[0].eval(s)
	stmt[1].eval(s)
}

func (ite IfThenElse) eval(s ValState) {
	v := ite.cond.eval(s)
	if v.flag == ValueBool {
		if v.valB {
			ite.thenStmt.eval(s)
		} else {
			ite.elseStmt.eval(s)
		}
	} else {
		fmt.Printf("if-then-else eval fail")
	}

}

func (e While) eval(s ValState) {
	v := e.cond.eval(s)
	if v.flag != ValueBool {
		fmt.Printf("while eval fail: condition does not return a boolean")
		return
	}
	for v.valB {
		s.startBlock()
		e.body.eval(s)
		s.endBlock()
		v = e.cond.eval(s)
	}
}

func (e Print) eval(s ValState) {
	panic("not yet implemented")
}

// Expressions

func (x Var) eval(s ValState) Val {
	panic("not yet implemented")
	return mkUndefined()
}

func (x Bool) eval(s ValState) Val {
	return mkBool((bool)(x))
}

func (x Num) eval(s ValState) Val {
	return mkInt((int)(x))
}

func (e Equal) eval(s ValState) Val {
	panic("not yet implemented")
	return mkUndefined()
}

func (e Less) eval(s ValState) Val {
	panic("not yet implemented")
	return mkUndefined()
}

func (e Mult) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI * n2.valI)
	}
	return mkUndefined()
}

func (e Plus) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI + n2.valI)
	}
	return mkUndefined()
}

func (e And) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && !b1.valB:
		return mkBool(false)
	case b1.flag == ValueBool && b2.flag == ValueBool:
		return mkBool(b1.valB && b2.valB)
	}
	return mkUndefined()
}

func (e Or) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && b1.valB:
		return mkBool(true)
	case b1.flag == ValueBool && b2.flag == ValueBool:
		return mkBool(b1.valB || b2.valB)
	}
	return mkUndefined()
}
