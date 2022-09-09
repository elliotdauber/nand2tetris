package main

/** NODE **/

type Node interface {
	VMPrint(writer *VMWriter)
}

/** CLASS **/

type Class struct {
	name            string
	class_var_decs  []ClassVarDec
	subroutine_decs []SubroutineDec
	Node
}

func (this *Class) VMPrint(writer *VMWriter) {
	writer.WriteComment("class")
	writer.SetClassName(this.name)
	writer.ClassSymbolTable().Reset()
	for _, class_var_dec := range this.class_var_decs {
		class_var_dec.VMPrint(writer)
	}

	for _, subroutine_dec := range this.subroutine_decs {
		subroutine_dec.VMPrint(writer)
	}
}

/** ClassVarDec **/

type ClassVarDec struct {
	dec_type  string
	var_type  string
	var_names []string
	Node
}

func (this *ClassVarDec) VMPrint(writer *VMWriter) {
	writer.WriteComment("classVarDec")
	for _, var_name := range this.var_names {
		writer.DefineSymbol(var_name, this.var_type, get_symbol_kind(this.dec_type))
	}
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

func (this *SubroutineDec) VMPrint(writer *VMWriter) {
	writer.WriteComment("subroutineDec")
	writer.SubroutineSymbolTable().Reset()

	num_vars := this.subroutine_body.num_vars()
	writer.WriteFunction(writer.ClassName()+"."+this.name, num_vars)

	if this.dec_type == "method" {
		writer.DefineSymbolWithClassType("this", SymbolArg)
		writer.WritePush("argument", 0)
		writer.WritePop("pointer", 0)
	} else if this.dec_type == "constructor" {
		class_size := writer.ClassSymbolTable().VarCount(SymbolField)
		writer.WritePush("constant", class_size)
		writer.WriteCall("Memory.alloc", 1)
		writer.WritePop("pointer", 0)
	}

	this.parameter_list.VMPrint(writer)
	this.subroutine_body.VMPrint(writer)
}

/** ParameterList **/

type ParameterList struct {
	parameters []Parameter
	Node
}

type Parameter struct {
	typ  string
	name string
}

func (this *ParameterList) VMPrint(writer *VMWriter) {
	writer.WriteComment("parameterList")
	for _, parameter := range this.parameters {
		writer.WriteComment("parameter")
		writer.DefineSymbol(parameter.name, parameter.typ, SymbolArg)
	}
}

/** SubroutineBody **/

type SubroutineBody struct {
	var_decs   []VarDec
	statements Statements
	Node
}

func (this *SubroutineBody) VMPrint(writer *VMWriter) {
	writer.WriteComment("subroutineBody")
	for _, var_dec := range this.var_decs {
		var_dec.VMPrint(writer)
	}

	this.statements.VMPrint(writer)
}

func (this *SubroutineBody) num_vars() int {
	num_vars := 0
	for _, var_dec := range this.var_decs {
		num_vars += len(var_dec.var_names)
	}
	return num_vars
}

/** VarDec **/

type VarDec struct {
	var_type  string
	var_names []string
	Node
}

func (this *VarDec) VMPrint(writer *VMWriter) {
	writer.WriteComment("varDec")
	for _, var_name := range this.var_names {
		writer.DefineSymbol(var_name, this.var_type, SymbolVar)
	}
}

/** Statements **/

type Statement interface {
	Node
}

type Statements struct {
	statements []Statement
	Node
}

func (this *Statements) VMPrint(writer *VMWriter) {
	writer.WriteComment("statments")
	for _, statement := range this.statements {
		statement.VMPrint(writer)
	}
}

/** LetStatement **/

type LetStatement struct {
	var_name       string
	index_expr     Expression
	has_index_expr bool
	assigment_expr Expression
	Statement
}

func (this *LetStatement) VMPrint(writer *VMWriter) {
	writer.WriteComment("letStatement")
	symbol, _ := writer.FindSymbol(this.var_name)

	if this.has_index_expr {
		writer.WritePush(get_symbol_kind_string(symbol.kind), symbol.index)
		this.index_expr.VMPrint(writer)
		writer.WriteArithmetic("add")

		this.assigment_expr.VMPrint(writer)
		writer.WritePop("temp", 0)

		writer.WritePop("pointer", 1)
		writer.WritePush("temp", 0)
		writer.WritePop("that", 0)
	} else {
		this.assigment_expr.VMPrint(writer)
		writer.WritePop(get_symbol_kind_string(symbol.kind), symbol.index)
	}
}

/** DoStatement **/

type DoStatement struct {
	subroutine_call SubroutineCall
	Statement
}

func (this *DoStatement) VMPrint(writer *VMWriter) {
	writer.WriteComment("doStatement")
	this.subroutine_call.VMPrint(writer)
	writer.WritePop("temp", 0)
}

/** ReturnStatement **/

type ReturnStatement struct {
	expr     Expression
	has_expr bool
	Statement
}

// todo: push void if no expr??
func (this *ReturnStatement) VMPrint(writer *VMWriter) {
	writer.WriteComment("returnStatement")
	if this.has_expr {
		this.expr.VMPrint(writer)
	} else {
		writer.WritePush("constant", 0)
	}
	writer.WriteReturn()
}

/** WhileStatement **/

type WhileStatement struct {
	check_expr Expression
	statements Statements
	Statement
}

func (this *WhileStatement) VMPrint(writer *VMWriter) {
	writer.WriteComment("whileStatement")
	check_label := writer.GetAndIncLabel()
	writer.WriteLabel(check_label)

	this.check_expr.VMPrint(writer)
	writer.WriteArithmetic("not")

	loopend_label := writer.GetAndIncLabel()
	writer.WriteIf(loopend_label)

	this.statements.VMPrint(writer)
	writer.WriteGoto(check_label)

	writer.WriteLabel(loopend_label)
}

/** IfStatement **/

type IfStatement struct {
	cond_expr           Expression
	if_statements       Statements
	else_statements     Statements
	has_else_statements bool
	Statement
}

func (this *IfStatement) VMPrint(writer *VMWriter) {
	writer.WriteComment("ifStatement")
	this.cond_expr.VMPrint(writer)
	writer.WriteArithmetic("not")

	false_label := writer.GetAndIncLabel()
	writer.WriteIf(false_label)

	this.if_statements.VMPrint(writer)
	if this.has_else_statements {
		endif_label := writer.GetAndIncLabel()
		writer.WriteGoto(endif_label)
		writer.WriteLabel(false_label)

		this.else_statements.VMPrint(writer)
		writer.WriteLabel(endif_label)
	} else {
		writer.WriteLabel(false_label)
	}
}

/** Expressions **/

type Expression struct {
	terms []Term
	ops   []string
	Node
}

// todo: how do we handle & and |
func (this *Expression) VMPrint(writer *VMWriter) {
	writer.WriteComment("expression")
	for i, term := range this.terms {
		term.VMPrint(writer)
		if i > 0 {
			op := this.ops[i-1]
			if op == "*" {
				writer.WriteCall("Math.multiply", 2)
			} else if op == "/" {
				writer.WriteCall("Math.divide", 2)
			} else {
				op_instruction := get_binop_instruction(op)
				writer.WriteArithmetic(op_instruction)
			}
		}
	}
}

/** Expressionlist **/

type ExpressionList struct {
	exprs []Expression
	Node
}

func (this *ExpressionList) VMPrint(writer *VMWriter) {
	writer.WriteComment("expressionList")
	for _, expr := range this.exprs {
		expr.VMPrint(writer)
	}
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

func (this *IntegerConstant) VMPrint(writer *VMWriter) {
	writer.WriteComment("integerConstant")
	writer.WritePush("constant", this.val)
}

/** StringConstant **/

type StringConstant struct {
	val string
	Term
}

func (this *StringConstant) VMPrint(writer *VMWriter) {
	writer.WriteComment("stringConstant")
	strlen := len(this.val)
	writer.WritePush("constant", strlen)

	writer.WriteCall("String.new", 1)

	for _, char := range this.val {
		writer.WritePush("constant", int(char))
		writer.WriteCall("String.appendChar", 2)
	}
}

/** KeywordConst **/

type KeywordConst struct {
	val string
	Term
}

func (this *KeywordConst) VMPrint(writer *VMWriter) {
	writer.WriteComment("keywordConst")
	if this.val == "true" {
		writer.WritePush("constant", 1)
		writer.WriteArithmetic("neg")
	} else if this.val == "false" || this.val == "null" {
		writer.WritePush("constant", 0)
	} else if this.val == "this" {
		_this, found := writer.FindSymbol("this")
		if found {
			writer.WritePush(get_symbol_kind_string(_this.kind), _this.index)
		} else {
			writer.WritePush("pointer", 0)
		}
	}
}

/** VarTerm **/

type VarTerm struct {
	name           string
	index_expr     Expression
	has_index_expr bool
	Term
}

func (this *VarTerm) VMPrint(writer *VMWriter) {
	writer.WriteComment("varTerm")
	symbol, _ := writer.FindSymbol(this.name)

	if this.has_index_expr {
		writer.WritePush(get_symbol_kind_string(symbol.kind), symbol.index)
		this.index_expr.VMPrint(writer)
		writer.WriteArithmetic("add")
		writer.WritePop("pointer", 1)

		writer.WritePush("that", 0)
	} else {
		writer.WritePush(get_symbol_kind_string(symbol.kind), symbol.index)
	}
}

/** ParenTerm **/

type ParenTerm struct {
	expr Expression
	Term
}

func (this *ParenTerm) VMPrint(writer *VMWriter) {
	writer.WriteComment("parenTerm")
	this.expr.VMPrint(writer)
}

/** UnaryOpTerm **/

type UnaryOpTerm struct {
	op   string
	term Term
	Term
}

func (this *UnaryOpTerm) VMPrint(writer *VMWriter) {
	writer.WriteComment("unaryOpTerm")
	this.term.VMPrint(writer)
	writer.WriteArithmetic(get_unaryop_instruction(this.op))
}

/** SubroutineCall **/

type SubroutineCall struct {
	object          string
	has_object      bool
	subroutine_name string
	expr_list       ExpressionList
	Term
}

func (this *SubroutineCall) VMPrint(writer *VMWriter) {
	writer.WriteComment("subroutineCall")
	var callee_type string
	num_exprs := len(this.expr_list.exprs)
	if this.has_object {
		symbol, found := writer.FindSymbol(this.object)
		if found {
			//this is a method call
			num_exprs += 1
			callee_type = symbol.typ
			writer.WritePush(get_symbol_kind_string(symbol.kind), symbol.index)
		} else {
			//this is a static or constructor call
			callee_type = this.object
		}
	} else {
		//this is a method call on the current object
		num_exprs += 1
		callee_type = writer.ClassName()
		_this, found := writer.FindSymbol("this")
		if found {
			writer.WritePush(get_symbol_kind_string(_this.kind), _this.index)
		} else {
			writer.WritePush("pointer", 0)
		}
	}

	this.expr_list.VMPrint(writer)

	writer.WriteCall(callee_type+"."+this.subroutine_name, num_exprs)
}

/** Helper functions **/

func get_symbol_kind(kind_as_string string) SymbolKind {
	if val, found := SymbolKinds()[kind_as_string]; found {
		return val
	}
	return SymbolNone
}

func get_symbol_kind_string(kind SymbolKind) string {
	if val, found := SymbolKindStrings()[kind]; found {
		return val
	}
	return ""
}

func get_binop_instruction(op string) string {
	if val, found := BinaryOpsInstructions()[op]; found {
		return val
	}
	return ""
}

func get_unaryop_instruction(op string) string {
	if val, found := UnaryOpsInstructions()[op]; found {
		return val
	}
	return ""
}
