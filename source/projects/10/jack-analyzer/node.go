package main

import "fmt"

/** NODE **/

type Node interface {
	XMLPrint(writer XMLWriter)
	XMLTag() string
}

func XMLPrintNode(node Node, writer XMLWriter) {
	writer.WriteOpeningTagNewLn(node.XMLTag())
	writer.IndentLevelInc()
	node.XMLPrint(writer)
	writer.IndentLevelDec()
	writer.WriteClosingTagNewLn(node.XMLTag())
}

/** CLASS **/

type Class struct {
	name            string
	class_var_decs  []ClassVarDec
	subroutine_decs []SubroutineDec
}

func (this *Class) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", "class")
	writer.WriteTagLn("identifier", this.name)
	writer.WriteTagLn("symbol", "{")

	for _, class_var_dec := range this.class_var_decs {
		XMLPrintNode(&class_var_dec, writer)
	}

	for _, subroutine_dec := range this.subroutine_decs {
		XMLPrintNode(&subroutine_dec, writer)
	}

	writer.WriteTagLn("symbol", "}")
}

func (this *Class) XMLTag() string {
	return "class"
}

/** ClassVarDec **/

type ClassVarDec struct {
	dec_type  string
	var_type  string
	var_names []string
	Node
}

func (this *ClassVarDec) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", this.dec_type)
	var_type_tag := "identifier"
	if _, ok := KeyWordMap()[this.var_type]; ok {
		var_type_tag = "keyword"
	}
	writer.WriteTagLn(var_type_tag, this.var_type)

	for i, var_name := range this.var_names {
		if i > 0 {
			writer.WriteTagLn("symbol", ",")
		}
		writer.WriteTagLn("identifier", var_name)
	}
	writer.WriteTagLn("symbol", ";")
}

func (this *ClassVarDec) XMLTag() string {
	return "classVarDec"
}

/** SubroutineDec **/

type SubroutineDec struct {
	dec_type        string
	return_type     string
	name            string
	parameter_list  ParameterList
	subroutine_body SubroutineBody
	Node
}

func (this *SubroutineDec) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", this.dec_type)
	return_type_tag := "identifier"
	if _, ok := KeyWordMap()[this.return_type]; ok {
		return_type_tag = "keyword"
	}
	writer.WriteTagLn(return_type_tag, this.return_type)

	writer.WriteTagLn("identifier", this.name)
	writer.WriteTagLn("symbol", "(")
	XMLPrintNode(&this.parameter_list, writer)

	writer.WriteTagLn("symbol", ")")
	XMLPrintNode(&this.subroutine_body, writer)
}

func (this *SubroutineDec) XMLTag() string {
	return "subroutineDec"
}

/** ParameterList **/

type ParameterList struct {
	parameters []Parameter
	Node
}

func (this *ParameterList) XMLPrint(writer XMLWriter) {
	for i, param := range this.parameters {
		if i > 0 {
			writer.WriteTagLn("symbol", ",")
		}

		param_type_tag := "identifier"
		if _, ok := KeyWordMap()[param.typ]; ok {
			param_type_tag = "keyword"
		}
		writer.WriteTagLn(param_type_tag, param.typ)
		writer.WriteTagLn("identifier", param.name)
	}
}

func (this *ParameterList) XMLTag() string {
	return "parameterList"
}

type Parameter struct {
	typ  string
	name string
}

/** SubroutineBody **/

type SubroutineBody struct {
	var_decs   []VarDec
	statements Statements
	Node
}

func (this *SubroutineBody) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("symbol", "{")
	for _, var_dec := range this.var_decs {
		XMLPrintNode(&var_dec, writer)
	}
	XMLPrintNode(&this.statements, writer)
	writer.WriteTagLn("symbol", "}")
}

func (this *SubroutineBody) XMLTag() string {
	return "subroutineBody"
}

/** VarDec **/

type VarDec struct {
	var_type  string
	var_names []string
	Node
}

func (this *VarDec) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", "var")

	var_type_tag := "identifier"
	if _, ok := KeyWordMap()[this.var_type]; ok {
		var_type_tag = "keyword"
	}
	writer.WriteTagLn(var_type_tag, this.var_type)

	for i, var_name := range this.var_names {
		if i > 0 {
			writer.WriteTagLn("symbol", ",")
		}
		writer.WriteTagLn("identifier", var_name)
	}
	writer.WriteTagLn("symbol", ";")
}

func (this *VarDec) XMLTag() string {
	return "varDec"
}

/** Statements **/

type Statement interface {
	Node
}

type Statements struct {
	statements []Statement
	Node
}

func (this *Statements) XMLPrint(writer XMLWriter) {
	for _, statement := range this.statements {
		XMLPrintNode(statement, writer)
	}
}

func (this *Statements) XMLTag() string {
	return "statements"
}

/** LetStatement **/

type LetStatement struct {
	var_name       string
	index_expr     Expression
	has_index_expr bool
	assigment_expr Expression
	Statement
}

func (this *LetStatement) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", "let")
	writer.WriteTagLn("identifier", this.var_name)

	if this.has_index_expr {
		writer.WriteTagLn("symbol", "[")
		XMLPrintNode(&this.index_expr, writer)
		writer.WriteTagLn("symbol", "]")
	}

	writer.WriteTagLn("symbol", "=")
	XMLPrintNode(&this.assigment_expr, writer)
	writer.WriteTagLn("symbol", ";")
}

func (this *LetStatement) XMLTag() string {
	return "letStatement"
}

/** DoStatement **/

