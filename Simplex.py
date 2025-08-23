OPCODES = {
    'PUSH': 1,
    'ADD': 2,
    'STORE': 3,
    'LOAD': 4,
    'PRINT': 5,
    'HALT': 6,
    'SUB': 7,
    'MUL': 8,
    'DIV': 9,
    'EXP': 10,
    'SQRT': 11,
    'JMP': 12,
    'JZ': 13,
    'JNZ': 14,
    'TOINT': 15,
    'TOSTR': 16,
    'LEN': 17,
    'EQ': 18,
    'GT': 19,
    'LT': 20,
    'IN': 21,
    'MOD': 22,
    'AND': 23,
    'OR': 24,
    'NOT': 25,
    'IF': 26,
    'ELSE': 27,
    'INC': 28,
    'DEC': 29,
}


def assemble(asm):
    code = []
    vars = {}
    labels = {}
    next_slot = 0
    lines = []

    # Preprocess: remove comments and blank lines
    for line in asm.strip().split('\n'):
        line = line.split(';')[0].strip()
        if line:
            lines.append(line)

    # First pass: determine label positions
    pc = 0
    for line in lines:
        if line.endswith(':'):
            label = line[:-1]
            labels[label] = pc
        else:
            instr, *arg = line.split(maxsplit=1)
            instr = instr.upper()
            pc += 1  # opcode
            if arg:
                pc += 1  # argument

    # Second pass: generate code
    for line in lines:
        if line.endswith(':'):
            continue
        instr, *arg = line.split(maxsplit=1)
        instr = instr.upper()
        code.append(OPCODES[instr])
        if arg:
            val = arg[0].strip()
            if val.lstrip('-').isdigit():
                code.append(int(val))
            elif val.startswith('"') and val.endswith('"'):
                code.append(val[1:-1])
            elif val in labels:
                code.append(labels[val])
            elif instr in {'STORE', 'LOAD', 'IN', 'INC', 'DEC'}:
                if val not in vars:
                    vars[val] = next_slot
                    next_slot += 1
                code.append(vars[val])
            else:
                code.append(val)  # Possibly a runtime error
    return code, next_slot


def run_vm(code, slots):
    stack = []
    mem = [0] * slots
    i = 0

    while i < len(code):
        op = code[i]
        i += 1

        if op == 1:  # PUSH
            stack.append(code[i])
            i += 1
        elif op == 2:  # ADD
            b = stack.pop()
            a = stack.pop()
            if isinstance(a, str) and isinstance(b, str):
              stack.append(str(a) + str(b))
            else:
                stack.append(a+b)
        elif op == 3:  # STORE
            mem[code[i]] = stack.pop()
            i += 1
        elif op == 4:  # LOAD
            stack.append(mem[code[i]])
            i += 1
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
            i = code[i]
        elif op == 13:  # JZ
            addr = code[i]
            i += 1
            if stack.pop() == 0:
                i = addr
        elif op == 14:  # JNZ
            addr = code[i]
            i += 1
            if stack.pop() != 0:
                i = addr
        elif op == 15:  # TOINT
            val = stack.pop()
            stack.append(int(val))
        elif op == 16:  # TOSTR
            val = stack.pop()
            stack.append(str(val))
        elif op == 17:  # LEN
            val = stack.pop()
            stack.append(len(str(val)))
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
            var_slot = code[i]
            i += 1
            mem[var_slot] = input("Enter value: ")
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
            addr = code[i]
            i += 1
            cond = stack.pop()
            if not cond:
                i = addr
        elif op == 27:  # ELSE
            addr = code[i]
            i += 1
            i = addr
        elif op == 28:  # INC
            a = stack.pop()
            if isinstance(a, str):
                raise TypeError("Cannot use str with INC")
            else:
                stack.append(a + 1)
        elif op == 29:  # DEC
            a = stack.pop()
            if isinstance(a, str):
                raise TypeError("Cannot use str with DEC")
            else:
                stack.append(a - 1)

program = """

"""

if __name__ == "__main__":
    code, slots = assemble(program)
    run_vm(code, slots)
