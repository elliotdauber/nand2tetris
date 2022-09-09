package main

import (
	"os"
	"strconv"
)

type VMWriter struct {
	outfile                os.File
	class_name             string
	label_counter          int
	class_symboltable      SymbolTable
	subroutine_symboltable SymbolTable
}

func (this *VMWriter) Init(outfile os.File) VMWriter {
	this.outfile = outfile
	this.label_counter = 0
	this.class_symboltable = new(SymbolTable).Init()
	this.subroutine_symboltable = new(SymbolTable).Init()
	return *this
}

func (this *VMWriter) write(text string) {
	this.outfile.WriteString(text)
}

func (this *VMWriter) writeln(text string) {
	this.write(text + "\n")
}

func (this *VMWriter) SetClassName(class_name string) {
	this.class_name = class_name
}

func (this *VMWriter) ClassName() string {
	return this.class_name
}

// comment in this function's body to print comments (should make a cli input)
func (this *VMWriter) WriteComment(text string) {
	// this.writeln("//" + text)
}

func (this *VMWriter) WritePush(segment string, index int) {
	this.writeln("push " + segment + " " + strconv.Itoa(index))
}

func (this *VMWriter) WritePop(segment string, index int) {
	this.writeln("pop " + segment + " " + strconv.Itoa(index))
}

func (this *VMWriter) WriteArithmetic(command string) {
	this.writeln(command)
}

func (this *VMWriter) WriteLabel(label string) {
	this.writeln("label " + label)
}

func (this *VMWriter) WriteGoto(label string) {
	this.writeln("goto " + label)
}

func (this *VMWriter) WriteIf(label string) {
	this.writeln("if-goto " + label)
}

func (this *VMWriter) WriteCall(name string, nArgs int) {
	this.writeln("call " + name + " " + strconv.Itoa(nArgs))
}

func (this *VMWriter) WriteFunction(name string, nVars int) {
	this.writeln("function " + name + " " + strconv.Itoa(nVars))
}

func (this *VMWriter) WriteReturn() {
	this.writeln("return")
}

func (this *VMWriter) ClassSymbolTable() *SymbolTable {
	return &this.class_symboltable
}

func (this *VMWriter) SubroutineSymbolTable() *SymbolTable {
	return &this.subroutine_symboltable
}

func (this *VMWriter) DefineSymbol(name string, typ string, kind SymbolKind) {
	if kind == SymbolField || kind == SymbolStatic {
		this.ClassSymbolTable().Define(name, typ, kind)
	} else {
		this.SubroutineSymbolTable().Define(name, typ, kind)
	}
}

func (this *VMWriter) DefineSymbolWithClassType(name string, kind SymbolKind) {
	this.DefineSymbol(name, this.class_name, kind)
}

func (this *VMWriter) FindSymbol(name string) (SymbolTableEntry, bool) {
	symbol, found := this.SubroutineSymbolTable().FindSymbol(name)
	if found {
		return symbol, found
	}
	symbol, found = this.ClassSymbolTable().FindSymbol(name)
	return symbol, found
}

func (this *VMWriter) GetAndIncLabel() string {
	label := "LABEL" + strconv.Itoa(this.label_counter)
	this.label_counter++
	return label
}
