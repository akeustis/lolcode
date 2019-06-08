package lang

import (
	"lol/parser"
	"lol/token"
)

// ids of grammar nodes
const (
	Expr = iota + token.NumTokens
	ExprList
	MoarList
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

	// ExprList
	D.Rule(ExprList, exprMoar, Expr, MoarList)
	// MoarList
	D.RepRule(MoarList, anExpr, -token.AN, Expr)

	// Expr
	// literal
	D.Rule(Expr, literal, token.Literal)
	// variable lookup
	D.Rule(Expr, ident, token.Ident)
	// boolean
	D.Rule(Expr, notExpr, token.NOT, Expr)
	D.Rule(Expr, bothofXAnY, token.BOTHOF, Expr, -token.AN, Expr)
	D.Rule(Expr, eitherofXAnY, token.EITHEROF, Expr, -token.AN, Expr)
	D.Rule(Expr, wonofXAnY, token.WONOF, Expr, -token.AN, Expr)
	D.Rule(Expr, allofList, token.ALLOF, ExprList, -token.MKAY)
	D.Rule(Expr, anyofList, token.ANYOF, ExprList, -token.MKAY)

	// comparison
	D.Rule(Expr, bothsaemXAnY, token.BOTHSAEM, Expr, -token.AN, Expr)
	D.Rule(Expr, diffrintXAnY, token.DIFFRINT, Expr, -token.AN, Expr)
	// math
	D.Rule(Expr, sumofXAnY, token.SUMOF, Expr, token.AN, Expr)
	D.Rule(Expr, diffofXAnY, token.DIFFOF, Expr, token.AN, Expr)
	D.Rule(Expr, prodofXAnY, token.PRODUKTOF, Expr, token.AN, Expr)
	D.Rule(Expr, quoshofXAnY, token.QUOSHUNTOF, Expr, token.AN, Expr)
	D.Rule(Expr, modofXAnY, token.MODOF, Expr, token.AN, Expr)
}
