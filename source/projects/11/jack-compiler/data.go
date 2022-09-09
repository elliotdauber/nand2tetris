package main

import mapset "github.com/deckarep/golang-set/v2"

func SymbolSet() mapset.Set[rune] {
	set := mapset.NewSet[rune]()
	set.Add('{')
	set.Add('}')
	set.Add('(')
	set.Add(')')
	set.Add('[')
	set.Add(']')
	set.Add('.')
	set.Add(',')
	set.Add(';')
	set.Add('+')
	set.Add('-')
	set.Add('*')
	set.Add('/')
	set.Add('&')
	set.Add('|')
	set.Add('<')
	set.Add('>')
	set.Add('=')
	set.Add('~')
	return set
}

func KeyWordMap() map[string]KeyWord {
	return map[string]KeyWord{
		"class":       _Class,
		"method":      Method,
		"function":    Function,
		"constructor": Constructor,
		"int":         Int,
		"boolean":     Boolean,
		"char":        Char,
		"void":        Void,
		"var":         Var,
		"static":      Static,
		"field":       Field,
		"let":         Let,
		"do":          Do,
		"if":          If,
		"else":        Else,
		"while":       While,
		"return":      Return,
		"true":        True,
		"false":       False,
		"null":        Null,
		"this":        This,
	}
}

func BinaryOps() mapset.Set[rune] {
	set := mapset.NewSet[rune]()
	set.Add('+')
	set.Add('-')
	set.Add('*')
	set.Add('/')
	set.Add('&')
	set.Add('|')
	set.Add('<')
	set.Add('>')
	set.Add('=')
	return set
}

func BinaryOpsInstructions() map[string]string {
	return map[string]string{
		"+": "add",
		"-": "sub",
		"&": "and",
		"|": "or",
		"<": "lt",
		">": "gt",
		"=": "eq",
	}
}

func UnaryOpsInstructions() map[string]string {
	return map[string]string{
		"-": "neg",
		"~": "not",
	}
}

func UnaryOps() mapset.Set[rune] {
	set := mapset.NewSet[rune]()
	set.Add('~')
	set.Add('-')
	return set
}

func XMLEscapes() map[string]string {
	return map[string]string{
		"<": "&lt;",
		">": "&gt;",
		"&": "&amp;",
	}
}

func SymbolKinds() map[string]SymbolKind {
	return map[string]SymbolKind{
		"static": SymbolStatic,
		"field":  SymbolField,
		"arg":    SymbolArg,
		"var":    SymbolVar,
	}
}

func SymbolKindStrings() map[SymbolKind]string {
	return map[SymbolKind]string{
		SymbolStatic: "static",
		SymbolField:  "this",
		SymbolArg:    "argument",
		SymbolVar:    "local",
	}
}
