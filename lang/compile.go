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

// pred
func emptyPredicate(args []interface{}) interface{} {
	return pred(func(ident string, ns *namespace) {
		ns.vars["IT"] = ns.getOrPanic(ident)
	})
}

func rExpr(args []interface{}) interface{} {
	expr := args[1].(expr)
	return pred(func(ident string, ns *namespace) {
		ns.putOrPanic(ident, expr(ns))
	})
}

// statement
func varPredicate(args []interface{}) interface{} {
	ident := args[0].(string)
	pred := args[1].(pred)
	return statement(func(ns *namespace) {
		pred(ident, ns)
	})
}

func ihasaVarItz(args []interface{}) interface{} {
	ident := args[1].(string)
	if args[2] == nil { // ITZ is optional
		return statement(func(ns *namespace) {
			ns.vars[ident] = nil
		})
	}
	expr := args[2].(expr)
	return statement(func(ns *namespace) {
		ns.vars[ident] = expr(ns)
	})
}

// Expressions
func itzExpr(args []interface{}) interface{} {
	return args[1]
}

func bareExpr(args []interface{}) interface{} {
	expr := args[0].(expr)
	return statement(func(ns *namespace) {
		ns.vars["IT"] = expr(ns)
	})
}

func literal(args []interface{}) interface{} {
	return expr(func(ns *namespace) interface{} {
		return args[0]
	})
}

func ident(args []interface{}) interface{} {
	ident := args[0].(string)
	return expr(func(ns *namespace) interface{} {
		return ns.getOrPanic(ident)
	})
}

// If x is int64 and y is float64, promote x to float64 (and vice versa)
func saem(x, y interface{}) bool {
	switch v := x.(type) {
	case float64:
		switch w := y.(type) {
		case int64:
			y = float64(w)
		}
	case int64:
		switch y.(type) {
		case float64:
			x = float64(v)
		}
	}
	return x == y
}

func bothsaemXAnY(args []interface{}) interface{} {
	x, y := args[1].(expr), args[3].(expr)
	return expr(func(ns *namespace) interface{} {
		return saem(x(ns), y(ns))
	})
}

func diffrintXAnY(args []interface{}) interface{} {
	x, y := args[1].(expr), args[3].(expr)
	return expr(func(ns *namespace) interface{} {
		return !saem(x(ns), y(ns))
	})
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

func modofXAnY(args []interface{}) interface{} {
	return makeMathExpr(args[1].(expr), args[3].(expr),
		func(a, b int64) int64 { return a % b },
		func(a, b float64) float64 { panic("Cannot use MOD OF with type NUMBAR") },
	)
}

func castToYarn(x interface{}, isExplicit bool) string {
	switch x := x.(type) {
	case nil:
		if isExplicit {
			return ""
		}
		panic("Cannot implicitly cast NOOB to YARN")
	case bool:
		if x {
			return "WIN"
		}
		return "FAIL"
	default:
		return fmt.Sprintf("%v", x)
	}
}

func troof(x interface{}) bool {
	switch x := x.(type) {
	case bool:
		return x
	case nil:
		return false
	case string:
		return x != ""
	case int64:
		return x != 0
	case float64:
		return x != 0
	default:
		panic(fmt.Sprintf("Cannot cast type %t to TROOF", x))
	}
}

func exprMoar(args []interface{}) interface{} {
	rest := args[1].([]interface{})
	list := make([]expr, len(rest)+1)
	list[0] = args[0].(expr)
	for i, e := range rest {
		list[i+1] = e.(expr)
	}
	return list
}

func anExpr(args []interface{}) interface{} {
	return args[1]
}

func notExpr(args []interface{}) interface{} {
	e := args[1].(expr)
	return expr(func(ns *namespace) interface{} {
		return !troof(e(ns))
	})
}

func bothofXAnY(args []interface{}) interface{} {
	x, y := args[1].(expr), args[3].(expr)
	return expr(func(ns *namespace) interface{} {
		return troof(x(ns)) && troof(y(ns))
	})
}

func eitherofXAnY(args []interface{}) interface{} {
	x, y := args[1].(expr), args[3].(expr)
	return expr(func(ns *namespace) interface{} {
		return troof(x(ns)) || troof(y(ns))
	})
}

func wonofXAnY(args []interface{}) interface{} {
	x, y := args[1].(expr), args[3].(expr)
	return expr(func(ns *namespace) interface{} {
		return !troof(x(ns)) == troof(y(ns))
	})
}

func allofList(args []interface{}) interface{} {
	exprs := args[1].([]expr)
	return expr(func(ns *namespace) interface{} {
		for _, e := range exprs {
			if !troof(e(ns)) {
				return false
			}
		}
		return true
	})
}

func anyofList(args []interface{}) interface{} {
	exprs := args[1].([]expr)
	return expr(func(ns *namespace) interface{} {
		for _, e := range exprs {
			if troof(e(ns)) {
				return true
			}
		}
		return false
	})
}
