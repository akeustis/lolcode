package tokenizer

import (
	"bufio"
	"strings"
	"testing"
)

func TestEmitFragments(t *testing.T) {
	const lolCode = `HAI 1.2
	tok1 BTW comments here, including some commas
	tok2
	BTW, full line comment
	tok3 BTW
	tok4	 OBTW   illegal comment
	tok5, OBTW legal comment TLDR
	BTW that TLDR didn't work lulz, TLDR
	that didn't work either lulz, OBTW, TLDR tok6
	KTHXBYE
	`
	expected := []string{"HAI", "1.2", "\n",
		"tok1", "\n",
		"tok2", "\n",
		"tok3", "\n",
		"tok4", "OBTW", "illegal", "comment", "\n",
		"tok5", "\n",
		"tok6", "\n",
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
