package main

import (
    "bufio"
    "fmt"
    "math/big"
    "os"
    "strconv"
    "strings"
)

var (
    registers   = make(map[string]*big.Int)
    regBits     = map[string]int{}
    regCycles   = map[string]int{}
    labels      = map[string]int{}
    functions   = map[string]int{}
    cmpFlag     = false
    pc          = 0
    cycleCount  = 0
    program     []string
    callStack   []int
)

func initRegisters() {
    for i := 1; i <= 6; i++ {
        r := fmt.Sprintf("r%d", i)
        registers[r] = big.NewInt(0)
        regBits[r] = 16
        regCycles[r] = 1
    }
    for i := 7; i <= 10; i++ {
        r := fmt.Sprintf("r%d", i)
        registers[r] = big.NewInt(0)
        regBits[r] = 32
        regCycles[r] = 2
    }
    for i := 11; i <= 13; i++ {
        r := fmt.Sprintf("r%d", i)
        registers[r] = big.NewInt(0)
        regBits[r] = 64
        regCycles[r] = 4
    }
    for i := 14; i <= 16; i++ {
        r := fmt.Sprintf("r%d", i)
        registers[r] = big.NewInt(0)
        regBits[r] = 128
        regCycles[r] = 8
    }
}

func maskValue(reg string) {
    bits := regBits[reg]
    max := new(big.Int).Lsh(big.NewInt(1), uint(bits))
    max.Sub(max, big.NewInt(1))
    registers[reg].And(registers[reg], max)
}

func parseProgram(lines []string) {
    for i, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, ";") {
            continue
        }
        if strings.HasSuffix(line, ":") {
            label := strings.TrimSuffix(line, ":")
            labels[label] = len(program)
            continue
        }
        parts := strings.Fields(line)
        if len(parts) >= 2 && parts[0] == "FUNC" {
            functions[parts[1]] = len(program)
        }
        program = append(program, line)
        if line == "ENDFUNC" {
            program = append(program, "RET") // auto-ret for funcs
        }
    }
}

func getVal(s string) *big.Int {
    if v, ok := registers[s]; ok {
        return new(big.Int).Set(v)
    }
    n, _ := strconv.ParseInt(s, 10, 64)
    return big.NewInt(n)
}

func setReg(r string, v *big.Int) {
    registers[r].Set(v)
    maskValue(r)
    cycleCount += regCycles[r]
}

func execute() {
    for pc < len(program) {
        line := program[pc]
        pc++
        parts := strings.Fields(line)
        if len(parts) == 0 {
            continue
        }

        op := strings.ToUpper(parts[0])
        switch op {
        case "MOV":
            dst, src := parts[1], parts[2]
            setReg(dst, getVal(src))
        case "ADD":
            dst, src := parts[1], parts[2]
            registers[dst].Add(registers[dst], getVal(src))
            maskValue(dst); cycleCount += regCycles[dst]
        case "SUB":
            dst, src := parts[1], parts[2]
            registers[dst].Sub(registers[dst], getVal(src))
            maskValue(dst); cycleCount += regCycles[dst]
        case "MUL":
            dst, src := parts[1], parts[2]
            registers[dst].Mul(registers[dst], getVal(src))
            maskValue(dst); cycleCount += regCycles[dst]
        case "DIV":
            dst, src := parts[1], parts[2]
            registers[dst].Div(registers[dst], getVal(src))
            maskValue(dst); cycleCount += regCycles[dst]
        case "MOD":
            dst, src := parts[1], parts[2]
            registers[dst].Mod(registers[dst], getVal(src))
            maskValue(dst); cycleCount += regCycles[dst]
        case "NEG":
            dst := parts[1]
            registers[dst].Neg(registers[dst])
            maskValue(dst); cycleCount += regCycles[dst]
        case "INC":
            dst := parts[1]
            registers[dst].Add(registers[dst], big.NewInt(1))
            maskValue(dst); cycleCount += regCycles[dst]
        case "DEC":
            dst := parts[1]
            registers[dst].Sub(registers[dst], big.NewInt(1))
            maskValue(dst); cycleCount += regCycles[dst]
        case "PRINT":
            fmt.Println(registers[parts[1]])
        case "INPUT":
            reader := bufio.NewReader(os.Stdin)
            fmt.Print("? ")
            text, _ := reader.ReadString('\n')
            text = strings.TrimSpace(text)
            n, _ := strconv.ParseInt(text, 10, 64)
            setReg(parts[1], big.NewInt(n))
        case "LT", "LE", "GT", "GE", "EQ", "NE":
            a, b := parts[1], parts[2]
            va, vb := getVal(a), getVal(b)
            switch op {
            case "LT": cmpFlag = va.Cmp(vb) < 0
            case "LE": cmpFlag = va.Cmp(vb) <= 0
            case "GT": cmpFlag = va.Cmp(vb) > 0
            case "GE": cmpFlag = va.Cmp(vb) >= 0
            case "EQ": cmpFlag = va.Cmp(vb) == 0
            case "NE": cmpFlag = va.Cmp(vb) != 0
            }
        case "JMP":
            pc = labels[parts[1]]
        case "JZ":
            if cmpFlag { pc = labels[parts[1]] }
        case "JNZ":
            if !cmpFlag { pc = labels[parts[1]] }
        case "CALL":
            callStack = append(callStack, pc)
            pc = functions[parts[1]]
        case "RET":
            if len(callStack) == 0 {
                return
            }
            pc = callStack[len(callStack)-1]
            callStack = callStack[:len(callStack)-1]
        case "HALT":
            return
        }
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: asp <file.asp>")
        return
    }
    initRegisters()
    data, _ := os.ReadFile(os.Args[1])
    lines := strings.Split(string(data), "\n")
    parseProgram(lines)
    execute()
    fmt.Printf("[Program finished in %d cycles]\n", cycleCount)
}
