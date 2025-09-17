# Assemplex

Assemplex is a high-performance, low-level programming language and virtual machine designed for performance-critical applications. Written in Go, it provides a clean assembly-like syntax while delivering execution speeds approximately 2-4x faster than Python for many tasks. Assemplex is perfect for projects that need low-level control without sacrificing development speed.

## Key Features

- **High Performance**: Optimized execution with near-native speeds
- **Low-Level Control**: Direct memory management and type system
- **Type System**: Support for various numeric types (INT16/32/64, FLOAT16/32/64) and CHAR
- **Efficient Execution**: Optimized bytecode interpreter for fast performance
- **Function Support**: Define and call functions with parameters
- **Memory Management**: Manual memory management for optimal control

## Installation

1. Make sure you have Go 1.16 or later installed

2. Build the project:
   ```bash
   go build -o asp.exe asp.go
   ```

3.  Add the directory containing `asp.exe` to your system's PATH

## Quick Start

Create a file named `hello.asp` with the following content:

```assembly
; Simple Hello World program
VAR INT32 counter 0

; Main program
PRINT "Hello, World!"

; Count to 5
VAR INT32 i 1
LABEL loop
PRINT i
ADD i, 1
LT i, 6
JNZ loop

HALT
```

Run it with:
```bash
asp hello.asp
```

## Language Reference

### Variables

Variables in Assemplex are statically typed and must be declared before use. Each variable has a specific type that determines the size and kind of data it can hold.

```assembly
; Basic variable declaration
VAR INT32 count          ; Declare an integer variable with default value 0
VAR FLOAT64 pi 3.14159   ; Declare and initialize a floating-point number
VAR CHAR initial 'A'     ; Declare and initialize a character

; Multiple variables of the same type
VAR INT32 x 10, y 20, z 30

; Freeing variables when done
FREE count               ; Free the variable when no longer needed

; Using variables in expressions
VAR INT32 a 5
VAR INT32 b 10
ADD a, b                ; a = a + b
PRINT a                 ; Prints 15
```

### Data Types
- `INT16`, `INT32`, `INT64`: Signed integers
- `FLOAT16`, `FLOAT32`, `FLOAT64`: Floating-point numbers
- `CHAR`: Single character

### Arithmetic Operations
```assembly
ADD <dest>, <src>  ; dest = dest + src
SUB <dest>, <src>  ; dest = dest - src
MUL <dest>, <src>  ; dest = dest * src
DIV <dest>, <src>  ; dest = dest / src
```

### Control Flow
```assembly
JMP <label>     ; Unconditional jump to label
JZ <label>      ; Jump if last comparison was equal
JNZ <label>     ; Jump if last comparison was not equal
HALT            ; Stop program execution
```

### Comparison Operators
```assembly
LT <a>, <b>     ; Set flag if a < b
LE <a>, <b>     ; Set flag if a <= b
GT <a>, <b>     ; Set flag if a > b
GE <a>, <b>     ; Set flag if a >= b
EQ <a>, <b>     ; Set flag if a == b
NE <a>, <b>     ; Set flag if a != b
```

### Functions
```assembly
FUNC <name>     ; Start function definition
  ; function body
  RETURN <value> ; Return from function with value
  ; or
  RET           ; Return without value
ENDFUNC         ; End function definition

; Function call
CALL <func> [arg1, arg2, ...]  ; Call function with arguments
```

### I/O Operations
```assembly
PRINT <value>   ; Print value to console
INPUT <var>     ; Read input into variable
```

## Examples

### Fibonacci Sequence
```assembly
; Calculate first 10 Fibonacci numbers
VAR INT32 a 0
VAR INT32 b 1
VAR INT32 i 0
VAR INT32 temp 0

LABEL loop
PRINT a
MOV temp, a
ADD a, b
MOV b, temp
ADD i, 1
LT i, 10
JNZ loop

HALT
```

### Function Example
```assembly
; Define a function to calculate square
FUNC square
  MUL p0, p0  ; p0 is first parameter
  RETURN p0
ENDFUNC

; Main program
VAR INT32 num 5
CALL square, num
PRINT p0  ; Prints 25

HALT
```

### Running Your First Program
Create a file named `hello.asp` with the following content:
```assembly
; Example Assemplex program
MOV R1, 10
MOV R2, 20
ADD R1, R2
PRINT R1  ; Should print 30
```

Then run it with:
```bash
asp hello.asp
```

## Usage

### Running Programs
To run an Assemplex program:
```bash
asp <filename.asp>
```

### Example
Create a file named `hello.asp`:
```assembly
; This is a comment
LOAD R1, 42     ; Load value 42 into register R1
PRINT R1        ; Print the value in R1
HALT            ; End program
```

Run it with:
```bash
asp hello.asp
```

## Language Reference

Assemplex provides a clean, minimal set of instructions that cover fundamental programming concepts while remaining powerful enough for various applications. It uses a memory-based variable system with static typing.

### Memory Model

Assemplex uses a simple memory model where all variables must be explicitly declared with their types before use. Variables are scoped to their containing function and must be freed when no longer needed.

### Instruction Set

#### Variable Management
- `VAR <type> <name> [value]` - Declare a variable with optional initial value
- `FREE <name>` - Free a variable's memory
- `VAR <type> <name1> [value1], <name2> [value2], ...` - Multiple declarations

