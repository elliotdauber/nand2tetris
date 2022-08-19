
all:
	g++ -c -o parser.o src/parser.cpp -I include
	g++ -c -o code.o --std=c++11 src/code.cpp -I include 
	g++ -c -o symboltable.o --std=c++11 src/symboltable.cpp -I include 
	g++ -c -o assembler.o --std=c++11 src/assembler.cpp -I include 
	g++ -c -o main.o src/main.cpp -I include 
	g++ -g -o assembler parser.o code.o symboltable.o assembler.o main.o

	rm -f *.o

clean:
	rm -f *.o