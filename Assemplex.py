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
    'LOADLIB': 33, 'IN': 34, 'DEBUG': 35
}

# --------------------------
# Parse source to bytecode
# --------------------------
def parse(source, vars=None, next_slot=0):
    code = []
    if vars is None:
        vars = {}
    labels = {}
    lines = []

    # strip comments and empty lines
    for line in source.strip().split("\n"):
        line = line.split(";")[0].strip()
        if line:
            lines.append(line)

    # gather labels
    pc = 0
    for line in lines:
        if line.startswith("FUNC "):
            func_name = line.split()[1].rstrip(":")
            labels[func_name] = pc
        elif not line.endswith(":"):
            pc += 1

    # assemble instructions
    for line in lines:
        if line.startswith("FUNC "):
            continue

        parts = line.split(maxsplit=1)
        instr = parts[0].upper()
        arg = parts[1].strip() if len(parts) > 1 else None

        if instr not in OPCODES:
            print(f"Unknown instruction: {instr}")
            sys.exit(EXIT_RUNTIME)

        opcode = OPCODES[instr]

        if arg:
            val = arg

            # LOADLIB arg is string filename
            if instr == "LOADLIB":
                if val.startswith('"') and val.endswith('"'):
                    arg_val = val[1:-1]
                else:
                    arg_val = val
                code.append((opcode, arg_val))
                continue

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
            code.append((opcode, arg_val))
        else:
            code.append((opcode, None))

    return code, next_slot, labels, vars

# --------------------------
# VM Interpreter
# --------------------------
def run_vm(source):
    code, memsize, labels, vars = parse(source)
    stack = []
    mem = [0] * memsize
    call_stack = []
    pc = 0

    def popn(n):
        if len(stack) < n:
            print("Stack underflow!")
            sys.exit(EXIT_STACKERR)
        vals = [stack.pop() for _ in range(n)]
        return vals[::-1]

    while pc < len(code):
        op, arg = code[pc]
        pc += 1

        try:
            if op == OPCODES['PUSH']:
                stack.append(arg)
            elif op == OPCODES['STORE']:
                mem[arg] = stack.pop()
            elif op == OPCODES['LOAD']:
                stack.append(mem[arg])
            elif op == OPCODES['PRINT']:
                print(stack.pop())
            elif op == OPCODES['HALT']:
                sys.exit(EXIT_OK)
            elif op == OPCODES['IN']:
                prompt = ">> "
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
            elif op == OPCODES['ADD']:
                a, b = popn(2)
                stack.append(a + b)
            elif op == OPCODES['SUB']:
                a, b = popn(2)
                stack.append(a - b)
            elif op == OPCODES['MUL']:
                a, b = popn(2)
                stack.append(a * b)
            elif op == OPCODES['DIV']:
                a, b = popn(2)
                if b == 0:
                    print("Error: Division by zero")
                    sys.exit(EXIT_DIVZERO)
                stack.append(a / b)
            elif op == OPCODES['MOD']:
                a, b = popn(2)
                if b == 0:
                    print("Error: Modulo by zero")
                    sys.exit(EXIT_MODZERO)
                stack.append(a % b)
            elif op == OPCODES['EQ']:
                a, b = popn(2)
                stack.append(1 if a == b else 0)
            elif op == OPCODES['NEQ']:
                a, b = popn(2)
                stack.append(1 if a != b else 0)
            elif op == OPCODES['GT']:
                a, b = popn(2)
                stack.append(1 if a > b else 0)
            elif op == OPCODES['LT']:
                a, b = popn(2)
                stack.append(1 if a < b else 0)
            elif op == OPCODES['GE']:
                a, b = popn(2)
                stack.append(1 if a >= b else 0)
            elif op == OPCODES['LE']:
                a, b = popn(2)
                stack.append(1 if a <= b else 0)
            elif op == OPCODES['AND']:
                a, b = popn(2)
                stack.append(1 if a and b else 0)
            elif op == OPCODES['OR']:
                a, b = popn(2)
                stack.append(1 if a or b else 0)
            elif op == OPCODES['NOT']:
                a, = popn(1)
                stack.append(0 if a else 1)
            elif op == OPCODES['DUP']:
                stack.append(stack[-1])
            elif op == OPCODES['SWAP']:
                a, b = popn(2)
                stack.extend([b, a])
            elif op == OPCODES['DROP']:
                stack.pop()
            elif op == OPCODES['JMP']:
                pc = arg
            elif op == OPCODES['JZ']:
                a, = popn(1)
                if a == 0:
                    pc = arg
            elif op == OPCODES['JNZ']:
                a, = popn(1)
                if a != 0:
                    pc = arg
            elif op == OPCODES['CALL']:
                call_stack.append(pc)
                if isinstance(arg, str) and arg in labels:
                    pc = labels[arg]
                else:
                    pc = arg
            elif op == OPCODES['RET']:
                if not call_stack:
                    print("Bad return")
                    sys.exit(EXIT_BADRET)
                pc = call_stack.pop()
            elif op == OPCODES['CALLFN']:
                func_ref = stack.pop()
                if isinstance(func_ref, tuple) and func_ref[0] == "FUNC":
                    call_stack.append(pc)
                    pc = func_ref[1]
                elif isinstance(func_ref, int):
                    call_stack.append(pc)
                    pc = func_ref
                else:
                    print("Bad function call")
                    sys.exit(EXIT_BADFUNC)
            elif op == OPCODES['TOINT']:
                a, = popn(1)
                stack.append(int(a))
            elif op == OPCODES['TOSTR']:
                a, = popn(1)
                stack.append(str(a))
            elif op == OPCODES['LEN']:
                a, = popn(1)
                stack.append(len(str(a)))
            elif op == OPCODES['ROUND']:
                a, = popn(1)
                stack.append(round(a))
            elif op == OPCODES['LOADLIB']:
                # For simplicity, ignore library loading in this interactive mode
                print(f"Loadlib '{arg}' ignored in interactive mode.")
            elif op == OPCODES['DEBUG']:
                print("=== DEBUG STATE ===")
                print(f"PC: {pc}")
                print(f"STACK: {stack}")
                print(f"CALLSTACK: {call_stack}")
                mem_state = {i:v for i,v in enumerate(mem) if v != 0}
                print(f"MEM (nonzero): {mem_state}")
                print("===================")
            else:
                print(f"Unknown opcode {op}")
                sys.exit(EXIT_RUNTIME)

        except SystemExit:
            raise
        except Exception as e:
            print(f"Runtime error: {e}")
            sys.exit(EXIT_RUNTIME)

# --------------------------
# Interactive REPL to input code and run
# --------------------------
def main():
    print("Enter your VM code below.")
    print("Type ':run' on a new line to finish input and run the code.")
    print("Press Ctrl+C to quit.\n")

    lines = []
    while True:
        try:
            line = input('>>> ')
        except KeyboardInterrupt:
            print("\nExiting.")
            sys.exit(EXIT_OK)
        except EOFError:
            print("\nEOF received. Exiting.")
            sys.exit(EXIT_OK)

        if line.strip() == ':run':
            break

        lines.append(line)

    source = "\n".join(lines)
    print("\n--- Running your code ---\n")
    run_vm(source)
    print("\n--- Done ---")

if __name__ == "__main__":
    main()
