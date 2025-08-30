# ⚡ Assemplex

**Assemplex** is a **Pythonic Assembly Language** built on a lightweight stack-based virtual machine.  
It is designed for learning low-level programming concepts while still running on Python.

With Assemplex you can:
- Write simple assembly-like programs with readable opcodes  
- Work with arithmetic, logic, branching, and memory instructions  
- Execute them directly on the provided Python interpreter

---

## 🚀 Features

- Stack-based architecture  
- Human-readable opcodes (e.g., `PUSH`, `ADD`, `PRINT`)  
- Memory operations (`STORE`, `LOAD`)  
- Control flow with jumps and conditionals  
- Type conversion, string operations, and math helpers  
- Easy to extend with new instructions  

---

## 🧾 Language Basics

Assemplex code is written line-by-line using **opcodes**.  
Comments start with `;`.  
Labels end with `:` and can be used as jump targets.

### Example: Hello World

```asm
PUSH "Hello, World!"
PRINT
HALT
```

### Example: Arithmetic

```asm
PUSH 5
PUSH 3
ADD
PRINT    ; prints 8
HALT
```

### Example: Variables

```asm
PUSH 42
STORE x

LOAD x
PRINT     ; prints 42
HALT
```

### Example: Loop

```asm
PUSH 5
STORE counter

loop_start:
LOAD counter
PRINT
LOAD counter
DEC
STORE counter
LOAD counter
JNZ loop_start
HALT
```

---

## 🛠 Opcodes Reference

| Opcode             | Description                                |
|-------------------|--------------------------------------------|
| **PUSH X**         | Push value `X` onto stack                 |
| **ADD**            | Pop two values, push their sum            |
| **SUB**            | Pop two values, push subtraction          |
| **MUL**            | Pop two values, push multiplication       |
| **DIV**            | Pop two values, push division             |
| **MOD**            | Pop two values, push modulo               |
| **STORE var**      | Store top of stack into variable          |
| **LOAD var**       | Load variable value onto stack            |
| **PRINT**          | Pop and print top of stack                |
| **HALT**           | Stop program                              |
| **JMP label**      | Jump to label                             |
| **JZ label**       | Jump if top of stack == 0                 |
| **JNZ label**      | Jump if top of stack != 0                 |
| **CALL label**     | Call a function at label                  |
| **RET**            | Return from function call                 |
| **CALLFN**         | Call function reference from stack        |
| **EQ / NEQ / GT / LT / GE / LE** | Comparison ops, push 1 or 0      |
| **AND / OR / NOT** | Boolean logic                             |
| **DUP / SWAP / DROP** | Stack manipulation                     |
| **TOINT / TOSTR**  | Type conversion                            |
| **LEN / ROUND**    | Utilities for strings and numbers        |
| **LOADLIB**        | Load an external `.asp` library           |
| **IN var**         | Input from user, store into variable     |

---

## 📂 Example: Factorial

```asm
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
```

Output:

```
120
```

---

## 💻 Installing and Running Assemplex on Windows

You can distribute Assemplex as a **Python script** or a **standalone executable**.

### 1️⃣ Using Python Version

1. Install **Python 3** and add it to your PATH.
2. Download and unzip the `ASP-Run` folder anywhere.
3. Open Command Prompt.
4. Run ASP programs with:

```cmd
C:\path\to\ASP-Run\runasp.bat myprogram
```

- Or, if you added `ASP-Run` to your PATH:

```cmd
runasp myprogram
```

- The batch file automatically appends `.asp` if omitted.

---

### 2️⃣ Using Executable Version (No Python Needed)

1. Install **PyInstaller**:

```cmd
pip install pyinstaller
```

2. In the `ASP-Run` folder, create an executable:

```cmd
pyinstaller --onefile runasp.py
```

3. Replace the Python call in `runasp.bat` with:

```bat
"%~dp0runasp.exe" %file%
```

4. Users can now run ASP programs without Python:

```cmd
runasp myprogram
```

---

### 3️⃣ Example Programs

The `examples/` folder contains sample `.asp` programs:

- `hello.asp` – prints "Hello, World!"  
- `test.asp` – demonstrates arithmetic operations  

---

## 📜 License

MIT License © 2025 Prime-Lasking
