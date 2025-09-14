# Assemplex

Assemplex is a simple assembly-like interpreter written in Go. It allows you to execute programs with registers, arithmetic operations, comparisons, jumps, and input/output in a syntax similar to assembly language.

> Note: A precompiled `.exe` file is included in the repository. You can use it as either `assemplex.exe` or `asp.exe`.

## Features

* **16 General Purpose Registers:** r1 to r16
* **Arithmetic Operations:** `MOV`, `ADD`, `SUB`, `MUL`, `DIV`, `MOD`, `INC`, `DEC`, `NEG`
* **Comparison Operations:** `LT`, `LE`, `GT`, `GE`, `EQ`, `NE` (set a comparison flag)
* **Conditional Jumps:** `JZ`, `JNZ` (check the comparison flag)
* **Unconditional Jump:** `JMP` (jump to label)
* **Input/Output:** `INPUT`, `PRINT`
* **Program Flow:** Labels, loops, and branching
* **HALT:** Stop the program execution

## Syntax

* **MOV dest, src**: Move value from src (register or number) to dest register.
* **ADD dest, src**: Add src to dest.
* **SUB dest, src**: Subtract src from dest.
* **MUL dest, src**: Multiply dest by src.
* **DIV dest, src**: Divide dest by src.
* **MOD dest, src**: dest modulo src.
* **INC dest / DEC dest / NEG dest**: Increment, decrement, or negate register.
* **PRINT reg**: Print value of register.
* **INPUT reg**: Read value from user input into register.
* **Comparisons:** `LT`, `LE`, `GT`, `GE`, `EQ`, `NE` set the comparison flag.
* **JZ label / JNZ label**: Jump if comparison flag is zero / non-zero.
* **JMP label**: Unconditional jump.
* **Labels:** Defined as `label:`
* **HALT:** Stop program execution.
* **Comments:** Start with `//`

## Example Program

```asm
MOV r1, 1
MOV r2, 1000000

loop:
PRINT r1
MUL r1, 2
EQ r1, r2
JNZ end
JMP loop

end:
HALT
```

This program doubles `r1` each iteration and prints it until it reaches `r2`.

## Running Programs

If you want to use the precompiled executable:

```bash
assemplex.exe example.asp
# or
asp.exe example.asp
```

Or compile the Go interpreter:

```bash
go build -o assemplex main.go
./assemplex example.asp
```

## Notes

* All numbers are treated as `float64`.
* Comparison instructions affect only the internal comparison flag, not the registers.
* Use labels for loops and branching.
* Comments start with `//`.

## License

MIT License

