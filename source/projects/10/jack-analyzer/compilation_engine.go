package main

import (
	"os"
	"strconv"
)

type CompilationEngine struct {
	xml_writer XMLWriter
	tokenizer  JackTokenizer
}

func (this *CompilationEngine) Init(tokenizer JackTokenizer, outfile os.File) CompilationEngine {
	this.tokenizer = tokenizer
	this.xml_writer = new(XMLWriter).Init(outfile)
	return *this
}

func (this *CompilationEngine) Compile() Class {
	this.tokenizer.Advance()
	return this.compile_class()
}

func (this *CompilationEngine) XMLPrint(class Class) {
	XMLPrintNode(&class, this.xml_writer)
}

func (this *CompilationEngine) compile_class() Class {
	name := this.tokenizer.Advance().val
	this.tokenizer.Advance()

	var class_var_decs []ClassVarDec
	var subroutine_decs []SubroutineDec

	for this.tokenizer.HasMoreTokens() {
		next_token := this.tokenizer.Advance()
		if next_token.typ == Symbol && next_token.val == "}" {
			break
		} else if next_token.typ == Keyword {
			keyword := this.tokenizer.Keyword()
			if keyword == Static || keyword == Field {
				class_var_dec := this.compile_class_var_dec()
				class_var_decs = append(class_var_decs, class_var_dec)
			} else {
				subroutine_dec := this.compile_subroutine_dec()
				subroutine_decs = append(subroutine_decs, subroutine_dec)
			}
		}
	}

	return Class{
		name:            name,
		class_var_decs:  class_var_decs,
		subroutine_decs: subroutine_decs,
	}
}

func (this *CompilationEngine) compile_class_var_dec() ClassVarDec {
	dec_type := this.tokenizer.TokenVal()
	var_type := this.tokenizer.Advance().val

	var var_names []string
	for {
		var_name := this.tokenizer.Advance().val
		var_names = append(var_names, var_name)

		symbol := this.tokenizer.Advance().val
		if symbol == ";" {
			break
		}
	}

	return ClassVarDec{
		dec_type:  dec_type,
		var_type:  var_type,
		var_names: var_names,
	}
}

func (this *CompilationEngine) compile_subroutine_dec() SubroutineDec {
	dec_type := this.tokenizer.TokenVal()
	return_type := this.tokenizer.Advance().val
	fn_name := this.tokenizer.Advance().val

	this.tokenizer.Advance()
	this.tokenizer.Advance()
	parameter_list := this.compile_parameter_list()
	this.tokenizer.Advance()

	this.tokenizer.Advance()
	subroutine_body := this.compile_subroutine_body()

	return SubroutineDec{
		dec_type:        dec_type,
		return_type:     return_type,
		name:            fn_name,
		parameter_list:  parameter_list,
		subroutine_body: subroutine_body,
	}
}

func (this *CompilationEngine) compile_parameter_list() ParameterList {
	var parameters []Parameter

	for {
		token_type := this.tokenizer.TokenType()
		token_val := this.tokenizer.TokenVal()
		if token_type == Symbol {
			if token_val == ")" {
				this.tokenizer.Reverse()
				break
			} else {
				this.tokenizer.Advance()
			}
		} else {
			param_type := this.tokenizer.TokenVal()
			param_name := this.tokenizer.Advance().val
			param := Parameter{typ: param_type, name: param_name}
			parameters = append(parameters, param)
			this.tokenizer.Advance()
		}
	}

	return ParameterList{
		parameters: parameters,
	}
}

func (this *CompilationEngine) compile_subroutine_body() SubroutineBody {
	var var_decs []VarDec

	for {
		token_type := this.tokenizer.Advance().typ
		if token_type == Keyword && this.tokenizer.Keyword() == Var {
			var_decs = append(var_decs, this.compile_var_dec())
		} else {
			break
		}
	}

	statements := this.compile_statements()
	this.tokenizer.Advance()

	return SubroutineBody{
		var_decs:   var_decs,
		statements: statements,
	}
}

func (this *CompilationEngine) compile_var_dec() VarDec {
	var_type := this.tokenizer.Advance().val

	var var_names []string

	for {
		var_name := this.tokenizer.Advance().val
		var_names = append(var_names, var_name)

		symbol := this.tokenizer.Advance().val
		if symbol == ";" {
			break
		}
	}

	return VarDec{
		var_type:  var_type,
		var_names: var_names,
	}
}

func (this *CompilationEngine) compile_statements() Statements {
	var statements []Statement

	for {
		token_type := this.tokenizer.TokenType()
		if token_type == Keyword {
			statement_type := this.tokenizer.Keyword()

			if statement_type == Let {
				let_statement := this.compile_let()
				statements = append(statements, &let_statement)
			} else if statement_type == Do {
				do_statement := this.compile_do()
				statements = append(statements, &do_statement)
			} else if statement_type == Return {
				return_statement := this.compile_return()
				statements = append(statements, &return_statement)
			} else if statement_type == While {
				while_statement := this.compile_while()
				statements = append(statements, &while_statement)
			} else if statement_type == If {
				if_statement := this.compile_if()
				statements = append(statements, &if_statement)
			}

			this.tokenizer.Advance()
		} else {
			if this.tokenizer.TokenVal() == "}" {
				this.tokenizer.Reverse()
			}
			break
		}
	}

	return Statements{
		statements: statements,
	}
}

