; Factorial of 5
PUSH 5
STORE n
PUSH 1
STORE result

loop_start:
LOAD n
JZ end
LOAD result
LOAD n
MUL
STORE result
LOAD n
DEC
STORE n
JMP loop_start

end:
LOAD result
PRINT
HALT
