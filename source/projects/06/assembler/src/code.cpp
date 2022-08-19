#include <map>

#include "code.h"

Code::Code() {}

string Code::dest(string destStr) {
    string binStr = "000";
    if (destStr.find("A") != -1) {
        binStr[0] = '1';
    }
    if (destStr.find("D") != -1) {
        binStr[1] = '1';
    }
    if (destStr.find("M") != -1) {
        binStr[2] = '1';
    }
    return binStr;
}

string Code::comp(string compStr) {
    return _compBins.at(compStr);
}

string Code::jump(string jumpStr) {
    return _jumpBins.at(jumpStr);
}