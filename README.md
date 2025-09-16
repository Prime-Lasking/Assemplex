# Assemplex

Assemplex is a high-performance, low-level programming language and virtual machine designed for performance-critical applications. Written in Go, it provides a clean assembly-like syntax while delivering execution speeds approximately 2-4x faster than Python for many tasks. Assemplex is perfect for projects that need low-level control without sacrificing development speed.

## Key Features

- **High Performance**: Approximately 2-4x faster than Python for many workloads
- **Low-Level Control**: Direct access to registers and memory management
- **Simple Syntax**: Clean, assembly-like language that's easy to learn
- **Cross-Platform**: Runs anywhere Go is supported
- **Efficient**: Optimized for performance with minimal overhead

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Language Reference](#language-reference)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## System Requirements

- **Go 1.16 or later** (for building from source)
- **Windows, macOS, or Linux** (for running the binary)

## Installation

### Option 1: Using the Pre-built Binary

1. Download the `asp` binary from the [GitHub Releases](https://github.com/Prime-Lasking/Assemplex/releases) page.

2. Make the binary executable:
   - **Linux/macOS**:
     ```bash
     chmod +x asp
     ```
   - **Windows**: The binary should be ready to use as `asp.exe`

   3. Move the binary to a directory in your system's PATH for global access

### Option 2: Building from Source

1. Make sure you have Go 1.16 or later installed

2. Clone the repository:
   ```bash
   git clone https://github.com/Prime-Lasking/Assemplex.git
   cd Assemplex
   ```

3. Build the project:
   ```bash
   go build -o asp.exe asp.go
   ```

### Verifying Your Installation
Run the following command to verify the installation:
```bash
# On Linux/macOS
./asp --version

# On Windows
asp.exe --version
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
./asp hello.asp
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

Assemplex provides a clean, minimal set of instructions that cover fundamental programming concepts while remaining powerful enough for various applications.

### Register Architecture

Assemplex features a modern register architecture with different bit-widths and performance characteristics. Each register type has its own performance characteristics and operation costs.

#### General Purpose Registers
- `r1-r6`: 16-bit registers (1 cycle operations)
  - 16-bit unsigned integers (0-65,535)
  - Fastest operations, minimal memory usage
- `r7-r10`: 32-bit registers (2 cycle operations)
  - 32-bit unsigned integers (0-4,294,967,295)
  - Good balance of speed and capacity
- `r11-r13`: 64-bit registers (4 cycle operations)
  - 64-bit unsigned integers (0-18,446,744,073,709,551,615)
  - Larger capacity, moderate speed
- `r14-r16`: 128-bit registers (8 cycle operations)
  - 128-bit arbitrary-precision integers (using math/big)
  - Maximum capacity, slower operations

### Instruction Set

#### Core Instructions
- `MOV Rx, Ry` - Copy value between registers
- `MOV Rx, value` - Load immediate value into register
- `PRINT Rx` - Print value of register to console
- `JMP label` - Unconditional jump to label
- `JZ label` - Jump if zero flag is set
- `JNZ label` - Jump if zero flag is not set
- `HALT` - Stop program execution

#### Arithmetic Operations
- `ADD Rx, Ry` - Add registers (Rx = Rx + Ry)
- `SUB Rx, Ry` - Subtract registers (Rx = Rx - Ry)
- `MUL Rx, Ry` - Multiply registers (Rx = Rx * Ry)
- `DIV Rx, Ry` - Divide registers (Rx = Rx / Ry)
- `MOD Rx, Ry` - Modulo operation (Rx = Rx % Ry)
- `NEG Rx` - Negate register value
- `INC Rx` - Increment register by 1
- `DEC Rx` - Decrement register by 1

#### Comparison Operations
- `LT` - Less than (sets zero flag if true)
- `LE` - Less than or equal (sets zero flag if true)
- `GT` - Greater than (sets zero flag if true)
- `GE` - Greater than or equal (sets zero flag if true)
- `EQ` - Equal to (sets zero flag if true)
- `NE` - Not equal to (sets zero flag if true)

#### Function Handling
- `FUNC name` - Define a function
- `CALL func` - Call a function
- `ENDFUNC` - End function definition

### Program Structure

A typical Assemplex program follows this structure:

```assembly
; This is a comment
FUNC main
  MOV r1, 10      ; Load immediate 10 into r1
  MOV r2, 20      ; Load immediate 20 into r2
  ADD r1, r2      ; r1 = r1 + r2
  PRINT r1        ; Should print 30
  HALT            ; End program
ENDFUNC
```

### Features

- **Type-Safe Operations**: Operations respect register bit-widths
- **Cycle Counting**: Tracks execution cycles for performance analysis
- **Function Support**: Define and call functions
- **Labels**: Support for named code locations with labels
- **Comments**: Use `;` for single-line comments
- **Error Handling**: Detailed error messages for invalid operations
### Example Program

Here's a simple program that adds two numbers and prints the result:

```assembly
; Simple addition program
FUNC main
  MOV r1, 10      ; Load first number
  MOV r2, 20      ; Load second number
  ADD r1, r2      ; Add them together
  PRINT r1        ; Print result (30)
  HALT            ; End program
ENDFUNC
```

### Notes
- Use `;` for single-line comments
- Programs must have a `main` function
- All code must be inside functions
- Register operations respect bit-widths automatically

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
    ; Arguments: r1, r2
    ; Returns: r3 = r1 + r2
    MOV r3, r1        ; Copy first argument to r3
    ADD r3, r2        ; Add second argument
    RET               ; Return with result in r3
ENDFUNC

; Main program
FUNC main
    MOV r1, 10        ; First argument
    MOV r2, 20        ; Second argument
    CALL add_numbers   ; Call function
    PRINT r3          ; Print result (30)
    HALT
ENDFUNC
```

### Conditional Logic with Jumps
```assembly
; Find maximum of two numbers
FUNC main
    MOV r1, 42        ; First number
    MOV r2, 27        ; Second number
    
    ; Compare values using subtraction
    MOV r3, r1        ; Copy r1 to r3
    SUB r3, r2        ; r3 = r1 - r2
    JGT r1_greater    ; Jump if r1 > r2
    
    ; r2 is greater or equal
    PRINT r2          ; Print r2
    JMP end_compare
    
r1_greater:
    ; r1 is greater
    PRINT r1          ; Print r1
    
end_compare:
    HALT
ENDFUNC
```

### Loop with Counter
```assembly
; Count from 10 to 1
FUNC main
    MOV r1, 10        ; Initialize counter to 10
    
loop_start:
    PRINT r1          ; Print current count
    DEC r1            ; Decrement counter
    JZ loop_end       ; If zero, exit loop
    JMP loop_start    ; Otherwise, continue loop
    
loop_end:
    HALT              ; End program
ENDFUNC
```

### Countdown with Delay
```assembly
; Countdown from 10 with delay
FUNC main
    MOV r1, 10        ; Initialize counter
    MOV r2, 1000000   ; Delay counter
    
count_loop:
    PRINT r1          ; Print current number
    
    ; Delay loop
    MOV r3, r2
    
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
