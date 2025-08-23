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
    'READ': 21,
    'MOD': 22,
    'AND': 23,
    'OR': 24,
    'NOT': 25,
    'IF': 26,
    'ELSE': 27,
}

def assemble(asm):
    code = []
    vars = {}
    labels = {}
    next_slot = 0
    lines = []

    for line in asm.strip().split('\n'):
        line = line.split(';')[0].strip()
        if not line:
            continue
        lines.append(line)

    i = 0
    for line in lines:
        if line.endswith(':'):
            label = line[:-1]
            labels[label] = i
        else:
            instr, *arg = line.split(maxsplit=1)
            instr = instr.upper()
            i += 1
            if arg:
                i += 1

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
            else:
                if instr in {'STORE', 'LOAD', 'READ'}:
                    if val not in vars:
                        vars[val] = next_slot
                        next_slot += 1
                    code.append(vars[val])
                else:
                    code.append(val)
    return code, next_slot

def run_vm(code, slots):
    stack = []
    mem = [0] * slots
    i = 0

    while i < len(code):
        op = code[i]
        i += 1

        if op == 1:
            stack.append(code[i])
            i += 1
        elif op == 2:
            b = stack.pop()
            a = stack.pop()
            stack.append(a + b)
        elif op == 3:
            mem[code[i]] = stack.pop()
            i += 1
        elif op == 4:
            stack.append(mem[code[i]])
            i += 1
        elif op == 5:
            print(stack.pop())
        elif op == 6:
            break
        elif op == 7:
            b = stack.pop()
            a = stack.pop()
            stack.append(a - b)
        elif op == 8:
            b = stack.pop()
            a = stack.pop()
            stack.append(a * b)
        elif op == 9:
            b = stack.pop()
            a = stack.pop()
            stack.append(a / b)
        elif op == 10:
            b = stack.pop()
            a = stack.pop()
            stack.append(a ** b)
        elif op == 11:
            a = stack.pop()
            stack.append(a ** 0.5)
        elif op == 12:
            i = code[i]
        elif op == 13:
            addr = code[i]
            i += 1
            if stack.pop() == 0:
                i = addr
        elif op == 14:
            addr = code[i]
            i += 1
            if stack.pop() != 0:
                i = addr
        elif op == 15:
            val = stack.pop()
            stack.append(int(val))
        elif op == 16:
            val = stack.pop()
            stack.append(str(val))
        elif op == 17:
            val = stack.pop()
            stack.append(len(str(val)))
        elif op == 18:
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a == b else 0)
        elif op == 19:
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a > b else 0)
        elif op == 20:
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a < b else 0)
        elif op == 21:
            var_slot = code[i]
            i += 1
            mem[var_slot] = input("Enter value: ")
        elif op == 22:
            b = stack.pop()
            a = stack.pop()
            stack.append(a % b)
        elif op == 23:
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a and b else 0)
        elif op == 24:
            b = stack.pop()
            a = stack.pop()
            stack.append(1 if a or b else 0)
        elif op == 25:
            a = stack.pop()
            stack.append(0 if a else 1)
        elif op == 26:
            addr = code[i]
            i += 1
            cond = stack.pop()
            if not cond:
                i = addr
        elif op == 27:
            addr = code[i]
            i += 1
            i = addr

program = """
PUSH 10
PUSH 5
EQ
IF Else_Block
PUSH "10 is equal to 5"
PRINT
JMP END
Else_Block:
PUSH "10 is not equal to 5"
PRINT
JMP END
END:
HALT
"""

if __name__ == "__main__":
    code, slots = assemble(program)
    run_vm(code, slots)
