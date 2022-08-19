#ifndef PARSER_HEADER
#define PARSER_HEADER

#include <iostream>
#include <string>
#include <fstream>

using namespace std;

enum InstructionType {
    A_INSTRUCTION, C_INSTRUCTION, L_INSTRUCTION
};

class Parser {
    public:
        Parser(ifstream &infile);
        int hasMoreLines();
        void advance();
        int instructionType();
        string symbol();
        string dest();
        string comp();
        string jump();
    private:
        istream &_infile;
        InstructionType _instructionType;
        string _symbol;
        string _dest;
        string _comp;
        string _jump;

        void parseLine(string line);
        void parseAInstruction(string instruction);
        void parseCInstruction(string instruction);
        void parseLInstruction(string instruction);
        void reset();
        int isComment(string line);
};

#endif