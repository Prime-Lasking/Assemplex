package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const version = "Assemplex v1.2"

var (
	registers  = make(map[string]*big.Int)
	regBits    = map[string]int{}
	regCycles  = map[string]int{}
	labels     = map[string]int{}
	functions  = map[string]int{}
	cmpFlag    = false
	pc         = 0
	cycleCount = 0
	program    []string
	callStack  []int
)

// --- Register Initialization ---
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

// --- Parser ---
func parseProgram(lines []string) {
	for _, line := range lines {
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
			program = append(program, "RET") // implicit return
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

// --- Executor ---
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
			setReg(parts[1], getVal(parts[2]))
		case "ADD":
			registers[parts[1]].Add(registers[parts[1]], getVal(parts[2]))
			maskValue(parts[1])
			cycleCount += regCycles[parts[1]]
		case "SUB":
			registers[parts[1]].Sub(registers[parts[1]], getVal(parts[2]))
			maskValue(parts[1])
			cycleCount += regCycles[parts[1]]
		case "MUL":
			registers[parts[1]].Mul(registers[parts[1]], getVal(parts[2]))
			maskValue(parts[1])
			cycleCount += regCycles[parts[1]]
		case "DIV":
			registers[parts[1]].Div(registers[parts[1]], getVal(parts[2]))
			maskValue(parts[1])
			cycleCount += regCycles[parts[1]]
		case "MOD":
			registers[parts[1]].Mod(registers[parts[1]], getVal(parts[2]))
			maskValue(parts[1])
			cycleCount += regCycles[parts[1]]
		case "NEG":
			registers[parts[1]].Neg(registers[parts[1]])
			maskValue(parts[1])
			cycleCount += regCycles[parts[1]]
		case "INC":
			registers[parts[1]].Add(registers[parts[1]], big.NewInt(1))
			maskValue(parts[1])
			cycleCount += regCycles[parts[1]]
		case "DEC":
			registers[parts[1]].Sub(registers[parts[1]], big.NewInt(1))
			maskValue(parts[1])
			cycleCount += regCycles[parts[1]]
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
			va, vb := getVal(parts[1]), getVal(parts[2])
			switch op {
			case "LT":
				cmpFlag = va.Cmp(vb) < 0
			case "LE":
				cmpFlag = va.Cmp(vb) <= 0
			case "GT":
				cmpFlag = va.Cmp(vb) > 0
			case "GE":
				cmpFlag = va.Cmp(vb) >= 0
			case "EQ":
				cmpFlag = va.Cmp(vb) == 0
			case "NE":
				cmpFlag = va.Cmp(vb) != 0
			}
		case "JMP":
			pc = labels[parts[1]]
		case "JZ":
			if cmpFlag {
				pc = labels[parts[1]]
			}
		case "JNZ":
			if !cmpFlag {
				pc = labels[parts[1]]
			}
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

// --- Help Screen ---
func showHelp() {
	fmt.Println(version)
	fmt.Println("\nUsage:")
	fmt.Println("  asp <file.asp>       Run program")
	fmt.Println("  asp --version        Show version")
	fmt.Println("  asp --help           Show this help\n")

	fmt.Println("Registers:")
	fmt.Println("  r1–r6   = 16-bit (fastest)")
	fmt.Println("  r7–r10  = 32-bit")
	fmt.Println("  r11–r13 = 64-bit")
	fmt.Println("  r14–r16 = 128-bit (slowest)\n")

	fmt.Println("Instructions:")
	fmt.Println("  MOV rX val       Move value into register")
	fmt.Println("  ADD rX val       Add")
	fmt.Println("  SUB rX val       Subtract")
	fmt.Println("  MUL rX val       Multiply")
	fmt.Println("  DIV rX val       Divide")
	fmt.Println("  MOD rX val       Modulo")
	fmt.Println("  NEG rX           Negate")
	fmt.Println("  INC rX           Increment")
	fmt.Println("  DEC rX           Decrement")
	fmt.Println("  PRINT rX         Print register")
	fmt.Println("  INPUT rX         Input integer into register")
	fmt.Println("  LT/LE/GT/GE/EQ/NE rX val   Compare (sets flag)")
	fmt.Println("  JMP label        Jump to label")
	fmt.Println("  JZ label         Jump if last compare was true")
	fmt.Println("  JNZ label        Jump if last compare was false")
	fmt.Println("  CALL func        Call function")
	fmt.Println("  RET              Return from function")
	fmt.Println("  FUNC name        Function start")
	fmt.Println("  ENDFUNC          Function end")
	fmt.Println("  HALT             Stop program\n")

	fmt.Println("Labels:   <name>:")
	fmt.Println("Comments: ; text after semicolon")
}

// --- Main ---
func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "--version" {
			fmt.Println(version)
			return
		}
		if os.Args[1] == "--help" {
			showHelp()
			return
		}
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: asp <program.asp>   (or asp --help)")
		return
	}

	initRegisters()

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	parseProgram(lines)
	execute()
	fmt.Printf("[Program finished in %d cycles]\n", cycleCount)
}
