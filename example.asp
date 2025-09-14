MOV r1, 1
MOV r2, 1000000000
loop:
PRINT r1
MUL r1, 2
GT r1, r2
JNZ end
JMP loop 
end:
HALT

