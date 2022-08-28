package vmtranslator;

import java.util.Scanner;

public class Parser {
    private Scanner _inputScanner;
    private int _commandType = -1;
    private String _arg1 = "";
    private int _arg2 = -1;
    private String _line = "";

    public Parser(Scanner inputScanner) {
        _inputScanner = inputScanner;
    }

    public boolean hasMoreLines() {
        return _inputScanner.hasNext();
    }

    public void advance() {
        reset();
        String line = "";
        while (line.isBlank() || isComment(line)) {
            line = _inputScanner.nextLine();

            int commentIndex = line.indexOf("//");
            if (commentIndex != -1) {
                line = line.substring(0, commentIndex);
            }

            line = line.strip();
        }
        _line = line;
        parseLine(line);
    }

    public int commandType() {
        return _commandType;
    }

    public String arg1() {
        return _arg1;
    }

    public int arg2() {
        return _arg2;
    }

    public String line() {
        return _line;
    }

    private boolean isComment(String line) {
        return line.startsWith("//");
    }

    private void parseLine(String line) {
        String[] parts = line.split(" ");

        _commandType = CommandType.fromString(parts[0]);
        if (_commandType == CommandType.C_ARITHMETIC) {
            _arg1 = parts[0];
            return;
        }
        if (parts.length > 1) {
            _arg1 = parts[1];
        }
        if (parts.length > 2) {
            _arg2 = Integer.valueOf(parts[2]);
        }
    }

    private void reset() {
        _commandType = -1;
        _arg1 = "";
        _arg2 = -1;
        _line = "";
    }
}
