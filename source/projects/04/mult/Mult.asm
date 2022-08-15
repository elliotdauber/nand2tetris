// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

//set n = 0
@i
M = 0

//will add up R1 iterations of R0
@sum
M = 0

(LOOP)
//if (i == R1) goto SAVE
@i
D = M
@R1
D = D - M
@SAVE
D;JEQ

//get R0
@R0
D = M

//add R0 to sum
@sum
M = D + M

@i
M = M + 1

@LOOP
0;JMP

(SAVE)

//get sum and store in R2
@sum
D = M

@R2
M = D


//end
(END)
@END
0;JMP