type DoStatement struct {
	subroutine_call SubroutineCall
	Statement
}

func (this *DoStatement) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", "do")
	this.subroutine_call.XMLPrint(writer)
	writer.WriteTagLn("symbol", ";")
}

func (this *DoStatement) XMLTag() string {
	return "doStatement"
}

/** ReturnStatement **/

type ReturnStatement struct {
	expr     Expression
	has_expr bool
	Statement
}

func (this *ReturnStatement) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", "return")

	if this.has_expr {
		XMLPrintNode(&this.expr, writer)
	}

	writer.WriteTagLn("symbol", ";")
}

func (this *ReturnStatement) XMLTag() string {
	return "returnStatement"
}

/** WhileStatement **/

type WhileStatement struct {
	check_expr Expression
	statements Statements
	Statement
}

func (this *WhileStatement) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", "while")
	writer.WriteTagLn("symbol", "(")
	XMLPrintNode(&this.check_expr, writer)
	writer.WriteTagLn("symbol", ")")
	writer.WriteTagLn("symbol", "{")
	XMLPrintNode(&this.statements, writer)
	writer.WriteTagLn("symbol", "}")
}

func (this *WhileStatement) XMLTag() string {
	return "whileStatement"
}

/** IfStatement **/

type IfStatement struct {
	cond_expr           Expression
	if_statements       Statements
	else_statements     Statements
	has_else_statements bool
	Statement
}

func (this *IfStatement) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", "if")
	writer.WriteTagLn("symbol", "(")
	XMLPrintNode(&this.cond_expr, writer)
	writer.WriteTagLn("symbol", ")")
	writer.WriteTagLn("symbol", "{")
	XMLPrintNode(&this.if_statements, writer)
	writer.WriteTagLn("symbol", "}")

	if this.has_else_statements {
		writer.WriteTagLn("keyword", "else")
		writer.WriteTagLn("symbol", "{")
		XMLPrintNode(&this.else_statements, writer)
		writer.WriteTagLn("symbol", "}")
	}
}

func (this *IfStatement) XMLTag() string {
	return "ifStatement"
}

/** Expressions **/

type Expression struct {
	terms []Term
	ops   []string
	Node
}

func (this *Expression) XMLPrint(writer XMLWriter) {
	for i, term := range this.terms {
		if i > 0 {
			op := this.ops[i-1]
			writer.WriteTagLn("symbol", op)
		}
		XMLPrintNode(term, writer)
	}
}

func (this *Expression) XMLTag() string {
	return "expression"
}

/** Expressionlist **/

type ExpressionList struct {
	exprs []Expression
	Node
}

func (this *ExpressionList) XMLPrint(writer XMLWriter) {
	for i, expr := range this.exprs {
		if i > 0 {
			writer.WriteTagLn("symbol", ",")
		}
		XMLPrintNode(&expr, writer)
	}
}

func (this *ExpressionList) XMLTag() string {
	return "expressionList"
}

/** Terms **/

type Term interface {
	Node
}

/** IntegerConstant **/

type IntegerConstant struct {
	val int
	Term
}

func (this *IntegerConstant) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("integerConstant", fmt.Sprint(this.val))
}

func (this *IntegerConstant) XMLTag() string {
	return "term"
}

/** StringConstant **/

type StringConstant struct {
	val string
	Term
}

func (this *StringConstant) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("stringConstant", this.val)
}

func (this *StringConstant) XMLTag() string {
	return "term"
}

/** KeywordConst **/

type KeywordConst struct {
	val string
	Term
}

func (this *KeywordConst) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("keyword", this.val)
}

func (this *KeywordConst) XMLTag() string {
	return "term"
}

/** VarTerm **/

type VarTerm struct {
	name           string
	index_expr     Expression
	has_index_expr bool
	Term
}

func (this *VarTerm) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("identifier", this.name)
	if this.has_index_expr {
		writer.WriteTagLn("symbol", "[")
		XMLPrintNode(&this.index_expr, writer)
		writer.WriteTagLn("symbol", "]")
	}
}

func (this *VarTerm) XMLTag() string {
	return "term"
}

/** ParenTerm **/

type ParenTerm struct {
	expr Expression
	Term
}

func (this *ParenTerm) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("symbol", "(")
	XMLPrintNode(&this.expr, writer)
	writer.WriteTagLn("symbol", ")")
}

func (this *ParenTerm) XMLTag() string {
	return "term"
}

/** ParenTerm **/

type UnaryOpTerm struct {
	op   string
	term Term
	Term
}

func (this *UnaryOpTerm) XMLPrint(writer XMLWriter) {
	writer.WriteTagLn("symbol", this.op)
	XMLPrintNode(this.term, writer)
}

func (this *UnaryOpTerm) XMLTag() string {
	return "term"
}

/** SubroutineCall **/

type SubroutineCall struct {
	object          string
	has_object      bool
	subroutine_name string
	expr_list       ExpressionList
	Term
}

func (this *SubroutineCall) XMLPrint(writer XMLWriter) {
	if this.has_object {
		var_tag := "identifier"
		if _, ok := KeyWordMap()[this.object]; ok {
			var_tag = "keyword"
		}
		writer.WriteTagLn(var_tag, this.object)
		writer.WriteTagLn("symbol", ".")
	}
	writer.WriteTagLn("identifier", this.subroutine_name)
	writer.WriteTagLn("symbol", "(")
	XMLPrintNode(&this.expr_list, writer)
	writer.WriteTagLn("symbol", ")")
}

func (this *SubroutineCall) XMLTag() string {
	return "term"
}
