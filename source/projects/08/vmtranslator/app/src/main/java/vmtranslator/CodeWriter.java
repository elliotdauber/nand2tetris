package vmtranslator;

import java.io.*;
import java.util.Map;
import java.util.HashMap;

public class CodeWriter {
    public BufferedWriter _outputWriter;
    private int _labelNum = 0;
    private String _filename = "";
    private String _currFunction = "";
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
    private HashMap<String, Integer> _returnCounters = new HashMap<>();

    private int getReturnCounter(String function) {
        int counter = 0;
        if (_returnCounters.containsKey(function)) {
            counter = _returnCounters.get(function);
        }
        _returnCounters.put(function, counter + 1);
        return counter;
    }

    public CodeWriter(BufferedWriter outputWriter) {
        _outputWriter = outputWriter;
    }

    public void setFileName(String filename) {
        _filename = filename;
    }

    public void writeBootstrapCode() {
        //set SP to 256
        write("@256");
        write("D = A");
        write("@SP");
        write("M = D");

        //call Sys.init
        write("@Sys.init");
        write("0;JMP");
    }

    public void writeLabel(String label) {
        write("(" + _currFunction + "$" + label + ")");
    }

    public void writeGoto(String label) {
        write("@" + _currFunction + "$" + label);
        write("0;JMP");
    }

    public void writeIf(String label) {
        write__popFromStack();
        write("@" + _currFunction + "$" + label);
        write("D;JNE");
    }

    public void writeFunction(String functionName, int nVars) {
        _currFunction = functionName;
        write("(" + functionName + ")");
        for (int i = 0; i < nVars; i++) {
            write("@0");
            write("D = A");
            write__pushToStack();
        }
    }

    public void writeCall(String functionName, int nArgs) {
        String returnAddressLabel = _currFunction + "$ret." + getReturnCounter(_currFunction); 
        write("@" + returnAddressLabel);
        write("D = A");
        write__pushToStack();

        write("@LCL");
        write("D = M");
        write__pushToStack();

        write("@ARG");
        write("D = M");
        write__pushToStack();

        write("@THIS");
        write("D = M");
        write__pushToStack();

        write("@THAT");
        write("D = M");
        write__pushToStack();

        //set ARG (args were pushed to the stack just before this frame)
        int argOffset = 5 + nArgs;

        write("@" + argOffset);
        write("D = A");
        write("@SP");
        write("D = M - D");
        write("@ARG");
        write("M = D");

        //LCL = SP
        write("@SP");
        write("D = M");
        write("@LCL");
        write("M = D");

        write("@" + functionName);
        write("0;JMP");

        write("(" + returnAddressLabel + ")");
    }

    public void writeReturn() {
        //save LCL to R13
        write("@LCL");
        write("D = M");
        write("@R13");
        write("M = D");

        //save returnAddress to R14
        write("@5"); //A = 5
        write("D = A"); //D = 5
        write("@R13"); //M = frame
        write("A = M - D"); //A = frame - 5
        write("D = M"); //D = *(frame - 5)
        write("@R14");
        write("M = D"); //R14 = returnAddress

        //put return value in *ARG
        write__popFromStack(); //D = retval
        write("@ARG");
        write("A = M");
        write("M = D"); // *ARG = retval

        //reposition SP
        write("@ARG");
        write("D = M + 1");
        write("@SP");
        write("M = D");

        //restore THAT
        write("@1"); //A = 1
        write("D = A"); //D = 1
        write("@R13"); //M = frame
        write("A = M - D"); //A = frame - 1
        write("D = M"); //D = *(frame - 1)
        write("@THAT");
        write("M = D");

        //restore THIS
        write("@2"); 
        write("D = A");
        write("@R13");
        write("A = M - D");
        write("D = M"); 
        write("@THIS");
        write("M = D");

        //restore ARG
        write("@3"); 
        write("D = A");
        write("@R13");
        write("A = M - D");
        write("D = M"); 
        write("@ARG");
        write("M = D");

        //restore LCL
        write("@4"); 
        write("D = A");
        write("@R13");
        write("A = M - D");
        write("D = M"); 
        write("@LCL");
        write("M = D");

        //get and jump to returnAddress (stored in R14)
        write("@R14");
        write("A = M");
        write("0;JMP");
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
        } else if (commandType == CommandType.C_POP) {
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
            write("@" + _filename + "." + index);
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
