package tokenizer

import (
	"bufio"
	"strings"
	"testing"
)

const lolCode = `HAI 1.2
tok0,tok1 BTW comments here, including some commas
tok2

	BTW full line comment
tok3 BTW, OBTW doesnt work here
tok4	 OBTW   illegal comment
tok5, OBTW legal comment,, TLDR tok6,,,
tok7, OBTW legal comment TLDR
BTW that TLDR doesn't work lulz, TLDR BTW
KTHXBYE
`

func TestEmitFragments(t *testing.T) {
	expected := []string{"HAI", "1.2", "\n",
		"tok0", "\n", "tok1", "\n",
		"tok2", "\n",
		"tok3", "\n",
		"tok4", "OBTW", "illegal", "comment", "\n",
		"tok5", "\n",
		"tok6", "\n",
		"tok7", "\n",
		"KTHXBYE", "\n"}
	reader := bufio.NewReader(strings.NewReader(lolCode))
	fragments := make(chan string, 100)
	go emitFragments(reader, fragments)
	i := 0
	for fragment := range fragments {
		if fragment != expected[i] {
			t.Fatalf("Expected: %s Got: %s", expected[i], fragment)
		}
		i++
	}
	if i != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), i)
	}
}

const lolCode2 = `HAI 1.2
I HAS A FISH ITZ 5
BTW full line comment
OBTW,TLDR 
FISH R "foo"
KTHXBYE

`

func TestEmitTokens(t *testing.T) {
	expected := []Token{
		{TokHAI, nil}, {TokLiteral, float32(1.2)}, {TokEOL, nil},
		{TokIHASA, nil}, {TokIdent, "FISH"}, {TokITZ, nil}, {TokLiteral, int(5)}, {TokEOL, nil},
		{TokIdent, "FISH"}, {TokR, nil}, {TokLiteral, "foo"}, {TokEOL, nil},
		{TokKTHXBYE, nil}, {TokEOL, nil},
	}
	reader := bufio.NewReader(strings.NewReader(lolCode2))
	fragments := make(chan Token, 100)
	go EmitTokens(reader, fragments)
	i := 0
	for fragment := range fragments {
		if fragment != expected[i] {
			t.Fatalf("Expected: %v Got: %v", expected[i], fragment)
		}
		i++
	}
	if i != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), i)
	}
}
