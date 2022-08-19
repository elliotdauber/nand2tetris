#include "parser.h"

Parser::Parser(ifstream &infile) : _infile(infile) {}

int Parser::hasMoreLines() { 
    return _infile.peek() != EOF;
}

void Parser::advance() {
    reset();

    string line;
    while (line.length() == 0 || isComment(line)) {
        getline(_infile, line);
        line.erase(remove_if(line.begin(), line.end(), ::isspace), line.end());

        //remove trailing comment
        int commentIndex = line.find_first_of("//");
        if (commentIndex != -1) {
            line = line.substr(0, commentIndex);
        }
    }
    
    parseLine(line);
}

void Parser::parseLine(string line) {
    if (line[0] == '@') {
        _instructionType = A_INSTRUCTION;
        parseAInstruction(line);
    } else if (line[0] == '(') {
        _instructionType = L_INSTRUCTION;
        parseLInstruction(line);
    } else  {
        _instructionType = C_INSTRUCTION;
        parseCInstruction(line);
    }
}

void Parser::parseLInstruction(string instruction) {
    //gets xxx from (xxx)
    _symbol = instruction.substr(1, instruction.length() - 2);
};

void Parser::parseCInstruction(string instruction) {
    size_t equalsIndex = instruction.find("=");
    size_t semicolonIndex = instruction.find(";");
    if (equalsIndex != -1) {
        //gets the part of the instruction before the = sign
        _dest = instruction.substr(0, equalsIndex);
        instruction = instruction.substr(equalsIndex + 1);

        if (semicolonIndex != -1) {
            _comp = instruction.substr(0, semicolonIndex);
            _jump = instruction.substr(semicolonIndex + 1);
        } else {
            _comp = instruction;
        }
    } else if (semicolonIndex != -1) {
        _comp = instruction.substr(0, semicolonIndex);
        _jump = instruction.substr(semicolonIndex + 1);
    } else {
        cout << "ERROR: input string has wrong format" << endl; //TODO
    }
};

void Parser::parseAInstruction(string instruction) {
    //gets xxx from @xxx
    _symbol = instruction.substr(1);
};


int Parser::instructionType() { return _instructionType; }
string Parser::symbol() {return _symbol;}
string Parser::dest() {return _dest;}
string Parser::comp() {return _comp;}
string Parser::jump() {return _jump;}


void Parser::reset() {
    _symbol = "";
    _dest = "";
    _comp = "";
    _jump = "";
}

int Parser::isComment(string line) {
    return line.substr(0, 2) == "//";
}