package main

// Type inferencer/checker

// Expressions type inference

func (x Var) infer(t TyState) Type {
	return t.lookup(string(x))
}

func (x Bool) infer(t TyState) Type {
	return TyBool
}

func (x Num) infer(t TyState) Type {
	return TyInt
}

func (e Equal) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == t2 {
		return TyBool
	}
	return TyIllTyped
}

func (e Less) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyBool
	}
	return TyIllTyped
}

func (e Mult) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e Plus) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e And) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (e Or) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (e Not) infer(t TyState) Type {
	if e.exp.infer(t) == TyBool {
		return TyBool
	}
	return TyIllTyped
}

// Statement type checking

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
	t.declare(x, ty)
	return true
}

func (a Assign) check(t TyState) bool {
	x := (string)(a.lhs)
	return t.lookup(x) == a.rhs.infer(t)
}

func (w While) check(t TyState) bool {
	if w.cond.infer(t) != TyBool {
		return false
	}
	t.startBlock()
	b := w.body.check(t)
	t.endBlock()
	return b
}

func (ite IfThenElse) check(t TyState) bool {
	if ite.cond.infer(t) != TyBool {
		return false
	}
	t.startBlock()
	b := ite.thenStmt.check(t)
	t.endBlock()
	if !b {
		return false
	}
	t.startBlock()
	b = ite.elseStmt.check(t)
	t.endBlock()
	return b
}

func (print Print) check(t TyState) bool {
	return print.exp.infer(t) != TyIllTyped
}
