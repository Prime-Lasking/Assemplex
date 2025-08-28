from math import *
import time

OPCODES = {
    'PUSH': 1, 'ADD': 2, 'STORE': 3, 'LOAD': 4, 'PRINT': 5, 'HALT': 6,
    'SUB': 7, 'MUL': 8, 'DIV': 9, 'EXP': 10, 'SQRT': 11, 'JMP': 12,
    'JZ': 13, 'JNZ': 14, 'TOINT': 15, 'TOSTR': 16, 'LEN': 17, 'EQ': 18,
    'GT': 19, 'LT': 20, 'IN': 21, 'MOD': 22, 'AND': 23, 'OR': 24,
    'NOT': 25, 'IF': 26, 'FLOOR': 27, 'INC': 28, 'DEC': 29, 'CEIL': 30,
    'TIME': 31, 'DUP': 32, 'DEBUG': 33
}

def interpreter(asm_source):
    # --- Assemble ---
    code = []
    vars = {}
    labels = {}
    next_slot = 0
    lines = []

    for line in asm_source.strip().split('\n'):
        line = line.split(';')[0].strip()
        if line:
            lines.append(line)

    # First pass: label positions
    pc = 0
    for line in lines:
        if line.endswith(':'):
            labels[line[:-1]] = pc
        else:
            instr, *arg = line.split(maxsplit=1)
            pc += 1

    # Second pass: generate bytecode as (opcode, arg)
    for line in lines:
        if line.endswith(':'):
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
                    elif val in labels:
                        arg_val = labels[val]
                    elif instr in {'STORE', 'LOAD', 'IN', 'INC', 'DEC'}:
                        if val not in vars:
                            vars[val] = next_slot
                            next_slot += 1
                        arg_val = vars[val]
                    else:
                        arg_val = val
            code.append((opcode, arg_val))
        else:
            code.append((opcode, None))

    # --- Run VM ---
    stack = []
    mem = [0] * next_slot
    i = 0

    def popn(n):
        if len(stack) < n:
            raise IndexError(f"Stack underflow: tried to pop {n} item(s) from stack of size {len(stack)}")
        return [stack.pop() for _ in range(n)][::-1]

    while i < len(code):
        op, arg = code[i]
        i += 1

        if op == 1:  # PUSH
            stack.append(arg)

        elif op == 2:  # ADD
            a, b = popn(2)
            stack.append(a + b)

        elif op == 3:  # STORE
            mem[arg] = stack.pop()

        elif op == 4:  # LOAD
            stack.append(mem[arg])

        elif op == 5:  # PRINT
            print(stack.pop())

        elif op == 6:  # HALT
            break

        elif op == 7:  # SUB
            a, b = popn(2)
            stack.append(a - b)

        elif op == 8:  # MUL
            a, b = popn(2)
            stack.append(a * b)

        elif op == 9:  # DIV
            a, b = popn(2)
            if b == 0:
                raise ZeroDivisionError("Division by zero")
            stack.append(a / b)

        elif op == 10:  # EXP
            a, b = popn(2)
            stack.append(a ** b)

        elif op == 11:  # SQRT
            a, = popn(1)
            stack.append(a ** 0.5)

        elif op == 12:  # JMP
            i = arg

        elif op == 13:  # JZ
            a, = popn(1)
            if a == 0:
                i = arg

        elif op == 14:  # JNZ
            a, = popn(1)
            if a != 0:
                i = arg

        elif op == 15:  # TOINT
            a, = popn(1)
            stack.append(int(a))

        elif op == 16:  # TOSTR
            a, = popn(1)
            stack.append(str(a))

        elif op == 17:  # LEN
            a, = popn(1)
            stack.append(len(str(a)))

        elif op == 18:  # EQ
            a, b = popn(2)
            stack.append(1 if a == b else 0)

        elif op == 19:  # GT
            a, b = popn(2)
            stack.append(1 if a > b else 0)

        elif op == 20:  # LT
            a, b = popn(2)
            stack.append(1 if a < b else 0)

        elif op == 21:  # IN
            inp = input("Enter value: ")
            try:
                val = int(inp) if inp.isdigit() or (inp.startswith('-') and inp[1:].isdigit()) else float(inp)
            except ValueError:
                val = inp
            mem[arg] = val

        elif op == 22:  # MOD
            a, b = popn(2)
            if b == 0:
                raise ZeroDivisionError("Modulo by zero")
            stack.append(a % b)

        elif op == 23:  # AND
            a, b = popn(2)
            stack.append(1 if a and b else 0)

        elif op == 24:  # OR
            a, b = popn(2)
            stack.append(1 if a or b else 0)

        elif op == 25:  # NOT
            a, = popn(1)
            stack.append(0 if a else 1)

        elif op == 26:  # IF
            a, = popn(1)
            if not a:
                i = arg

        elif op == 27:  # FLOOR
            a, = popn(1)
            stack.append(floor(a))

        elif op == 28:  # INC
            a, = popn(1)
            if isinstance(a, str):
                raise TypeError("Cannot use INC on string")
            stack.append(a + 1)

        elif op == 29:  # DEC
            a, = popn(1)
            if isinstance(a, str):
                raise TypeError("Cannot use DEC on string")
            stack.append(a - 1)

        elif op == 30:  # CEIL
            a, = popn(1)
            stack.append(ceil(a))

        elif op == 31:  # TIME
            a, = popn(1)
            time.sleep(a)

        elif op == 32:  # DUP
            if not stack:
                raise IndexError("Stack underflow: cannot DUP on empty stack")
            stack.append(stack[-1])
        elif op == 33:  # DEBUG
            print("=== DEBUG ===")
            print("STACK:", stack)
            print("MEM:", mem)
            print("=============")

