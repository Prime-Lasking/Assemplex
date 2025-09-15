MOV R1, 400000
EQ r1, 0
JNZ fib0

EQ r1, 1
JNZ fib1

MOV r2, 0       
MOV r3, 1       
MOV r4, 2       

fib_loop:
ADD r5, r2    
MOV r2, r3     
MOV r3, r5      

INC r4
EQ r4, r1
JNZ fib_end
JMP fib_loop

fib0:
PRINT r2
HALT

fib1:
PRINT r3
HALT

fib_end:
PRINT r3
HALT
