package main

import (
	"os"
)

type JackAnalyzer struct {
}

func (this *JackAnalyzer) Init() JackAnalyzer {
	return *this
}

func (this *JackAnalyzer) analyze_file(filename string) {
	tokenizer := new(JackTokenizer).Init(filename)
	tokenizer.num_tokens()

	outfilename := filename[0:len(filename)-5] + ".X.xml"
	outfile, err := os.Create(outfilename)
	check(err)
	defer outfile.Close()

	comp_engine := new(CompilationEngine).Init(tokenizer, *outfile)
	class := comp_engine.Compile()
	comp_engine.XMLPrint(class)
}
