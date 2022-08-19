#include "symboltable.h"

SymbolTable::SymbolTable(ifstream &infile) : _nextVariableRegister(16) {
    _symboltable = {
        {"R0", 0},
        {"R1", 1},
        {"R2", 2},
        {"R3", 3},
        {"R4", 4},
        {"R5", 5},
        {"R6", 6},
        {"R7", 7},
        {"R8", 8},
        {"R9", 9},
        {"R10", 10},
        {"R11", 11},
        {"R12", 12},
        {"R13", 13},
        {"R14", 14},
        {"R15", 15},

        {"SP", 0},
        {"LCL", 1},
        {"ARG", 2},
        {"THIS", 3},
        {"THAT", 4},

        {"SCREEN", 16384},
        {"KBD", 24576}
    };

    int lineno = 0;
    string line;
    while (getline(infile, line)) {
        //remove whitespace
        line.erase(remove_if(line.begin(), line.end(), ::isspace), line.end());

        //skip blank lines and comments
        if (line.length() == 0 || isComment(line)) {
            continue;
        } 

        //add labels to the symboltable
        if (line[0] == '(') {
            string label = line.substr(1, line.length() - 2);
            addEntry(label, lineno);
            continue;   
        }
        lineno ++;
    }
}

void SymbolTable::addEntry(string symbol, int address) {
    pair<string, int> newEntry = {symbol, address};
    _symboltable.insert(newEntry);
}

void SymbolTable::addVariableEntry(string symbol) {
    if (isdigit(symbol[0])) {
        return;
    }
    pair<string, int> newEntry = {symbol, _nextVariableRegister++};
    _symboltable.insert(newEntry);
}

int SymbolTable::contains(string symbol) {
    return _symboltable.count(symbol);
}

//means an address was already passed in
int SymbolTable::getAddress(string symbol) {
    if (isdigit(symbol[0])) {
        return atoi(symbol.c_str());
    }
    return _symboltable.at(symbol);
}

//TODO: put in utils?
int SymbolTable::isComment(string line) {
    return line.substr(0, 2) == "//";
}