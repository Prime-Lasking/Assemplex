from math import *
OPCODES = {
    'PUSH': 1, 'ADD': 2, 'STORE': 3, 'LOAD': 4, 'PRINT': 5, 'HALT': 6,
    'SUB': 7, 'MUL': 8, 'DIV': 9, 'EXP': 10, 'SQRT': 11, 'JMP': 12,
    'JZ': 13, 'JNZ': 14, 'TOINT': 15, 'TOSTR': 16, 'LEN': 17, 'EQ': 18,
    'GT': 19, 'LT': 20, 'IN': 21, 'MOD': 22, 'AND': 23, 'OR': 24,
    'NOT': 25, 'IF': 26, 'FLOOR': 27, 'INC': 28, 'DEC': 29,'CEIL': 30,
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
        opcode = OPCODES[instr]
        if arg:
            val = arg[0].strip()
            if val.lstrip('-').isdigit():
                arg_val = int(val)
            elif val.startswith('"') and val.endswith('"'):
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

    while i < len(code):
        op, arg = code[i]
        i += 1

        if op == 1:  # PUSH
            stack.append(arg)

        elif op == 2:  # ADD
            b = stack.pop()
            a = stack.pop()
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
            b = stack.pop()
            a = stack.pop()
            stack.append(a - b)

        elif op == 8:  # MUL
            b = stack.pop()
            a = stack.pop()
            stack.append(a * b)

        elif op == 9:  # DIV
            b = stack.pop()
            a = stack.pop()
            stack.append(a / b)

        elif op == 10:  # EXP
            b = stack.pop()
            a = stack.pop()
            stack.append(a ** b)

        elif op == 11:  # SQRT
            a = stack.pop()
            stack.append(a ** 0.5)

        elif op == 12:  # JMP
            i = arg

        elif op == 13:  # JZ
            if stack.pop() == 0:
                i = arg

        elif op == 14:  # JNZ
            if stack.pop() != 0:
                i = arg

        elif op == 15:  # TOINT
            stack.append(int(stack.pop()))

        elif op == 16:  # TOSTR
            stack.append(str(stack.pop()))

        elif op == 17:  # LEN
            stack.append(len(str(stack.pop())))

        elif op == 18:  # EQ
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a == b else 0)

        elif op == 19:  # GT
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a > b else 0)

        elif op == 20:  # LT
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a < b else 0)

        elif op == 21:  # IN
            mem[arg] = input("Enter value: ")

        elif op == 22:  # MOD
            b = stack.pop()
            a = stack.pop()
            stack.append(a % b)

        elif op == 23:  # AND
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a and b else 0)

        elif op == 24:  # OR
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a or b else 0)

        elif op == 25:  # NOT
            a = stack.pop()
            stack.append(0 if a else 1)

        elif op == 26:  # IF
            if not stack.pop():
                i = arg

        elif op == 27:  # FLOOR
            a = stack.pop()
            stack.append(floor(a))

        elif op == 28:  # INC
            a = stack.pop()
            if isinstance(a, str):
                raise TypeError("Cannot use INC on string")
            stack.append(a + 1)

        elif op == 29:  # DEC
            a = stack.pop()
            if isinstance(a, str):
                raise TypeError("Cannot use DEC on string")
            stack.append(a - 1)
        elif op == 30: # CEIL
            a = stack.pop()
            stack.append(ceil(a))