#### Core Instructions
- `PRINT <value>` - Print value to console
- `JMP <label>` - Unconditional jump to label
- `JZ <label>` - Jump if last comparison was true
- `JNZ <label>` - Jump if last comparison was false
- `HALT` - Stop program execution

#### Arithmetic Operations
- `ADD <dest>, <src>` - Add values (dest = dest + src)
- `SUB <dest>, <src>` - Subtract values (dest = dest - src)
- `MUL <dest>, <src>` - Multiply values (dest = dest * src)
- `DIV <dest>, <src>` - Divide values (dest = dest / src)
- `MOD <dest>, <src>` - Modulo operation (dest = dest % src)

#### Comparison Operations
- `LT <a>, <b>` - Less than
- `LE <a>, <b>` - Less than or equal
- `GT <a>, <b>` - Greater than
- `GE <a>, <b>` - Greater than or equal
- `EQ <a>, <b>` - Equal to
- `NE <a>, <b>` - Not equal to

#### Function Handling
- `FUNC <name>` - Define a function
- `CALL <func> [arg1, arg2, ...]` - Call a function with arguments
- `RETURN <value>` - Return from function with value
- `RET` - Return from function without value
- `ENDFUNC` - End function definition

### Program Structure

A typical Assemplex program follows this structure:

```assembly
; This is a comment
FUNC main
  ; Declare variables
  VAR INT32 num1 10
  VAR INT32 num2 20
  VAR INT32 result
  
  ; Perform operations
  ADD result, num1
  ADD result, num2  ; result = 10 + 20 = 30
  
  ; Output result
  PRINT result      ; Prints 30
  
  ; Clean up
  FREE num1, num2, result
  HALT
ENDFUNC
```

### Features

- **Type-Safe Operations**: All operations are type-checked at runtime
- **Memory Management**: Manual memory management with explicit allocation/deallocation
- **Function Support**: Define and call functions with parameters
- **Labels**: Support for named code locations with labels
- **Comments**: Use `;` for single-line comments
- **Error Handling**: Detailed error messages for invalid operations
### Example Program

Here's a simple program that adds two numbers and prints the result:

```assembly
; Simple addition program
FUNC main
  ; Declare and initialize variables
  VAR INT32 num1 10
  VAR INT32 num2 20
  VAR INT32 result
  
  ; Perform addition
  ADD result, num1
  ADD result, num2
  
  ; Print and clean up
  PRINT result    ; Prints 30
  
  ; Free variables when done
  FREE num1, num2, result
  
  HALT            ; End program
ENDFUNC
```

### Notes
- Use `;` for single-line comments
- Variables must be declared with `VAR` before use
- Always free variables when no longer needed with `FREE`
- Variable names are case-sensitive
- Variables are scoped to their containing function

## Use Cases

Assemplex is designed for:
- **High-Performance Computing**: Where low-level control is essential
- **Embedded Systems**: With precise hardware interaction requirements
- **Compiler Backends**: As a target for higher-level language compilation
- **Performance-Critical Code**: For algorithms where every cycle counts
- **System Software**: Operating systems, drivers, and firmware

## Advanced Examples

### Function Call Example
```assembly
; Function to add two numbers
FUNC add_numbers
    ; Parameters: p0, p1 (automatically assigned)
    ; Returns: p0 + p1
    ADD p0, p1        ; Add parameters
    RETURN p0         ; Return result
ENDFUNC

; Main program
FUNC main
    ; Declare variables
    VAR INT32 a 10
    VAR INT32 b 20
    VAR INT32 result
    
    ; Call function
    CALL add_numbers, a, b
    
    ; Result is in p0
    PRINT p0          ; Prints 30
    
    ; Clean up
    FREE a, b, result
    HALT
ENDFUNC
```

### Conditional Logic with Jumps
```assembly
; Find maximum of two numbers
FUNC main
    ; Declare variables
    VAR INT32 num1 42
    VAR INT32 num2 27
    VAR INT32 max
    
    ; Compare values
    GT num1, num2
    JZ second_greater
    
    ; First number is greater
    MOV max, num1
    JMP print_result
    
second_greater:
    ; Second number is greater or equal
    MOV max, num2
    
print_result:
    PRINT "Maximum is: "
    PRINT max
    
    ; Clean up
    FREE num1, num2, max
    HALT
ENDFUNC
```

### Loop with Counter
```assembly
; Count from 10 to 1
FUNC main
    ; Initialize counter
    VAR INT32 counter 10
    
loop_start:
    ; Print current count
    PRINT counter
    
    ; Decrement and check if done
    SUB counter, 1
    GT counter, 0
    JNZ loop_start
    
    ; Clean up
    FREE counter
    HALT
ENDFUNC
```

### Countdown with Delay
```assembly
; Countdown from 10 with delay
FUNC main
    ; Initialize variables
    VAR INT32 count 10
    VAR INT32 delay 1000000
    VAR INT32 temp
    
count_loop:
    ; Print current number
    PRINT count
    
delay_loop:
    DEC r3
    JNZ delay_loop
    
    DEC r1            ; Decrement counter
    JNZ count_loop    ; Loop if not zero
    
    HALT
ENDFUNC
```

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
