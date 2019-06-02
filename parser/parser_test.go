package parser

import (
	"bufio"
	"lol/tokenizer"
	"strings"
	"testing"
)

const (
	Expr = iota + tokenizer.NumTokens
	ExprList
	Sum
	Prod
	NumNodes
)

var d = NewDialect(tokenizer.NumTokens, NumNodes)

func init() {
	// Define a basic calculator grammar for testing
	fetchFirst := func(args []interface{}) interface{} { return args[0] }
	fetchSecond := func(args []interface{}) interface{} { return args[1] }
	sum := func(args []interface{}) interface{} {
		sum := int64(0)
		for _, a := range args[1].([]interface{}) {
			sum += a.(int64)
		}
		return sum
	}
	prod := func(args []interface{}) interface{} {
		prod := int64(1)
		for _, a := range args[1].([]interface{}) {
			prod *= a.(int64)
		}
		return prod
	}
	d.Rule(Expr, fetchFirst, tokenizer.TokLiteral)
	d.Rule(Expr, sum, tokenizer.TokSUMOF, ExprList, -tokenizer.TokMKAY)
	d.Rule(Expr, prod, tokenizer.TokPRODUKTOF, ExprList, -tokenizer.TokMKAY)
	d.RepRule(ExprList, fetchSecond, -tokenizer.TokAN, Expr)
}

func TestMath(t *testing.T) {
	str := "SUM OF 3 4 AN 5\n"
	reader := bufio.NewReader(strings.NewReader(str))
	tokens := make(chan tokenizer.Token, 100)
	go tokenizer.EmitTokens(reader, tokens)
	cur, val, ok := d.Parse(Expr, tokens)
	switch {
	case !ok:
		t.Fatalf("Parse unsuccessful")
	case *cur != tokenizer.Token{Type: tokenizer.TokEOL, Value: nil}:
		t.Fatalf("Expected token mismatch")
	case val.(int64) != 12:
		t.Fatalf("Parse returned wrong value")
	}
}
