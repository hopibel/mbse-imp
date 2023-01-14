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
		fmt.Printf("assign eval fail: tried to assign %s to %s\n", showValType(v), showValType(s.lookup(x)))
	}
}

func (stmt Seq) eval(s ValState) {
	stmt[0].eval(s)
	stmt[1].eval(s)
}

func (ite IfThenElse) eval(s ValState) {
	v := ite.cond.eval(s)
	if v.flag == ValueBool {
		s.startBlock()
		if v.valB {
			ite.thenStmt.eval(s)
		} else {
			ite.elseStmt.eval(s)
		}
		s.endBlock()
	} else {
		fmt.Printf("if-then-else eval fail: condition has type %s instead of boolean\n", showValType(v))
	}

}

func (e While) eval(s ValState) {
	v := e.cond.eval(s)
	if v.flag != ValueBool {
		fmt.Printf("while eval fail: condition has type %s instead of boolean\n", showValType(v))
		return
	}
	// evaluate body in a new scope as long as condition holds
	for v.valB {
		s.startBlock()
		e.body.eval(s)
		s.endBlock()
		v = e.cond.eval(s)
	}
}

func (e Print) eval(s ValState) {
	x := e.exp.eval(s)
	fmt.Println(showVal(x))
}

// Expressions

func (x Var) eval(s ValState) Val {
	return s.lookup(string(x))
}

func (x Bool) eval(s ValState) Val {
	return mkBool((bool)(x))
}

func (x Num) eval(s ValState) Val {
	return mkInt((int)(x))
}

func (e Equal) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == n2.flag && n1.flag != Undefined {
		switch n1.flag {
		case ValueBool:
			return mkBool(n1.valB == n2.valB)
		case ValueInt:
			return mkBool(n1.valI == n2.valI)
		}
	}
	return mkUndefined()
}

func (e Less) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkBool(n1.valI < n2.valI)
	}
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
	if b1.flag == ValueBool {
		// short circuit: false && _ => false
		if !b1.valB {
			return mkBool(false)
		}
		b2 := e[1].eval(s)
		if b2.flag == ValueBool {
			// true && V => V
			return mkBool(b2.valB)
		}
	}
	return mkUndefined()
}

func (e Or) eval(s ValState) Val {
	b1 := e[0].eval(s)
	if b1.flag == ValueBool {
		// short circuit: true || _ => true
		if b1.valB {
			return mkBool(true)
		}
		b2 := e[1].eval(s)
		if b2.flag == ValueBool {
			// false || V => V
			return mkBool(b2.valB)
		}
	}
	return mkUndefined()
}

func (e Not) eval(s ValState) Val {
	val := e.exp.eval(s)
	if val.flag == ValueBool {
		return mkBool(!val.valB)
	}
	return mkUndefined()
}
