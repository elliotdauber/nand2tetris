package vmtranslator;

import java.io.*;
import java.util.Map;

public class CodeWriter {
    public BufferedWriter _outputWriter;
    private String _staticId;
    private int _labelNum = 0;
    private Map<String, String> _segments = Map.of(
        "local", "LCL",
        "argument", "ARG",
        "this", "THIS",
        "that", "THAT"
    );
    private Map<String, String> _binaryArithCommands = Map.of(
        "add", "M = M + D",
        "sub", "M = M - D",
        "and", "M = M & D",
        "or", "M = M | D"
    );
    private Map<String, String> _unaryArithCommands = Map.of(
        "neg", "D = -D",
        "not", "D = !D",
        "and", "M = M & D",
        "or", "M = M | D"
    );
    private Map<String, String> _equalityArithCommands = Map.of(
        "eq", "JEQ",
        "lt", "JLT",
        "gt", "JGT"
    );

    public CodeWriter(BufferedWriter outputWriter, String staticId) {
        _outputWriter = outputWriter;
        _staticId = staticId;
    }

    public void writeArithmetic(String command) {
        if (_binaryArithCommands.containsKey(command)) {
            write__popFromStack();

            write("@SP");
            write("A = M - 1");
            write(_binaryArithCommands.get(command));
        } else if (_unaryArithCommands.containsKey(command)) {
            write__popFromStack();
            write(_unaryArithCommands.get(command));
            write__pushToStack();
        } else if (_equalityArithCommands.containsKey(command)) {
            write__popFromStack();

            //evaluate
            write("@SP");
            write("A = M - 1");
            write("D = M - D");

            int trueLabel = _labelNum;
            //jump if true
            write("@LABEL" + _labelNum++);
            write("D;" + _equalityArithCommands.get(command));
            
            //load false, then jump
            write__false();

            int falseLabel = _labelNum;
            write("@LABEL" + _labelNum++);
            write("0;JMP");

            write("(LABEL" + trueLabel + ")");
            
            //load "true"
            write__true();

            write("(LABEL" + falseLabel + ")");

            //put result onto stack
            write("@SP");
            write("A = M - 1");
            write("M = D");
        }
    }

    public void writePushPop(int commandType, String segment, int index) {

        if (commandType == CommandType.C_PUSH) {
            write__load(segment, index);
            if (!segment.equals("constant")) {
                write("A = M");
            }
            write("D = A");
            write__pushToStack(); 
        } else if (commandType == CommandType.C_POP) { //TODO: USES R13
            write__load(segment, index); //gets address
            write("D = A"); //saves address to D

            write("@R13"); //save D (this is the address to write to)
            write("M = D");

            write__popFromStack(); //gets item and puts in D

            write("@R13");
            write("A = M");
            write("M = D"); //put D into memory at the stored address

        }
    }

    //puts data in D
    private void write__popFromStack() {
        //decrement SP
        write("@SP");
        write("M = M - 1");

        write("@SP");
        write("A = M");
        write("D = M");
    }

    //assumes data to push is in D
    private void write__pushToStack() {
        write("@SP");
        write("A = M");
        write("M = D");

        //increment SP
        write("@SP");
        write("M = M + 1");
    }

    private void write__load(String segment, int index) {
        if (segment.equals("constant")) {
            write("@" + index);
        } else if (segment.equals("pointer")) {
            if (index == 0) {
                write("@THIS");
            } else if (index == 1) {
                write("@THAT");
            }
        } else if (segment.equals("static")) {
            write("@" + _staticId + "." + index);
        } else if (segment.equals("temp")) {
            write("@" + (5 + index));
        } else if (_segments.containsKey(segment)) {
            write("@" + index);
            write("D = A");

            write("@" + _segments.get(segment));
            write("A = M + D");
        }
    }

    private void write(String line) {
        try {
            _outputWriter.write(line + "\n");
        } catch (IOException e) {
            System.out.println("Could not write to output file");
            e.printStackTrace();
            System.exit(1);
        }
    }

    public void writeInfiniteLoop() {
        write("(INFINITE_LOOP)");
        write("@INFINITE_LOOP");
        write("0;JMP");
    }

    private void write__true() {
        write("@1");
        write("D = -A");
    }

    private void write__false() {
        write("@0");
        write("D = A");
    }

    public void writeComment(String line) {
        write("//" + line);
    }
}
