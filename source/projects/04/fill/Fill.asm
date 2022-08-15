// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

//set the initial values
@8192
D = A
@numPixelRegisters
M = D

@previousInput
M = -1

(LOOP)

//get the keyboard input
@KBD
D = M

//store keyboard input as a variable
@input
M = D

//if input == previousInput, goto LOOP (no need to draw)
@previousInput
D = D - M
@LOOP
D;JEQ

//save previousInput for the next iteration
@input
D = M 
@previousInput
M = D

//initialize 0 for the drawGroup
@i
M = 0

(DRAWLOOP)

//if i == numPixelRegisters, goto drawend
@i
D = M
@numPixelRegisters
D = D - M
@DRAWEND
D;JEQ

//load the input
@input
D = M

//jump to WHITE if input == 0
@WHITE
D;JEQ

//set as all 1s (represented as -1)
@SCREEN
D = A
@i
A = D + M
M = -1

//jump to end of loop
@DRAWLOOPEND
0;JMP

(WHITE)

//set to all 0s
@SCREEN
D = A
@i
A = D + M
M = 0

(DRAWLOOPEND)

//increment i
@i
M = M + 1

//return to beginning of draw loop
@DRAWLOOP
0;JMP

//just restarts the outer loop 
(DRAWEND)
@LOOP
0;JMP