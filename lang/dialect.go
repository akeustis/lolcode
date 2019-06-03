package lang

import (
	"lol/parser"
	"lol/token"
)

// ids of grammar nodes
const (
	Expr = iota + token.NumTokens
	Statement
	VarPredicate
	Itz
	NumNodes
)

// Dialect is The *token.Dialect which implements the Lolcode language
var D = parser.NewDialect(token.NumTokens, NumNodes)

func init() {
	// Statement
	D.Rule(Statement, varPredicate, token.Ident, VarPredicate)
	D.Rule(Statement, ihasaVarItz, token.IHASA, token.Ident, -Itz, token.EOL)
	D.Rule(Statement, bareExpr, Expr, token.EOL)

	// Itz
	D.Rule(Itz, itzExpr, token.ITZ, Expr)

	//VarPredicate
	D.Rule(VarPredicate, emptyPredicate, token.EOL)
	D.Rule(VarPredicate, rExpr, token.R, Expr)

	// Expr
	D.Rule(Expr, literal, token.Literal)
	D.Rule(Expr, ident, token.Ident)
	D.Rule(Expr, sumofXAnY, token.SUMOF, Expr, token.AN, Expr)
	D.Rule(Expr, diffofXAnY, token.DIFFOF, Expr, token.AN, Expr)
	D.Rule(Expr, prodofXAnY, token.PRODUKTOF, Expr, token.AN, Expr)
	D.Rule(Expr, quoshofXAnY, token.QUOSHUNTOF, Expr, token.AN, Expr)
}
