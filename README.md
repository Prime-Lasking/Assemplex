# Assemplex


Assemplex is a high-performance assembly language implementation designed for modern software development. Significantly faster than interpreted languages like Python, Assemplex provides direct hardware-level control with support for variable bit-width registers, modular code organization, and efficient memory management. Built for real-world applications where performance is non-negotiable, Assemplex delivers execution speeds that can be orders of magnitude faster than Python for compute-intensive tasks.

## Key Features

- **Blazing Fast**: Native execution outperforms Python by 10-100x for compute-bound tasks
- **Heterogeneous Registers**: 16-bit to 128-bit registers with variable operation costs
- **Modular Code**: Support for includes, imports, and function calls
- **Arbitrary Precision Math**: Built-in support for large integers with minimal overhead
- **Cross-Platform**: Runs anywhere Go is supported
- **Performance Optimized**: Cycle counting and instruction-level optimization
- **Extensible**: Easy to add custom instructions and operations for specialized tasks

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Language Reference](#language-reference)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Supported Platforms

Assemplex is built with Go's excellent cross-platform support and runs natively on:

### Operating Systems
- **Windows** (7/8/10/11, both 32-bit and 64-bit)
- **macOS** (10.13 High Sierra and later, Intel and Apple Silicon)
- **Linux** (Most distributions with glibc 2.17+ or musl)
- **FreeBSD** (12.0 and later)
- **Android** (via Termux)
- **iOS/iPadOS** (via iSH or similar terminal emulators)

### Architectures
- x86 (32-bit and 64-bit)
- ARM (32-bit and 64-bit)
- ARM64 (Apple Silicon)
- MIPS (experimental)
- RISC-V (experimental)

## Installation

### Method 1: Using Go (Recommended for Developers)

#### Prerequisites
- Go 1.16 or higher
- Git (for cloning the repository)

#### Installation Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Prime-Lasking/Assemplex.git
   cd Assemplex
   ```

2. **Run the setup script**:
   ```bash
   chmod +x setup_asp.sh
   ./setup_asp.sh
   ```

3. **Verify installation**:
   ```bash
   asp --version
   ```

> **Performance Note**: While both methods provide the same functionality, the Go-based installation allows for platform-specific optimizations during compilation. The pre-built binaries are compiled with conservative settings to ensure maximum compatibility across different systems.

### Method 2: Binary-Only Installation (No Go Required)

#### Prerequisites
- A supported operating system (see Supported Platforms)
- Basic terminal/shell access

#### Installation Steps

1. **Download the latest binary** from the [GitHub Releases](https://github.com/Prime-Lasking/Assemplex/releases) page.
   - For Windows: Download the `.zip` file
   - For Linux/macOS: Download the `.tar.gz` file

2. **Extract the binary** to a directory in your PATH:
   - **Linux/macOS**:
     ```bash
     mkdir -p ~/bin
     tar -xzf assemplex-*.tar.gz -C ~/bin/
     ```
   - **Windows**:
     1. Create a directory for binaries (e.g., `C:\bin`)
     2. Extract the ZIP file to this directory

3. **Add to PATH** (if not already):
   - **Linux/macOS**: Add to `~/.bashrc` or `~/.zshrc`:
     ```bash
     export PATH="$HOME/bin:$PATH"
     ```
   - **Windows**:
     1. Open System Properties > Advanced > Environment Variables
     2. Add `C:\bin` to the PATH variable

4. **Verify installation**:
   ```bash
   asp --version
   ```

### Verifying Your Installation
After installation, open a new terminal and run:
```bash
asp --help
```
This should display the help message with available commands and options.

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

Assemplex features a modern register architecture with different bit-widths and performance characteristics:

#### General Purpose Registers
- `r1-r6`: 16-bit registers (1 cycle operations)
- `r7-r10`: 32-bit registers (2 cycle operations)
- `r11-r13`: 64-bit registers (4 cycle operations)
- `r14-r16`: 128-bit registers (8 cycle operations)

#### Special Registers
- `sp`: Stack pointer
- `pc`: Program counter

### Memory Model
- **Global Scope**: Direct memory access for high-performance operations
- **Stack-Based**: Efficient function calls and local variable management
- **Modular Memory**: Support for includes and external module linking

### Instruction Set

#### Data Movement
- `LOAD Rx, value` - Load immediate value or memory into register
- `MOVE Rx, Ry` - Copy value between registers
- `PUSH Rx` - Push register onto stack
- `POP Rx` - Pop value from stack into register

#### Arithmetic
- `ADD Rx, Ry` - Add registers (Rx = Rx + Ry)
- `SUB Rx, Ry` - Subtract registers (Rx = Rx - Ry)
- `MUL Rx, Ry` - Multiply registers (Rx = Rx * Ry)
- `DIV Rx, Ry` - Divide registers (Rx = Rx / Ry)
- `MOD Rx, Ry` - Modulo operation (Rx = Rx % Ry)
- `INC Rx` - Increment register
- `DEC Rx` - Decrement register

#### Control Flow
- `JMP label` - Unconditional jump
- `JEQ label` - Jump if equal (uses CMP result)
- `JNE label` - Jump if not equal
- `JGT label` - Jump if greater than
- `JLT label` - Jump if less than
- `CALL func` - Call function
- `RET` - Return from function
- `HALT` - Stop program execution

#### System
- `CMP Rx, Ry` - Compare two registers (sets flags)
- `PRINT Rx` - Print register value
- `;` - Line comment

## Use Cases

Assemplex is designed for:
- **High-Performance Computing**: Where low-level control is essential
- **Embedded Systems**: With precise hardware interaction requirements
- **Compiler Backends**: As a target for higher-level language compilation
- **Performance-Critical Code**: For algorithms where every cycle counts
- **System Software**: Operating systems, drivers, and firmware

## Advanced Examples

### Function Call with Stack Operations
```assembly
; Main program
LOAD r1, 10      ; Load first argument
LOAD r2, 20      ; Load second argument
CALL add_numbers ; Call function
PRINT r3         ; Print result (30)
HALT

; Function to add two numbers
add_numbers:
    PUSH r1      ; Save registers
    PUSH r2
    ADD r1, r2    ; Add arguments
    MOVE r3, r1   ; Store result in r3
    POP r2        ; Restore registers
    POP r1
    RET           ; Return from function
```

### Conditional Logic
```assembly
; Find maximum of two numbers
LOAD r1, 42      ; First number
LOAD r2, 27      ; Second number
CMP r1, r2       ; Compare values
JGT r1_greater   ; Jump if r1 > r2
MOVE r3, r2      ; r2 is greater
JMP end_compare

r1_greater:
    MOVE r3, r1  ; r1 is greater

end_compare:
    PRINT r3     ; Print maximum value
    HALT
```

### Loop with Counter
```assembly
; Count from 10 to 1
LOAD r1, 10      ; Initialize counter

do_count:
    PRINT r1     ; Print current count
    DEC r1       ; Decrement counter
    CMP r1, 0    ; Check if zero
    JGT do_count ; Loop if greater than zero

HALT
```

### Countdown Loop
```assembly
; Countdown from 10
LOAD R1, 10
loop:
PRINT R1
LOAD R2, 1
SUB R1, R2
JNZ R1, loop
HALT
```

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
