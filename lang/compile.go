package lang

import (
	"fmt"
	"strconv"
)

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

func getNumericValue(val interface{}, useFloat bool) interface{} {
	typePanicMsg := "Cannot perform numerical operation on type "
	switch v := val.(type) {
	case nil:
		panic(typePanicMsg + "NOOB")
	case bool:
		panic(typePanicMsg + "TROOF")
	case int64:
		if useFloat {
			return float64(v)
		}
		return v
	case float64:
		return v
	case string:
		if !useFloat {
			if i, err := strconv.ParseInt(v, 0, 64); err == nil {
				return i
			}
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
		panic("Failed to parse numeric value from string: " + v)
	default:
		panic(typePanicMsg + fmt.Sprintf("%T", v))
	}
}

func makeMathExpr(left, right expr, intOper func(int64, int64) int64,
	floatOper func(float64, float64) float64) expr {
	return func(ns *namespace) interface{} {
		switch v1 := getNumericValue(left(ns), false).(type) {
		case int64:
			switch v2 := getNumericValue(right(ns), false).(type) {
			case int64:
				return intOper(v1, v2)
			default: // float64
				return floatOper(float64(v1), v2.(float64))
			}
		default: // float64
			v2 := getNumericValue(right(ns), true)
			return floatOper(v1.(float64), v2.(float64))
		}
	}
}

func sumofXAnY(args []interface{}) interface{} {
	return makeMathExpr(args[1].(expr), args[3].(expr),
		func(a, b int64) int64 { return a + b },
		func(a, b float64) float64 { return a + b },
	)
}

func diffofXAnY(args []interface{}) interface{} {
	return makeMathExpr(args[1].(expr), args[3].(expr),
		func(a, b int64) int64 { return a - b },
		func(a, b float64) float64 { return a - b },
	)
}

func prodofXAnY(args []interface{}) interface{} {
	return makeMathExpr(args[1].(expr), args[3].(expr),
		func(a, b int64) int64 { return a * b },
		func(a, b float64) float64 { return a * b },
	)
}

func quoshofXAnY(args []interface{}) interface{} {
	return makeMathExpr(args[1].(expr), args[3].(expr),
		func(a, b int64) int64 { return a / b },
		func(a, b float64) float64 { return a / b },
	)
}
