package tokenizer

import (
	"bufio"
	"strings"
)

type TokenType uint
type Token struct {
	Type  TokenType
	Value interface{}
}

// Exported token types
const (
	TokErr = iota
	TokHAI
	TokKTHXBYE
	TokLiteral
	TokIdent
	TokEOL
	TokIHASA
	TokITZ
	TokR
	TokIIZ
)

type phraseNode struct {
	t     TokenType
	nodes map[string]*phraseNode
}

func newPhraseNode() *phraseNode {
	return &phraseNode{TokErr, make(map[string]*phraseNode)}
}

type phraseInit struct {
	t      TokenType
	phrase string
}

func initPhrases(phraseInits []phraseInit) *phraseNode {
	phraseRoot := newPhraseNode()
	for _, init := range phraseInits {
		initPhrase(phraseRoot, init.t, strings.Split(init.phrase, " "))
	}
	return phraseRoot
}

// Recursively build phrase tree
func initPhrase(root *phraseNode, t TokenType, words []string) {
	if len(words) == 0 {
		root.t = t
		return
	}
	word := words[0]
	node := root.nodes[word]
	if node == nil {
		node = newPhraseNode()
		root.nodes[word] = node
	}
	initPhrase(node, t, words[1:])
}

var phraseRoot = initPhrases([]phraseInit{
	{TokEOL, "\n"},
	{TokHAI, "HAI"},
	{TokKTHXBYE, "KTHXBYE"},
	{TokIHASA, "I HAS A"},
	{TokITZ, "ITZ"},
	{TokR, "R"},
	{TokIIZ, "I IZ"},
})

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
			start := 0
			if insideComment {
				if fragments[0] != "TLDR" {
					continue
				}
				insideComment = false
				if len(fragments) == 1 {
					continue
				}
				start = 1
			} else {
				if fragments[0] == "OBTW" {
					insideComment = true
					continue
				}
			}
			for i := start; i < len(fragments); i++ {
				if fragments[i] == "BTW" {
					if i > start {
						out <- "\n"
					}
					continue txtLine
				}
				out <- fragments[i]
			}
			out <- "\n"
		}
	}
}

// EmitTokens parses a Lolcode reader and emits a stream of Tokens
func EmitTokens(reader *bufio.Reader, out chan<- Token) {
	defer close(out)
	frags := make(chan string, 100)
	go emitFragments(reader, frags)
	for word := range frags {
		// Emit as many phrase tokens as we can
		if word, ok := emitPhraseToken(word, frags, out); !ok {
			switch {
			case word == "": // end of input
				return
			case word[0] == '"': // string literal
				out <- stringLiteralToToken(word)
			}
		}
	}
}

// Reads a phrase starting with the given fragment (word)
// uses single-word look-ahead to parse as long a phrase as possible
func emitPhraseToken(word string, frags <-chan string, out chan<- Token) (string, bool) {
	phraseNode := phraseRoot
	hasRead := false
	for {
		if phraseNode = phraseNode.nodes[word]; phraseNode == nil {
			if hasRead {
				token := Token{phraseNode.t, nil}
				if token.Type == TokErr {
					// If we have an TokError, fill in a parser error as the value
					token.Value = getErrMessageForPhrase(phraseNode, word)
				}
				out <- token
				return word, true
			} //else
			return word, false
		}
		hasRead = true
		word = <-frags
	}
}

// TODO: add string escaping
func stringLiteralToToken(str string) Token {
	// String literal must end with '"' and have length at least 2
	if l := len(str); l < 2 || str[l-1] != '"' {
		return Token{TokErr, "Invalid string literal: " + str}
	}
	// chop off the start and end quotes
	return Token{TokLiteral, str[1 : len(str)-1]}
}

func getErrMessageForPhrase(node *phraseNode, word string) string {
	var out strings.Builder
	w := out.WriteString
	w("Syntax error at ")
	w(word)
	w(": expected ")
	sep := ""
	for word := range node.nodes {
		w(sep)
		w(word)
		sep = " or "
	}
	return out.String()
}
