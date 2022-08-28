package vmtranslator;

import java.util.List;

public class CommandType {
    public static int C_ARITHMETIC = 0;
    public static int C_PUSH = 1;
    public static int C_POP = 2;
    public static int C_LABEL = 3;
    public static int C_GOTO = 4;
    public static int C_IF = 5;
    public static int C_FUNCTION = 6;
    public static int C_RETURN = 7;
    public static int C_CALL = 8;

    private static List<String> _arithCommands = List.of(
        "add", "sub", "and", "or", "neg", "not", "eq", "lt", "gt"
    );

    public static int fromString(String command) {
        if (_arithCommands.contains(command)) {
            return C_ARITHMETIC;
        } else if (command.equals("push")) {
            return C_PUSH;
        } else if (command.equals("pop")) {
            return C_POP;
        } else if (command.equals("call")) {
            return C_CALL;
        } else if (command.equals("function")) {
            return C_FUNCTION;
        } else if (command.equals("return")) {
            return C_RETURN;
        } else if (command.equals("goto")) {
            return C_GOTO;
        } else if (command.equals("if-goto")) {
            return C_IF;
        } else if (command.equals("label")) {
            return C_LABEL;
        }
        return -1;
    }
}
