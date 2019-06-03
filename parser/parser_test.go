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
	type testCase struct {
		code     string
		expected int64
	}
	testCases := []testCase{
		// 3+4+5
		{"SUM OF 3 4 AN 5\n", 12},
		// (3+4)*5
		{"PRODUKT OF SUM OF 3 AN 4 MKAY AN 5\n", 35},
		// 3 + (4*5) + ((2+1)*3)
		{"SUM OF 3 AN PRODUKT OF 4 AN 5 MKAY AN PRODUKT OF SUM OF 2 1 MKAY 3\n", 32},
		{"SUM OF MKAY\n", 0},
		{"PRODUKT OF\n", 1},
	}
	for _, tc := range testCases {
		reader := bufio.NewReader(strings.NewReader(tc.code))
		tokens := make(chan tokenizer.Token, 100)
		go tokenizer.EmitTokens(reader, tokens)
		cur, val, ok := d.Parse(Expr, tokens)
		switch {
		case !ok:
			t.Fatalf("Parse unsuccessful")
		case *cur != tokenizer.Token{tokenizer.TokEOL, tokenizer.EOLPhrase}:
			t.Fatalf("Expected token mismatch")
		case val.(int64) != tc.expected:
			t.Fatalf("Parse returned %d, expected %d", val.(int64), tc.expected)
		}
	}
}

func TestPanic(t *testing.T) {
	str := "SUM OF 3 AN\n"
	reader := bufio.NewReader(strings.NewReader(str))
	tokens := make(chan tokenizer.Token, 100)
	go tokenizer.EmitTokens(reader, tokens)
	_, _, ok := d.Parse(Expr, tokens)
	if ok {
		t.Fatalf("Expected failure")
	}
}
