package main

import (
	"os"
)

type XMLWriter struct {
	outfile      os.File
	indent_level int
}

func (this *XMLWriter) Init(outfile os.File) XMLWriter {
	this.outfile = outfile
	this.indent_level = 0
	return *this
}

func (this *XMLWriter) IndentLevelInc() {
	this.indent_level++
}

func (this *XMLWriter) IndentLevelDec() {
	this.indent_level--
}

func (this *XMLWriter) Write(text string) {
	this.outfile.WriteString(text)
}

func (this *XMLWriter) WriteLn(text string) {
	this.Write(text + "\n")
}

func (this *XMLWriter) WriteInnerText(inner_text string) {
	if val, ok := XMLEscapes()[inner_text]; ok {
		inner_text = val
	}
	this.Write(" " + inner_text + " ")
}

func (this *XMLWriter) CreateIndentedTag(tag string) string {
	opening_tag_str := ""
	for i := 0; i < this.indent_level; i++ {
		opening_tag_str += "\t"
	}
	return opening_tag_str + "<" + tag + ">"
}

func (this *XMLWriter) CreateIndentedClosingTag(tag string) string {
	opening_tag_str := ""
	for i := 0; i < this.indent_level; i++ {
		opening_tag_str += "\t"
	}
	return opening_tag_str + "</" + tag + ">"
}

func (this *XMLWriter) WriteOpeningTag(tag string) {
	this.Write(this.CreateIndentedTag(tag))
}

func (this *XMLWriter) WriteOpeningTagNewLn(tag string) {
	this.WriteLn(this.CreateIndentedTag(tag))
}

func (this *XMLWriter) WriteClosingTag(tag string) {
	this.WriteLn("</" + tag + ">")
}

func (this *XMLWriter) WriteClosingTagNewLn(tag string) {
	this.WriteLn(this.CreateIndentedClosingTag(tag))
}

func (this *XMLWriter) WriteTagLn(tag, inner_text string) {
	this.WriteOpeningTag(tag)
	this.WriteInnerText(inner_text)
	this.WriteClosingTag(tag)
}
