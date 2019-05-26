package tokenizer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TokenType uint
type Token struct {
	Type  TokenType
	Value interface{}
}

const (
	TokHAI     = iota
	TokKTHXBYE = iota
	TokLiteral = iota
	TokIdent   = iota
	TokEOL     = iota
	TokIHASA   = iota
	TokITZ     = iota
)

var phraseStore = make(map[string][]string)

func keyPhrase(t TokenType, phr string) {

}
func phrases() {
	keyPhrase(TokHAI, "HAI")
	keyPhrase(TokKTHXBYE, "KTHXBYE")
	keyPhrase(TokIHASA, "I HAS A")
	keyPhrase(TokITZ, "ITZ")
}

// emitFragments reads from a bufio.Reader and emits string fragments on the given channel,
// omitting all comments and converting line separators "," into "\n" fragments
func emitFragments(reader *bufio.Reader, out chan<- string) {
	defer close(out)
	insideComment := false
txtLine:
	for {
		txt, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		for _, line := range strings.Split(txt, ",") {
			fragments := strings.Fields(line)
			if len(fragments) == 0 {
				continue //ignore empty lines
			}
			i := 0
			switch fragments[0] {
			// OBTW and TLDR are required to be the first token on a Lolcode line, per spec
			case "OBTW":
				insideComment = true
			case "TLDR":
				insideComment = false
				i = 1 // Resume processing after the TLDR
			case "BTW":
				continue txtLine // Don't emit an extra "\n" for full-line comment
			}
			if !insideComment {
				for ; i < len(fragments); i++ {
					if fragments[i] == "BTW" {
						out <- "\n"
						continue txtLine
					}
					out <- fragments[i]
				}
				out <- "\n"
			}
		}
	}
}

// Parse echoes each line of stdin back out
func Parse() {
	phrases()
	stdin := bufio.NewReader(os.Stdin)
	for {
		text, err := stdin.ReadString('\n')
		if err != nil {
			return
		}
		tokens := strings.Fields(text)
		for _, token := range tokens {
			fmt.Printf("%s,", token)
		}
		fmt.Println()
	}
}
