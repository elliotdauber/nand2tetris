package vmtranslator;

import java.io.*;  
import java.util.Scanner;

public class VMTranslator {
    public void translate(String filename) throws IOException {
        Scanner inputScanner = createInputScanner(filename);
        Parser parser = new Parser(inputScanner);

        BufferedWriter fileWriter = createFileWriter(filename);

        String staticId = filename.substring(filename.lastIndexOf("/") + 1, filename.length() - ".vm".length());
        CodeWriter codeWriter = new CodeWriter(fileWriter, staticId);

        while (parser.hasMoreLines()) {
            parser.advance();
            codeWriter.writeComment(parser.line());

            int commandType = parser.commandType();
            String arg1 = parser.arg1();
            int arg2 = parser.arg2();

            if (commandType == CommandType.C_ARITHMETIC) {
                codeWriter.writeArithmetic(arg1);
            } else if (commandType == CommandType.C_PUSH || commandType == CommandType.C_POP) {
                codeWriter.writePushPop(commandType, arg1, arg2);
            }
        }
        codeWriter.writeInfiniteLoop();

        inputScanner.close();
        fileWriter.close();
    }

    private Scanner createInputScanner(String filename) {
        FileInputStream fileInputStream = null;
        try {
            fileInputStream = new FileInputStream(filename); 
        } catch (FileNotFoundException e) {
            System.out.println("Error: file not found");
            System.exit(1);
        }
        return new Scanner(fileInputStream);
    }

    private BufferedWriter createFileWriter(String filename) {
        String outputFilename = filename.substring(0, filename.length() - 2) + "asm";
        BufferedWriter fileWriter = null;
        try {
            fileWriter = new BufferedWriter(new FileWriter(outputFilename));
        } catch (IOException e) {
            System.out.println("Error: can't create output file");
            System.exit(1);
        }
        return fileWriter;
    }
}
