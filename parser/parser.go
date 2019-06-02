package parser

import "lol/tokenizer"

// Dialect is a collection of parser nodes that form a language
type Dialect struct {
	nodes   []parseNode
	names   []string
	numToks int
}

type parseNode struct {
	name  string
	rules []rule
}

// NewDialect constructs a Dialect.
// ids 0 through t-1 are reserved for single-token parsing.  tok.Value will be parsed up to higher-level functions.
// ids t through m-1 are reserved for grammar nodes.
// Caller is responsible for knowing ahead of time how many basic tokens and grammar nodes they want to have;
// use of const with iota for compile-time determination of t and m is highly recommended.
func NewDialect(t int, m int) *Dialect {
	return &Dialect{
		make([]parseNode, m-t),
		make([]string, m-t),
		t,
	}
}

// Name assigns a name to a given grammar node.  This name will be used to format syntax errors.
func (d *Dialect) Name(i int, name string) {
	d.names[i-d.numToks] = name
}

// Parser is the signature of the functions that must be supplied to Rule and RepRule
type Parser func(args []interface{}) interface{}

// Rule establishes a new parseRule for node i which will parse nodes off a stream
// as determined by the given args, then apply the given parser function.
func (d *Dialect) Rule(i int, p Parser, args ...int) {
	d.rule(i, false, p, args)
}

// RepRule is similar to Rule but the sequence of nodes will be parsed as many times as possible (0 is ok).
// p is applied to each cycle and a slice of results is forwarded up.
func (d *Dialect) RepRule(i int, p Parser, args ...int) {
	d.rule(i, true, p, args)
}

// Represents a rule by which a node may be parsed.  A single node allows multiple
// rules if they are distinguishable by their first token.
type rule struct {
	nodes       []int
	isRepeating bool
	parse       Parser
}

func (d *Dialect) rule(i int, isRepeating bool, p Parser, args []int) {
	node := &d.nodes[i-d.numToks]
	node.rules = append(node.rules, rule{
		args, isRepeating, p,
	})
}

// Parse will parse a channel of supplied tokens according to the rules of the dialect
func (d *Dialect) Parse(start int, tokens <-chan tokenizer.Token,
) (*tokenizer.Token, interface{}, bool) {
	first := <-tokens
	return d.parseNode(start, &first, tokens)
}

func (d *Dialect) parseNode(id int, curr *tokenizer.Token, more <-chan tokenizer.Token,
) (*tokenizer.Token, interface{}, bool) {
	// base case is a single token
	if id < d.numToks {
		if curr.Type == id {
			next := <-more
			return &next, curr.Value, true
		}
		return curr, nil, false
	}
	// recursively try each of the rules until we find a winner
	node := &d.nodes[id-d.numToks]
	for _, r := range node.rules {
		if curr, val, ok := d.parseRule(&r, curr, more); ok {
			return curr, val, true
		}
	}
	return curr, nil, false
}

// Attempt to parse the given rule
func (d *Dialect) parseRule(r *rule, curr *tokenizer.Token, more <-chan tokenizer.Token,
) (*tokenizer.Token, interface{}, bool) {
	if r.isRepeating {
		var result []interface{}
		for {
			curr, val, ok := d.parseRuleSingle(r, curr, more)
			if !ok {
				return curr, result, true
			}
			result = append(result, val)
		}
	}
	return d.parseRuleSingle(r, curr, more)
}

// Attempt to parse a single pass of the given rule
func (d *Dialect) parseRuleSingle(r *rule, curr *tokenizer.Token, more <-chan tokenizer.Token,
) (*tokenizer.Token, interface{}, bool) {
	var vals []interface{}
	for i := 0; i < len(r.nodes); i++ {
		id, optional := r.nodes[i], false
		if id < 0 {
			id = -id
			optional = true
		}
		curr, val, ok := d.parseNode(id, curr, more)
		switch {
		case ok:
			if vals == nil {
				vals = make([]interface{}, len(r.nodes))
			}
			vals[i] = val
		case !optional:
			if vals == nil {
				return curr, nil, false
			}
			panic(curr) // syntax error, we parsed one and failed another
		}
	}
	return curr, r.parse(vals), true
}
