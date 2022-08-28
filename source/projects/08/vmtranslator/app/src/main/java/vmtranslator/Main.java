package vmtranslator;

import java.io.IOException;

public class Main {
    public static void main(String[] args) throws IOException {
        if (args.length != 1) {
            System.out.println("Error: please supply a file path");
            System.exit(1);
        }

        String filepath = args[0];
        VMTranslator translator = new VMTranslator();

        if (filepath.substring(filepath.length() - 3, filepath.length()).equals(".vm")) {
            translator.translateFile(filepath);
        } else {
            translator.translateDir(filepath);
        }
    }
}
