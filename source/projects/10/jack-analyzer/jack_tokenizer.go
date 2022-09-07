package main

import (
	"os"
	"regexp"
	"unicode"
)

type JackTokenizer struct {
	tokens    []Token
	token_idx int
}

func (this *JackTokenizer) Init(filename string) JackTokenizer {
	this.parse_input(filename)
	this.token_idx = -1
	return *this
}

/* PUBLIC FNS */

func (this *JackTokenizer) HasMoreTokens() bool {
	return this.token_idx < len(this.tokens)-1
}

func (this *JackTokenizer) Advance() Token {
	this.token_idx++
	return this.tokens[this.token_idx]
}

func (this *JackTokenizer) Reverse() Token {
	this.token_idx--
	return this.tokens[this.token_idx]
}

func (this *JackTokenizer) TokenVal() string {
	return this.tokens[this.token_idx].val
}

func (this *JackTokenizer) TokenType() TokenType {
	return this.tokens[this.token_idx].typ
}

func (this *JackTokenizer) Keyword() KeyWord {
	input := this.tokens[this.token_idx].val
	return KeyWordMap()[input]
}

/* PRIVATE FNS */

func (this *JackTokenizer) parse_input(filename string) {
	file_text, err := os.ReadFile(filename)
	check(err)
	input_string := string(file_text)
	input_string = this.remove_comments(input_string)
	this.lex(input_string)
}

func (this *JackTokenizer) remove_comments(s string) string {
	single_line_comments_regex := regexp.MustCompile(`//.*(\n|$)`)
	ret_s := single_line_comments_regex.ReplaceAllString(s, "")

	multi_line_comments_regex := regexp.MustCompile(`/\*`) // /\*(.*)\*/
	for {
		find_result := multi_line_comments_regex.FindStringIndex(ret_s)
		if find_result == nil {
			break
		}
		comment_idx := find_result[0]
		curr_idx := comment_idx + 2
		for {
			next_slice := ret_s[curr_idx : curr_idx+2]
			if next_slice == "*/" {
				ret_s = ret_s[0:comment_idx] + ret_s[curr_idx+3:]
				break
			}
			curr_idx++
		}
	}

	return ret_s
}

func (_ JackTokenizer) parse_str_const(input string) Token {
	re := regexp.MustCompile(`"[^"]*"`)
	res := re.FindStringIndex(input)
	return new(Token).Init(input[res[0]+1:res[1]-1], StrConst)
}

func (_ JackTokenizer) parse_int_const(input string) Token {
	re := regexp.MustCompile(`[0-9]+`)
	res := re.FindStringIndex(input)
	return new(Token).Init(input[res[0]:res[1]], IntConst)
}

func (_ JackTokenizer) parse_keyword_or_identifier(input string) Token {
	re := regexp.MustCompile(`class|constructor|function|method|field|static|var|int|char|boolean|void|true|false|null|this|let|do|if|else|while|return|[a-zA-Z_]+[a-zA-Z0-9_]*`)
	res := re.FindStringIndex(input)
	val := input[res[0]:res[1]]
	typ := Identifier
	if _, ok := KeyWordMap()[val]; ok {
		typ = Keyword
	}
	return new(Token).Init(val, typ)
}

func (this *JackTokenizer) lex(input string) {
	index := 0

	for index < len(input) {
		curr_char := rune(input[index])
		var next_token Token
		if unicode.IsSpace(curr_char) {
			index++
			continue
		} else if string(curr_char) == "\"" {
			next_token = this.parse_str_const(input[index:])
			index += 2 //accounts for the quote marks, which are not included in the token
		} else if unicode.IsNumber(curr_char) {
			next_token = this.parse_int_const(input[index:])
		} else if SymbolSet().Contains(curr_char) {
			next_token = new(Token).Init(string(curr_char), Symbol)
		} else {
			next_token = this.parse_keyword_or_identifier(input[index:])
		}

		index += len(next_token.val)
		this.tokens = append(this.tokens, next_token)
	}

	if this.num_tokens() > 0 {
		this.token_idx = 0
	}
}

/** simple funcs :) */

func (this *JackTokenizer) num_tokens() int {
	return len(this.tokens)
}
