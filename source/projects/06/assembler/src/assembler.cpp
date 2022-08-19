#include <bitset>
#include "assembler.h"
#include "parser.h"
#include "code.h"
#include "symboltable.h"

Assembler::Assembler() {}

void Assembler::assemble(string filename) {
    ifstream infilePreparse(filename);
    SymbolTable symboltable(infilePreparse);
    infilePreparse.close();

    ifstream infile(filename);

    filename.replace(filename.length() - 4, 4, ".hack");
    ofstream outfile(filename);

    Parser parser(infile);
    Code code;
    while (parser.hasMoreLines()) {
        parser.advance();
        if (parser.instructionType() == A_INSTRUCTION) {
            string symbol = parser.symbol();
            if (!symboltable.contains(symbol)) {
                symboltable.addVariableEntry(symbol);
            }
            int address = symboltable.getAddress(symbol);
            string addressBin = "0" + bitset<15>(address).to_string();
            outfile << addressBin << endl;
        } else if (parser.instructionType() == C_INSTRUCTION) {
            outfile << "111" << code.comp(parser.comp()) << code.dest(parser.dest()) << code.jump(parser.jump()) << endl;
        }
    }
    infile.close();
    outfile.close();
}