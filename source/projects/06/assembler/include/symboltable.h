#include <iostream>
#include <fstream>
#include <string>
#include <map>

using namespace std;

class SymbolTable {
    public:
        SymbolTable(ifstream &infile);
        void addEntry(string symbol, int address);
        void addVariableEntry(string symbol);
        int contains(string symbol);
        int getAddress(string symbol);
    private:
        map<string, int> _symboltable;
        int _nextVariableRegister;
        
        int isComment(string line);
};