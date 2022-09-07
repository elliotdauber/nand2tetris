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

func UnaryOps() mapset.Set[rune] {
	set := mapset.NewSet[rune]()
	set.Add('~')
	set.Add('-')
	return set
}

// TODO: MORE??
func XMLEscapes() map[string]string {
	return map[string]string{
		"<": "&lt;",
		">": "&gt;",
		"&": "&amp;",
	}
}
