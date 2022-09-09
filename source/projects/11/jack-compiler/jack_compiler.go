package main

import (
	"os"
)

type JackCompiler struct {
}

func (this *JackCompiler) Init() JackCompiler {
	return *this
}

func (this *JackCompiler) compile_file(filename string) {
	tokenizer := new(JackTokenizer).Init(filename)
	tokenizer.num_tokens()

	outfilename := filename[0:len(filename)-5] + ".vm"
	outfile, err := os.Create(outfilename)
	check(err)
	defer outfile.Close()

	comp_engine := new(CompilationEngine).Init(tokenizer, *outfile)
	class := comp_engine.Compile()
	comp_engine.VMPrint(class)
}
