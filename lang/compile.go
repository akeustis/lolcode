package lang

type namespace struct {
	vars map[string]interface{}
}

func (ns *namespace) getOrPanic(ident string) interface{} {
	if v, ok := ns.vars[ident]; ok {
		return v
	}
	panic("Reference to undefined variable: " + ident)
}

func (ns *namespace) putOrPanic(ident string, val interface{}) {
	if _, ok := ns.vars[ident]; !ok {
		panic("Assignment to undefined variable: " + ident)
	}
	ns.vars[ident] = val
}

type statement func(*namespace)
type expr func(*namespace) interface{}
type pred func(string, *namespace)

func nilFunc(args []interface{}) interface{} {
	return nil
}

func varPredicate(args []interface{}) interface{} {
	ident := args[0].(string)
	pred := args[1].(pred)
	return func(ns *namespace) {
		pred(ident, ns)
	}
}

func emptyPredicate(args []interface{}) interface{} {
	return func(ident string, ns *namespace) {
		ns.vars["IT"] = ns.getOrPanic(ident)
	}
}

func rExpr(args []interface{}) interface{} {
	expr := args[1].(expr)
	return func(ident string, ns *namespace) {
		ns.putOrPanic(ident, expr(ns))
	}
}

func ihasaVarItz(args []interface{}) interface{} {
	ident := args[1].(string)
	if args[2] == nil { // ITZ is optional
		return func(ns *namespace) {
			ns.vars[ident] = nil
		}
	}
	expr := args[2].(expr)
	return func(ns *namespace) {
		ns.vars[ident] = expr(ns)
	}
}

func itzExpr(args []interface{}) interface{} {
	return args[1]
}

func bareExpr(args []interface{}) interface{} {
	expr := args[0].(expr)
	return func(ns *namespace) {
		ns.vars["IT"] = expr(ns)
	}
}

func literal(args []interface{}) interface{} {
	return func(ns *namespace) interface{} {
		return args[0]
	}
}

func ident(args []interface{}) interface{} {
	ident := args[0].(string)
	return func(ns *namespace) interface{} {
		return ns.getOrPanic(ident)
	}
}

func sumofXAnY(args []interface{}) interface{} {
	return nil
}

func diffofXAnY(args []interface{}) interface{} {
	return nil
}

func prodofXAnY(args []interface{}) interface{} {
	return nil
}

func quoshofXAnY(args []interface{}) interface{} {
	return nil
}
