import sys
import os

# --------------------------
# Exit codes
# --------------------------
EXIT_OK        = 0
EXIT_DIVZERO   = 1
EXIT_MODZERO   = 2
EXIT_STACKERR  = 3
EXIT_BADRET    = 4
EXIT_BADFUNC   = 5
EXIT_RUNTIME   = 6
EXIT_LIBERROR  = 7

# --------------------------
# Opcodes
# --------------------------
OPCODES = {
    'PUSH': 1, 'STORE': 2, 'LOAD': 3, 'PRINT': 4, 'HALT': 5,
    'ADD': 6, 'SUB': 7, 'MUL': 8, 'DIV': 9, 'MOD': 10,
    'EQ': 11, 'NEQ': 12, 'GT': 13, 'LT': 14, 'GE': 15, 'LE': 16,
    'AND': 17, 'OR': 18, 'NOT': 19,
    'DUP': 20, 'SWAP': 21, 'DROP': 22,
    'JMP': 23, 'JZ': 24, 'JNZ': 25,
    'CALL': 26, 'RET': 27, 'CALLFN': 28,
    'TOINT': 29, 'TOSTR': 30, 'LEN': 31, 'ROUND': 32,
    'LOADLIB': 33, 'IN': 34
}

# --------------------------
# Parse .asp source to bytecode
# --------------------------
def parse(source):
    code = []
    vars = {}
    labels = {}
    next_slot = 0
    lines = []

    for line in source.strip().split("\n"):
        line = line.split(";")[0].strip()
        if line:
            lines.append(line)

    pc = 0
    for line in lines:
        if line.startswith("FUNC "):
            func_name = line.split()[1].rstrip(":")
            labels[func_name] = pc
        elif not line.endswith(":"):
            pc += 1

    for line in lines:
        if line.startswith("FUNC "):
            continue

        instr, *arg = line.split(maxsplit=1)
        instr = instr.upper()

        if instr not in OPCODES:
            raise ValueError(f"Unknown instruction: {instr}")

        opcode = OPCODES[instr]

        if arg:
            val = arg[0].strip()
            try:
                arg_val = int(val)
            except ValueError:
                try:
                    arg_val = float(val)
                except ValueError:
                    if val.startswith('"') and val.endswith('"'):
                        arg_val = val[1:-1]
                    elif instr == "PUSH" and val in labels:
                        arg_val = ("FUNC", labels[val])
                    elif val in labels:
                        arg_val = labels[val]
                    elif instr in {"STORE", "LOAD", "IN"}:
                        if val not in vars:
                            vars[val] = next_slot
                            next_slot += 1
                        arg_val = vars[val]
                    else:
                        arg_val = val
            code.append((opcode,arg_val))
        else:
            code.append((opcode,None))

    return code, next_slot, labels