func (this *CompilationEngine) compile_let() LetStatement {
	var_name := this.tokenizer.Advance().val

	var index_expr Expression
	has_index_expr := false
	next_symbol := this.tokenizer.Advance().val
	if next_symbol == "[" {
		this.tokenizer.Advance()
		index_expr = this.compile_expression()
		has_index_expr = true

		this.tokenizer.Advance()
		this.tokenizer.Advance()
	}
	this.tokenizer.Advance()
	assignment_expr := this.compile_expression()
	this.tokenizer.Advance()

	return LetStatement{
		var_name:       var_name,
		index_expr:     index_expr,
		has_index_expr: has_index_expr,
		assigment_expr: assignment_expr,
	}
}

func (this *CompilationEngine) compile_do() DoStatement {
	this.tokenizer.Advance()
	subroutine_call := this.compile_subroutine_call()
	this.tokenizer.Advance()

	return DoStatement{
		subroutine_call: subroutine_call,
	}
}

func (this *CompilationEngine) compile_subroutine_call() SubroutineCall {
	object_or_subroutine := this.tokenizer.TokenVal()
	symbol := this.tokenizer.Advance()

	var subroutine_name string
	var object string
	has_object := false
	if symbol.val == "." {
		has_object = true
		object = object_or_subroutine
		subroutine_name = this.tokenizer.Advance().val
		this.tokenizer.Advance()
	} else {
		subroutine_name = object_or_subroutine
	}

	next_token := this.tokenizer.Advance() //first of expr, but might be ')'
	var expr_list ExpressionList
	if next_token.val != ")" {
		expr_list = this.compile_expression_list()
		this.tokenizer.Advance()
	}

	return SubroutineCall{
		object:          object,
		has_object:      has_object,
		subroutine_name: subroutine_name,
		expr_list:       expr_list,
	}
}

func (this *CompilationEngine) compile_return() ReturnStatement {
	next_token := this.tokenizer.Advance().val
	var expr Expression
	has_expr := false
	if next_token != ";" {
		expr = this.compile_expression()
		has_expr = true
		this.tokenizer.Advance()
	}
	return ReturnStatement{
		expr:     expr,
		has_expr: has_expr,
	}
}

func (this *CompilationEngine) compile_while() WhileStatement {
	this.tokenizer.Advance()

	this.tokenizer.Advance()
	check_expr := this.compile_expression()
	this.tokenizer.Advance()

	this.tokenizer.Advance()
	this.tokenizer.Advance()

	statements := this.compile_statements()
	this.tokenizer.Advance()

	return WhileStatement{
		check_expr: check_expr,
		statements: statements,
	}
}

func (this *CompilationEngine) compile_if() IfStatement {
	this.tokenizer.Advance()

	this.tokenizer.Advance()
	cond_expr := this.compile_expression()
	this.tokenizer.Advance()

	this.tokenizer.Advance()
	this.tokenizer.Advance()

	if_statements := this.compile_statements()
	this.tokenizer.Advance()

	var else_statements Statements
	has_else_statements := false
	this.tokenizer.Advance()
	if this.tokenizer.TokenType() == Keyword && this.tokenizer.Keyword() == Else {
		has_else_statements = true
		this.tokenizer.Advance()

		this.tokenizer.Advance()
		else_statements = this.compile_statements()
		this.tokenizer.Advance()
	} else {
		this.tokenizer.Reverse()
	}
	return IfStatement{
		cond_expr:           cond_expr,
		if_statements:       if_statements,
		else_statements:     else_statements,
		has_else_statements: has_else_statements,
	}
}

func (this *CompilationEngine) compile_expression() Expression {
	var terms []Term
	var ops []string

	for {
		term := this.compile_term()
		terms = append(terms, term)
		next_token := this.tokenizer.Advance().val
		if BinaryOps().Contains(rune(next_token[0])) {
			ops = append(ops, next_token)
			this.tokenizer.Advance()
		} else {
			this.tokenizer.Reverse()
			break
		}
	}
	return Expression{
		terms: terms,
		ops:   ops,
	}
}

func (this *CompilationEngine) compile_expression_list() ExpressionList {
	var exprs []Expression

	for {
		expr := this.compile_expression()
		exprs = append(exprs, expr)
		next_token := this.tokenizer.Advance().val
		if next_token == "," {
			this.tokenizer.Advance()
		} else {
			this.tokenizer.Reverse()
			break
		}
	}

	return ExpressionList{
		exprs: exprs,
	}
}

func (this *CompilationEngine) compile_term() Term {
	token_type := this.tokenizer.TokenType()
	token_val := this.tokenizer.TokenVal()
	if token_type == IntConst {
		int_val, _ := strconv.Atoi(token_val)
		return &IntegerConstant{val: int_val}
	} else if token_type == StrConst {
		return &StringConstant{val: token_val}
	} else if token_type == Keyword {
		return &KeywordConst{val: token_val}
	} else if UnaryOps().Contains(rune(token_val[0])) {
		this.tokenizer.Advance()
		term := this.compile_term()
		return &UnaryOpTerm{op: token_val, term: term}
	} else if token_val == "(" {
		this.tokenizer.Advance()
		expr := this.compile_expression()
		this.tokenizer.Advance()
		return &ParenTerm{expr: expr}
	} else {
		next_token := this.tokenizer.Advance().val
		if next_token == "(" || next_token == "." {
			this.tokenizer.Reverse()
			subroutine_call := this.compile_subroutine_call()
			return &subroutine_call
		} else {
			var index_expr Expression
			has_index_expr := false
			if next_token == "[" {
				has_index_expr = true
				this.tokenizer.Advance()
				index_expr = this.compile_expression()
				this.tokenizer.Advance()
			} else {
				this.tokenizer.Reverse()
			}
			return &VarTerm{name: token_val, has_index_expr: has_index_expr, index_expr: index_expr}
		}
	}
}
