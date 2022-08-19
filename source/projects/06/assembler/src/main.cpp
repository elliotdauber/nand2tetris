#include <iostream>
#include <fstream>
#include "assembler.h"

using namespace std;

int main (int argc, char *argv[]) {
    if (argc < 2) {
        cout << "Must supply a filename" << endl;
        return 1;
    }

    Assembler assembler;
    assembler.assemble(argv[1]);

    return 0;
}