# --------------------------
# VM Interpreter
# --------------------------
def interpreter(file_path):
    if not os.path.exists(file_path):
        print(f"File not found: {file_path}")
        sys.exit(EXIT_RUNTIME)

    with open(file_path, "r") as f:
        source = f.read()

    code, memsize, labels = parse(source)
    stack = []
    mem = [0] * memsize
    call_stack = []
    pc = 0
    exit_code = EXIT_OK
    loaded_libraries = set()

    def popn(n):
        if len(stack) < n:
            sys.exit(EXIT_STACKERR)
        vals = [stack.pop() for _ in range(n)]
        return vals[::-1]

    def load_library(lib_file):
        if not lib_file.endswith(".asp"):
            lib_file += ".asp"
        if lib_file in loaded_libraries:
            return
        if not os.path.exists(lib_file):
            sys.exit(EXIT_LIBERROR)
        with open(lib_file, "r") as f:
            lib_source = f.read()
        lib_code, lib_memsize, lib_labels = parse(lib_source)
        offset = len(code)
        for name, addr in lib_labels.items():
            labels[name] = addr + offset
        code.extend(lib_code)
        loaded_libraries.add(lib_file)

    while pc < len(code):
        op, arg = code[pc]
        pc += 1

        try:
            # --- Core ---
            if op == OPCODES['PUSH']:
                stack.append(arg)
            elif op == OPCODES['STORE']:
                mem[arg] = stack.pop()
            elif op == OPCODES['LOAD']:
                stack.append(mem[arg])
            elif op == OPCODES['PRINT']:
                print(stack.pop())
            elif op == OPCODES['HALT']:
                sys.exit(exit_code)
            elif op == OPCODES['IN']:
                prompt = ""
                if arg is not None:
                    prompt = f"{arg}: "
                inp = input(prompt)
                try:
                    val = int(inp) if inp.isdigit() or (inp.startswith('-') and inp[1:].isdigit()) else float(inp)
                except ValueError:
                    val = inp
                if arg is not None:
                    mem[arg] = val
                else:
                    stack.append(val)

            # --- Math ---
            elif op == OPCODES['ADD']:
                a, b = popn(2); stack.append(a + b)
            elif op == OPCODES['SUB']:
                a, b = popn(2); stack.append(a - b)
            elif op == OPCODES['MUL']:
                a, b = popn(2); stack.append(a * b)
            elif op == OPCODES['DIV']:
                a, b = popn(2)
                if b == 0: sys.exit(EXIT_DIVZERO)
                stack.append(a / b)
            elif op == OPCODES['MOD']:
                a, b = popn(2)
                if b == 0: sys.exit(EXIT_MODZERO)
                stack.append(a % b)

            # --- Comparison ---
            elif op == OPCODES['EQ']:
                a, b = popn(2); stack.append(1 if a == b else 0)
            elif op == OPCODES['NEQ']:
                a, b = popn(2); stack.append(1 if a != b else 0)
            elif op == OPCODES['GT']:
                a, b = popn(2); stack.append(1 if a > b else 0)
            elif op == OPCODES['LT']:
                a, b = popn(2); stack.append(1 if a < b else 0)
            elif op == OPCODES['GE']:
                a, b = popn(2); stack.append(1 if a >= b else 0)
            elif op == OPCODES['LE']:
                a, b = popn(2); stack.append(1 if a <= b else 0)

            # --- Boolean ---
            elif op == OPCODES['AND']:
                a, b = popn(2); stack.append(1 if a and b else 0)
            elif op == OPCODES['OR']:
                a, b = popn(2); stack.append(1 if a or b else 0)
            elif op == OPCODES['NOT']:
                a, = popn(1); stack.append(0 if a else 1)

            # --- Stack ---
            elif op == OPCODES['DUP']:
                stack.append(stack[-1])
            elif op == OPCODES['SWAP']:
                a, b = popn(2); stack.extend([b, a])
            elif op == OPCODES['DROP']:
                stack.pop()

            # --- Control Flow ---
            elif op == OPCODES['JMP']:
                pc = arg
            elif op == OPCODES['JZ']:
                a, = popn(1)
                if a == 0: pc = arg
            elif op == OPCODES['JNZ']:
                a, = popn(1)
                if a != 0: pc = arg
            elif op == OPCODES['CALL']:
                call_stack.append(pc); pc = arg
            elif op == OPCODES['RET']:
                if not call_stack: sys.exit(EXIT_BADRET)
                pc = call_stack.pop()
            elif op == OPCODES['CALLFN']:
                func_ref = stack.pop()
                if isinstance(func_ref, tuple) and func_ref[0] == "FUNC":
                    call_stack.append(pc); pc = func_ref[1]
                elif isinstance(func_ref, int):
                    call_stack.append(pc); pc = func_ref
                else:
                    sys.exit(EXIT_BADFUNC)

            # --- Utilities ---
            elif op == OPCODES['TOINT']:
                a, = popn(1); stack.append(int(a))
            elif op == OPCODES['TOSTR']:
                a, = popn(1); stack.append(str(a))
            elif op == OPCODES['LEN']:
                a, = popn(1); stack.append(len(str(a)))
            elif op == OPCODES['ROUND']:
                a, = popn(1); stack.append(round(a))

            # --- Libraries ---
            elif op == OPCODES['LOADLIB']:
                load_library(arg)

        except SystemExit:
            raise
        except:
            sys.exit(EXIT_RUNTIME)

# --------------------------
# Run VM from command line
# --------------------------
if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python vm.py main.asp")
        sys.exit(EXIT_RUNTIME)
    interpreter(sys.argv[1])
