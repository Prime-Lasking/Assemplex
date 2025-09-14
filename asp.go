package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Registers r1–r16
var registers = map[string]float64{}

// Comparison flag (like CPU flags)
var cmpFlag int = 0 // 0 = false, 1 = true

func initRegisters() {
	for i := 1; i <= 16; i++ {
		registers[fmt.Sprintf("r%d", i)] = 0
	}
}

// Get value of operand (number or register)
func getValue(operand string) (float64, error) {
	operand = strings.ToLower(operand)
	if val, ok := registers[operand]; ok {
		return val, nil
	}
	f, err := strconv.ParseFloat(operand, 64)
	if err != nil {
		return 0, fmt.Errorf("unknown operand: %s", operand)
	}
	return f, nil
}

// Execute a single instruction line (without jumps)
func run(line string) error {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "//") || strings.HasSuffix(line, ":") {
		return nil
	}

	tokens := strings.Fields(strings.ReplaceAll(line, ",", " "))
	if len(tokens) == 0 {
		return nil
	}
	instruction := strings.ToUpper(tokens[0])
	operands := []string{}
	if len(tokens) > 1 {
		for _, op := range tokens[1:] {
			operands = append(operands, strings.ToLower(strings.TrimSpace(op)))
		}
	}

	dest := ""
	if len(operands) > 0 {
		dest = operands[0]
	}
	var srcVal float64
	var err error
	if len(operands) > 1 {
		srcVal, err = getValue(operands[1])
		if err != nil {
			return err
		}
	}

	switch instruction {
	case "MOV":
		registers[dest] = srcVal
	case "ADD":
		registers[dest] += srcVal
	case "SUB":
		registers[dest] -= srcVal
	case "MUL":
		registers[dest] *= srcVal
	case "DIV":
		if srcVal == 0 {
			return fmt.Errorf("division by zero")
		}
		registers[dest] /= srcVal
	case "MOD":
		if srcVal == 0 {
			return fmt.Errorf("modulo by zero")
		}
		registers[dest] = float64(int(registers[dest]) % int(srcVal))
	case "NEG":
		registers[dest] = -registers[dest]
	case "INC":
		registers[dest] += 1
	case "DEC":
		registers[dest] -= 1
	case "PRINT":
		val, _ := getValue(dest)
		fmt.Println(val)
	case "INPUT":
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter value for %s: ", dest)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		f, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return fmt.Errorf("invalid input: %s", text)
		}
		registers[dest] = f

	// Comparison instructions: set cmpFlag
	case "LT":
		destVal, err := getValue(operands[0])
		if err != nil {
			return err
		}
		cmpFlag = 0
		if destVal < srcVal {
			cmpFlag = 1
		}
	case "LE":
		destVal, err := getValue(operands[0])
		if err != nil {
			return err
		}
		cmpFlag = 0
		if destVal <= srcVal {
			cmpFlag = 1
		}
	case "GT":
		destVal, err := getValue(operands[0])
		if err != nil {
			return err
		}
		cmpFlag = 0
		if destVal > srcVal {
			cmpFlag = 1
		}
	case "GE":
		destVal, err := getValue(operands[0])
		if err != nil {
			return err
		}
		cmpFlag = 0
		if destVal >= srcVal {
			cmpFlag = 1
		}
	case "EQ":
		destVal, err := getValue(operands[0])
		if err != nil {
			return err
		}
		cmpFlag = 0
		if destVal == srcVal {
			cmpFlag = 1
		}
	case "NE":
		destVal, err := getValue(operands[0])
		if err != nil {
			return err
		}
		cmpFlag = 0
		if destVal != srcVal {
			cmpFlag = 1
		}

	default:
		// JMP, JZ, JNZ handled in runProgram
	}

	return nil
}

// Build a label map: label -> line number
func findLabels(lines []string) map[string]int {
	labels := map[string]int{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasSuffix(line, ":") {
			label := strings.ToLower(strings.TrimSuffix(line, ":"))
			labels[label] = i
		}
	}
	return labels
}

// Run a full program (handles jumps)
func runProgram(program string) error {
	lines := strings.Split(program, "\n")
	labels := findLabels(lines)
	pc := 0

	for pc < len(lines) {
		line := strings.TrimSpace(lines[pc])
		if line == "" || strings.HasPrefix(line, "//") || strings.HasSuffix(line, ":") {
			pc++
			continue
		}

		tokens := strings.Fields(strings.ReplaceAll(line, ",", " "))
		instruction := strings.ToUpper(tokens[0])
		operands := []string{}
		if len(tokens) > 1 {
			for _, op := range tokens[1:] {
				operands = append(operands, strings.ToLower(strings.TrimSpace(op)))
			}
		}

		switch instruction {
		case "JMP":
			if lbl, ok := labels[operands[0]]; ok {
				pc = lbl
				continue
			} else {
				return fmt.Errorf("unknown label: %s", operands[0])
			}
		case "JZ":
			if cmpFlag == 0 {
				if lbl, ok := labels[operands[0]]; ok {
					pc = lbl
					continue
				} else {
					return fmt.Errorf("unknown label: %s", operands[0])
				}
			}
		case "JNZ":
			if cmpFlag != 0 {
				if lbl, ok := labels[operands[0]]; ok {
					pc = lbl
					continue
				} else {
					return fmt.Errorf("unknown label: %s", operands[0])
				}
			}
		case "HALT":
			return nil
		}

		if err := run(line); err != nil {
			return fmt.Errorf("line %d: %v", pc+1, err)
		}
		pc++
	}

	return nil
}

// Main: run a .asp file passed as command-line argument
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myassemp.exe <file.asp>")
		return
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".asp") {
		fmt.Println("Error: file must end with .asp")
		return
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	initRegisters()
	err = runProgram(string(content))
	if err != nil {
		fmt.Println("Execution error:", err)
	}
}
