To run the VM Translator, run: ./gradlew run --args="<filepath from vmtranslator/app/>"

NOTE: For some reason, most of the tests in the FunctionCalls directory expect the SP to start at 261, even though the manual says that it should start at 256. This causes most of those tests to fail, when they should be passing. If you run them, the SP (and any data relating to it) may be off by -5, but the data is correct with the assumption that the SP starts at 256.

This may be because the tests expect you to push the normal 5 things onto the stack when calling Sys.init, but this seems pointless because Sys.init never returns.