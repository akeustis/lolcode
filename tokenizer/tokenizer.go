package tokenizer

import (
	"bufio"
	"strconv"
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
	frags := make(chan string, 100)
	go emitFragments(reader, frags)
	for word := range frags {
		// Emit as many phrase tokens as possible
		for ok := true; ok; {
			word, ok = parsePhraseToken(word, frags, out)
		}
		switch {
		case word == "": // end of input
			close(out)
			return
		case word == "WIN": // TROOF literal
			out <- Token{TokLiteral, true}
		case word == "FAIL":
			out <- Token{TokLiteral, false}
		case word == "NOOB": // NOOB is a literal; casting to type NOOB is not allowed
			out <- Token{TokLiteral, nil}
		case word[0] == '"': // yarn literal
			out <- yarnLiteralToToken(word)
		case isIdentifier(word):
			out <- Token{TokIdent, word}
		default:
			if numbr, err := strconv.ParseInt(word, 0, 64); err == nil {
				out <- Token{TokLiteral, numbr}
				continue
			}
			if numbar, err := strconv.ParseFloat(word, 64); err == nil {
				out <- Token{TokLiteral, numbar}
				continue
			}
			out <- Token{TokErr, "Syntax error: unexpected token " + word}
		}
	}
}

// Reads a phrase starting with the given fragment (word)
// uses single-word look-ahead to parse as long a phrase as possible
func parsePhraseToken(word string, frags <-chan string, out chan<- Token) (string, bool) {
	phraseNode := phraseRoot
	hasRead := false
	for {
		nextNode := phraseNode.nodes[word]
		if nextNode == nil {
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
		phraseNode = nextNode
		word = <-frags
	}
}

func isIdentifier(s string) bool {
	if len(s) == 0 || !isLetter(s[0]) {
		return false
	}
	for i := 1; i < len(s); i++ {
		if l := s[i]; !isLetter(l) && !isNumberOrUnderscore(l) {
			return false
		}
	}
	return true
}

func isLetter(l byte) bool {
	return l >= 'a' && l <= 'z' || l >= 'A' && l <= 'Z'
}

func isNumberOrUnderscore(l byte) bool {
	return l >= '0' && l <= '9' || l == '_'
}

// TODO: add string escaping
func yarnLiteralToToken(str string) Token {
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
