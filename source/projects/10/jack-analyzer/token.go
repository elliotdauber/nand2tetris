package main

type Token struct {
	val string
	typ TokenType
}

func (this *Token) Init(val string, typ TokenType) Token {
	this.val = val
	this.typ = typ
	return *this
}

type TokenType int

const (
	Keyword TokenType = iota
	Symbol
	Identifier
	IntConst
	StrConst
)

type KeyWord int

const (
	_Class KeyWord = iota
	Method
	Function
	Constructor
	Int
	Boolean
	Char
	Void
	Var
	Static
	Field
	Let
	Do
	If
	Else
	While
	Return
	True
	False
	Null
	This
)
