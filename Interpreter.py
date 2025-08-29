import time

# --- Opcodes ---
OPCODES = {
    'PUSH': 1, 'ADD': 2, 'STORE': 3, 'LOAD': 4, 'PRINT': 5, 'HALT': 6,
    'SUB': 7, 'MUL': 8, 'DIV': 9, 'EXP': 10, 'SQRT': 11, 'JMP': 12,
    'JZ': 13, 'JNZ': 14, 'TOINT': 15, 'TOSTR': 16, 'LEN': 17, 'EQ': 18,
    'GT': 19, 'LT': 20, 'IN': 21, 'MOD': 22, 'AND': 23, 'OR': 24,
    'NOT': 25, 'IF': 26, 'ROUND': 27, 'INC': 28, 'DEC': 29,
    'TIME': 30, 'DUP': 31, 'DEBUG': 32,
    'CALL': 33, 'RET': 34, 'CALLFN': 35
}

# --- Interpreter ---
def interpreter(asm_source):
    code = []
    vars = {}
    labels = {}
    next_slot = 0
    lines = []

    # sanitize and strip comments
    for line in asm_source.strip().split('\n'):
        line = line.split(';')[0].strip()
        if line:
            lines.append(line)

    # First pass: record labels and FUNC addresses
    pc = 0
    for line in lines:
        if line.startswith("FUNC "):
            func_name = line.split()[1].rstrip(':')
            labels[func_name] = pc
        elif line.endswith(':'):
            raise ValueError("Labels without FUNC are not allowed. Functions must start with FUNC")
        else:
            pc += 1

    # Second pass: generate bytecode
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
                    elif instr == 'PUSH' and val in labels:
                        arg_val = ('FUNC', labels[val])
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
    call_stack = []

    def popn(n):
        if len(stack) < n:
            raise IndexError(f"Stack underflow: tried to pop {n} item(s) from stack of size {len(stack)}")
        return [stack.pop() for _ in range(n)][::-1]

    while i < len(code):
        op, arg = code[i]
        i += 1

        if op == OPCODES['PUSH']:
            stack.append(arg)
        elif op == OPCODES['ADD']:
            a, b = popn(2)
            stack.append(a + b)
        elif op == OPCODES['STORE']:
            mem[arg] = stack.pop()
        elif op == OPCODES['LOAD']:
            stack.append(mem[arg])
        elif op == OPCODES['PRINT']:
            print(stack.pop())
        elif op == OPCODES['HALT']:
            break
        elif op == OPCODES['SUB']:
            a, b = popn(2)
            stack.append(a - b)
        elif op == OPCODES['MUL']:
            a, b = popn(2)
            stack.append(a * b)
        elif op == OPCODES['DIV']:
            a, b = popn(2)
            if b == 0:
                raise ZeroDivisionError("Division by zero")
            stack.append(a / b)
        elif op == OPCODES['EXP']:
            a, b = popn(2)
            stack.append(a ** b)
        elif op == OPCODES['SQRT']:
            a, = popn(1)
            stack.append(a ** 0.5)
        elif op == OPCODES['JMP']:
            i = arg
        elif op == OPCODES['JZ']:
            a, = popn(1)
            if a == 0:
                i = arg
        elif op == OPCODES['JNZ']:
            a, = popn(1)
            if a != 0:
                i = arg
        elif op == OPCODES['TOINT']:
            a, = popn(1)
            stack.append(int(a))
        elif op == OPCODES['TOSTR']:
            a, = popn(1)
            stack.append(str(a))
        elif op == OPCODES['LEN']:
            a, = popn(1)
            stack.append(len(str(a)))
        elif op == OPCODES['EQ']:
            a, b = popn(2)
            stack.append(1 if a == b else 0)
        elif op == OPCODES['GT']:
            a, b = popn(2)
            stack.append(1 if a > b else 0)
        elif op == OPCODES['LT']:
            a, b = popn(2)
            stack.append(1 if a < b else 0)
        elif op == OPCODES['IN']:
            inp = input("Enter value: ")
            try:
                val = int(inp) if inp.isdigit() or (inp.startswith('-') and inp[1:].isdigit()) else float(inp)
            except ValueError:
                val = inp
            mem[arg] = val
        elif op == OPCODES['MOD']:
            a, b = popn(2)
            if b == 0:
                raise ZeroDivisionError("Modulo by zero")
            stack.append(a % b)
        elif op == OPCODES['AND']:
            a, b = popn(2)
            stack.append(1 if a and b else 0)
        elif op == OPCODES['OR']:
            a, b = popn(2)
            stack.append(1 if a or b else 0)
        elif op == OPCODES['NOT']:
            a, = popn(1)
            stack.append(0 if a else 1)
        elif op == OPCODES['IF']:
            a, = popn(1)
            if not a:
                i = arg
        elif op == OPCODES['ROUND']:
            a, = popn(1)
            stack.append(round(a))
        elif op == OPCODES['INC']:
            a, = popn(1)
            stack.append(a + 1)
        elif op == OPCODES['DEC']:
            a, = popn(1)
            stack.append(a - 1)
        elif op == OPCODES['TIME']:
            a, = popn(1)
            time.sleep(a)
        elif op == OPCODES['DUP']:
            if not stack:
                raise IndexError("Stack underflow: cannot DUP on empty stack")
            stack.append(stack[-1])
        elif op == OPCODES['DEBUG']:
            print("=== DEBUG ===")
            print("PC:", i)
            print("STACK:", stack)
            print("MEM:", mem)
            print("CALL STACK:", call_stack)
            print("=============")
        elif op == OPCODES['CALL']:
            call_stack.append(i)
            i = arg
        elif op == OPCODES['RET']:
            if not call_stack:
                raise RuntimeError("Return without CALL")
            i = call_stack.pop()
        elif op == OPCODES['CALLFN']:
            if not stack:
                raise IndexError("Stack underflow: CALLFN expects a function reference on top of stack")
            func_ref = stack.pop()
            if isinstance(func_ref, tuple) and func_ref[0] == 'FUNC' and isinstance(func_ref[1], int):
                call_stack.append(i)
                i = func_ref[1]
            elif isinstance(func_ref, int):
                call_stack.append(i)
                i = func_ref
            else:
                raise TypeError(f"Invalid function reference on top of stack: {func_ref!r}")
        else:
            raise NotImplementedError(f"Opcode {op} not implemented")
