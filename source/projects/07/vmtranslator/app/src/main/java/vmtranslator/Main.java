package vmtranslator;

import java.io.IOException;

public class Main {
    public static void main(String[] args) throws IOException {
        if (args.length != 1) {
            System.out.println("Error: please supply a file path");
            System.exit(1);
        }

        String filename = args[0];

        if (!filename.substring(filename.length() - 3, filename.length()).equals(".vm")) {
            System.out.println("Error: must be a .vm file");
            System.exit(1);
        }

        VMTranslator translator = new VMTranslator();
        translator.translate(filename);
    }
}
