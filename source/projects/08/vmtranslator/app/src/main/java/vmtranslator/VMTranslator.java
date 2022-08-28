package vmtranslator;

import java.io.*;  
import java.util.Scanner;

public class VMTranslator {
    private void translateSingleFile(CodeWriter codeWriter, String filename) {
        String filenameStripped = filename.substring(filename.lastIndexOf("/") + 1, filename.length() - ".vm".length());
        codeWriter.setFileName(filenameStripped);

        Scanner inputScanner = createInputScanner(filename);
        Parser parser = new Parser(inputScanner);

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
            } else if (commandType == CommandType.C_FUNCTION) {
                codeWriter.writeFunction(arg1, arg2);
            } else if (commandType == CommandType.C_CALL) {
                codeWriter.writeCall(arg1, arg2);
            } else if (commandType == CommandType.C_GOTO) {
                codeWriter.writeGoto(arg1);
            } else if (commandType == CommandType.C_IF) {
                codeWriter.writeIf(arg1);
            } else if (commandType == CommandType.C_LABEL) {
                codeWriter.writeLabel(arg1);
            } else if (commandType == CommandType.C_RETURN) {
                codeWriter.writeReturn();
            }
        }

        inputScanner.close();
    }

    public void translateDir(String dirpath) throws IOException {

        BufferedWriter fileWriter = createFileWriterFromDir(dirpath);
        CodeWriter codeWriter = new CodeWriter(fileWriter); 
        codeWriter.writeBootstrapCode();
        
        File dir = new File(dirpath);
        File[] dirFiles = dir.listFiles();
        for (File file : dirFiles) {
            String filename = file.getName();
            if (isVmFile(filename)) {
                translateSingleFile(codeWriter, dirpath + "/" + filename);
            }   
        }
        fileWriter.close();
    }

    public void translateFile(String filename) throws IOException {
        BufferedWriter fileWriter = createFileWriter(filename);
        CodeWriter codeWriter = new CodeWriter(fileWriter);

        codeWriter.writeBootstrapCode();
        translateSingleFile(codeWriter, filename);
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

    private BufferedWriter createFileWriterFromDir(String dirpath) {
        String dirname = dirpath.substring(dirpath.lastIndexOf("/") + 1, dirpath.length());
        String outputFilename = dirpath + "/" + dirname + ".asm";
        BufferedWriter fileWriter = null;
        try {
            fileWriter = new BufferedWriter(new FileWriter(outputFilename));
        } catch (IOException e) {
            System.out.println("Error: can't create output file");
            System.exit(1);
        }
        return fileWriter;
    }

    private boolean isVmFile(String filename) {
        return filename.substring(filename.length() - 3, filename.length()).equals(".vm");
    }
}
