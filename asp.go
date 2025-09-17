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
	regMask    = map[string]*big.Int{}
	labels     = make(map[string]int)
	functions  = make(map[string]int)
	cmpFlag    = false
	pc         = 0
	cycleCount = 0
	program    []Instruction
	callStack  []int
)

type Instruction struct {
	op   string
	args []string
}

// --- Register Initialization ---
func initRegisters() {
	// 16-bit
	for i := 1; i <= 6; i++ {
		r := fmt.Sprintf("r%d", i)
		registers[r] = big.NewInt(0)
		regBits[r] = 16
		regCycles[r] = 1
		regMask[r] = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 16), big.NewInt(1))
	}
	// 32-bit
	for i := 7; i <= 10; i++ {
		r := fmt.Sprintf("r%d", i)
		registers[r] = big.NewInt(0)
		regBits[r] = 32
		regCycles[r] = 2
		regMask[r] = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 32), big.NewInt(1))
	}
	// 64-bit
	for i := 11; i <= 13; i++ {
		r := fmt.Sprintf("r%d", i)
		registers[r] = big.NewInt(0)
		regBits[r] = 64
		regCycles[r] = 4
		regMask[r] = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 64), big.NewInt(1))
	}
	// 128-bit
	for i := 14; i <= 16; i++ {
		r := fmt.Sprintf("r%d", i)
		registers[r] = big.NewInt(0)
		regBits[r] = 128
		regCycles[r] = 8
		regMask[r] = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))
	}
}

func maskValue(reg string) {
	registers[reg].And(registers[reg], regMask[reg])
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
		program = append(program, Instruction{op: parts[0], args: parts[1:]})
		if parts[0] == "ENDFUNC" {
			program = append(program, Instruction{op: "RET"})
		}
	}
}

// --- Value Helpers ---
func getVal(s string) *big.Int {
	if v, ok := registers[s]; ok {
		return v
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Printf("Invalid literal: %s\n", s)
		os.Exit(1)
	}
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
		ins := program[pc]
		pc++
		args := ins.args
		switch strings.ToUpper(ins.op) {
		case "MOV":
			setReg(args[0], getVal(args[1]))
		case "ADD":
			registers[args[0]].Add(registers[args[0]], getVal(args[1]))
			maskValue(args[0])
			cycleCount += regCycles[args[0]]
		case "SUB":
			registers[args[0]].Sub(registers[args[0]], getVal(args[1]))
			maskValue(args[0])
			cycleCount += regCycles[args[0]]
		case "MUL":
			registers[args[0]].Mul(registers[args[0]], getVal(args[1]))
			maskValue(args[0])
			cycleCount += regCycles[args[0]]
		case "DIV":
			if getVal(args[1]).Sign() == 0 {
				fmt.Println("Runtime error: divide by zero")
				return
			}
			registers[args[0]].Div(registers[args[0]], getVal(args[1]))
			maskValue(args[0])
			cycleCount += regCycles[args[0]]
		case "MOD":
			if getVal(args[1]).Sign() == 0 {
				fmt.Println("Runtime error: modulo by zero")
				return
			}
			registers[args[0]].Mod(registers[args[0]], getVal(args[1]))
			maskValue(args[0])
			cycleCount += regCycles[args[0]]
		case "NEG":
			registers[args[0]].Neg(registers[args[0]])
			maskValue(args[0])
			cycleCount += regCycles[args[0]]
		case "INC":
			registers[args[0]].Add(registers[args[0]], big.NewInt(1))
			maskValue(args[0])
			cycleCount += regCycles[args[0]]
		case "DEC":
			registers[args[0]].Sub(registers[args[0]], big.NewInt(1))
			maskValue(args[0])
			cycleCount += regCycles[args[0]]
		case "PRINT":
			fmt.Println(registers[args[0]])
		case "INPUT":
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("? ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			n, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				fmt.Println("Invalid input, expected integer")
				pc--
				continue
			}
			setReg(args[0], big.NewInt(n))
		case "LT", "LE", "GT", "GE", "EQ", "NE":
			va, vb := getVal(args[0]), getVal(args[1])
			switch strings.ToUpper(ins.op) {
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
			pc = labels[args[0]]
		case "JZ":
			if cmpFlag {
				pc = labels[args[0]]
			}
		case "JNZ":
			if !cmpFlag {
				pc = labels[args[0]]
			}
		case "CALL":
			callStack = append(callStack, pc)
			pc = functions[args[0]]
		case "RET":
			if len(callStack) == 0 {
				return
			}
			pc = callStack[len(callStack)-1]
			callStack = callStack[:len(callStack)-1]
		case "HALT":
			return
		default:
			fmt.Printf("Unknown instruction at line %d: %s\n", pc, ins.op)
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
	fmt.Println("  DIV rX val       Divide (runtime error if val=0)")
	fmt.Println("  MOD rX val       Modulo (runtime error if val=0)")
	fmt.Println("  NEG rX           Negate register")
	fmt.Println("  INC rX           Increment register by 1")
	fmt.Println("  DEC rX           Decrement register by 1")
	fmt.Println("  PRINT rX         Print register value")
	fmt.Println("  INPUT rX         Input integer into register")
	fmt.Println("  LT rX val        Set flag if rX < val")
	fmt.Println("  LE rX val        Set flag if rX <= val")
	fmt.Println("  GT rX val        Set flag if rX > val")
	fmt.Println("  GE rX val        Set flag if rX >= val")
	fmt.Println("  EQ rX val        Set flag if rX == val")
	fmt.Println("  NE rX val        Set flag if rX != val")
	fmt.Println("  JMP label        Jump unconditionally")
	fmt.Println("  JZ label         Jump if last comparison was true")
	fmt.Println("  JNZ label        Jump if last comparison was false")
	fmt.Println("  CALL func        Call function")
	fmt.Println("  RET              Return from function")
	fmt.Println("  FUNC name        Function start")
	fmt.Println("  ENDFUNC          Function end (implicit RET added)")
	fmt.Println("  HALT             Stop program\n")

	fmt.Println("Labels:   <name>:")
	fmt.Println("Comments: ; text after semicolon")
}

// --- Main ---
func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--version":
			fmt.Println(version)
			return
		case "--help":
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
