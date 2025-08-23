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
}

def assemble(asm):
    code = []
    vars = {}
    labels = {}
    next_slot = 0
    lines = []

    # First pass: strip comments and store clean instructions
    for line in asm.strip().split('\n'):
        line = line.split(';')[0].strip()
        if not line:
            continue
        lines.append(line)

    # First pass: collect label positions
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

    # Second pass: generate code
    for line in lines:
        if line.endswith(':'):
            continue  # Skip labels

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
                code.append(labels[val])  # Replace label with address
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

        if op == 1:  # PUSH
            stack.append(code[i])
            i += 1

        elif op == 2:  # ADD
            b = stack.pop()
            a = stack.pop()
            if isinstance(a, str) or isinstance(b, str):
                stack.append(str(a) + str(b))
            else:
                stack.append(a + b)

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
            if b == 0:
                raise ZeroDivisionError("Division by zero")
            stack.append(a / b)

        elif op == 10:  # EXP
            b = stack.pop()
            a = stack.pop()
            stack.append(a ** b)

        elif op == 11:  # SQRT
            a = stack.pop()
            if a < 0:
                raise ValueError("Cannot take square root of negative number")
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
            try:
                stack.append(int(val))
            except ValueError:
                raise ValueError(f"Cannot convert {val} to int")

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

        elif op == 21:  # READ
            var_slot = code[i]
            i += 1
            user_input = input("Enter value: ")
            mem[var_slot] = user_input

program = """
; Countdown from input number to 0
PUSH 100
STORE x
LOAD x
TOINT
STORE counter

loop_start:
LOAD counter
PUSH 0
GT          ; while counter > 0
JZ end_loop

LOAD counter
PRINT

LOAD counter
PUSH 1
SUB
STORE counter

JMP loop_start

end_loop:
HALT

"""

if __name__ == "__main__":
    code, slots = assemble(program)
    run_vm(code, slots)
