package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

// --- Register structure ---
type Register struct {
	Name  string
	Bits  int // 16,32,64,128
	Value *big.Int
	Speed int // operation cost
}

var registers = map[string]*Register{}
var cmpFlag int = 0
var functions = map[string]int{}
var globalLines []string
var cycleCount int64 = 0

// --- Mask to simulate overflow ---
func maskValue(r *Register) {
	max := new(big.Int).Lsh(big.NewInt(1), uint(r.Bits))
	max.Sub(max, big.NewInt(1))
	r.Value.And(r.Value, max)
}

// --- Initialize heterogeneous registers ---
func initRegisters() {
	for _, r := range []struct {
		name  string
		bits  int
		speed int
	}{
		{"r1", 16, 1}, {"r2", 16, 1}, {"r3", 16, 1}, {"r4", 16, 1}, {"r5", 16, 1}, {"r6", 16, 1},
		{"r7", 32, 2}, {"r8", 32, 2}, {"r9", 32, 2}, {"r10", 32, 2},
		{"r11", 64, 4}, {"r12", 64, 4}, {"r13", 64, 4},
		{"r14", 128, 8}, {"r15", 128, 8}, {"r16", 128, 8},
	} {
		registers[r.name] = &Register{
			Name:  r.name,
			Bits:  r.bits,
			Value: big.NewInt(0),
			Speed: r.speed,
		}
	}
}

// --- Operand handling ---
func getValue(operand string) (*big.Int, error) {
	operand = strings.ToLower(operand)
	if val, ok := registers[operand]; ok {
		return val.Value, nil
	}
	bi, ok := new(big.Int).SetString(operand, 10)
	if !ok {
		return nil, fmt.Errorf("unknown operand: %s", operand)
	}
	return bi, nil
}

// --- Arithmetic helpers ---
func add(dest, src string) {
	registers[dest].Value.Add(registers[dest].Value, registers[src].Value)
	maskValue(registers[dest])
	cycleCount += int64(registers[dest].Speed)
}
func sub(dest, src string) {
	registers[dest].Value.Sub(registers[dest].Value, registers[src].Value)
	maskValue(registers[dest])
	cycleCount += int64(registers[dest].Speed)
}
func mul(dest, src string) {
	registers[dest].Value.Mul(registers[dest].Value, registers[src].Value)
	maskValue(registers[dest])
	cycleCount += int64(registers[dest].Speed)
}
func div(dest, src string) error {
	if registers[src].Value.Sign() == 0 {
		return fmt.Errorf("division by zero")
	}
	registers[dest].Value.Div(registers[dest].Value, registers[src].Value)
	maskValue(registers[dest])
	cycleCount += int64(registers[dest].Speed)
	return nil
}
func mod(dest, src string) error {
	if registers[src].Value.Sign() == 0 {
		return fmt.Errorf("modulo by zero")
	}
	registers[dest].Value.Mod(registers[dest].Value, registers[src].Value)
	maskValue(registers[dest])
	cycleCount += int64(registers[dest].Speed)
	return nil
}
func neg(dest string) {
	registers[dest].Value.Neg(registers[dest].Value)
	maskValue(registers[dest])
	cycleCount += int64(registers[dest].Speed)
}
func inc(dest string) {
	registers[dest].Value.Add(registers[dest].Value, big.NewInt(1))
	maskValue(registers[dest])
	cycleCount += int64(registers[dest].Speed)
}
func dec(dest string) {
	registers[dest].Value.Sub(registers[dest].Value, big.NewInt(1))
	maskValue(registers[dest])
	cycleCount += int64(registers[dest].Speed)
}

// --- Run single line ---
func run(line string) error {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "//") || strings.HasSuffix(line, ":") {
		return nil
	}
	tokens := strings.Fields(strings.ReplaceAll(line, ",", " "))
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
	var srcVal *big.Int
	var err error
	if len(operands) > 1 {
		srcVal, err = getValue(operands[1])
		if err != nil {
			return err
		}
	}

	switch instruction {
	case "MOV":
		registers[dest].Value.Set(srcVal)
		maskValue(registers[dest])
		cycleCount += int64(registers[dest].Speed)
	case "ADD":
		add(dest, operands[1])
	case "SUB":
		sub(dest, operands[1])
	case "MUL":
		mul(dest, operands[1])
	case "DIV":
		if err := div(dest, operands[1]); err != nil {
			return err
		}
	case "MOD":
		if err := mod(dest, operands[1]); err != nil {
			return err
		}
	case "NEG":
		neg(dest)
	case "INC":
		inc(dest)
	case "DEC":
		dec(dest)
	case "PRINT":
		fmt.Println(registers[dest].Value.String())
	case "INPUT":
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter value for %s: ", dest)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		val, ok := new(big.Int).SetString(text, 10)
		if !ok {
			return fmt.Errorf("invalid input: %s", text)
		}
		registers[dest].Value.Set(val)
		maskValue(registers[dest])
		cycleCount += int64(registers[dest].Speed)
	case "LT", "LE", "GT", "GE", "EQ", "NE":
		val1, err := getValue(operands[0])
		if err != nil {
			return err
		}
		val2 := srcVal
		cmpFlag = 0
		switch instruction {
		case "LT":
			if val1.Cmp(val2) < 0 {
				cmpFlag = 1
			}
		case "LE":
			if val1.Cmp(val2) <= 0 {
				cmpFlag = 1
			}
		case "GT":
			if val1.Cmp(val2) > 0 {
				cmpFlag = 1
			}
		case "GE":
			if val1.Cmp(val2) >= 0 {
				cmpFlag = 1
			}
		case "EQ":
			if val1.Cmp(val2) == 0 {
				cmpFlag = 1
			}
		case "NE":
			if val1.Cmp(val2) != 0 {
				cmpFlag = 1
			}
		}
	}
	return nil
}

// --- Find labels and functions ---
func findLabelsAndFunctions(lines []string) (map[string]int, map[string]int) {
	labels := map[string]int{}
	funcs := map[string]int{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasSuffix(line, ":") {
			label := strings.ToLower(strings.TrimSuffix(line, ":"))
			if strings.HasPrefix(label, "func ") {
				funcName := strings.TrimSpace(strings.TrimPrefix(label, "func "))
				funcs[funcName] = i + 1
			} else {
				labels[label] = i
			}
		}
	}
	return labels, funcs
}

// --- Run program from a line ---
func runProgramFrom(startLine int) error {
	lines := globalLines
	labels, _ := findLabelsAndFunctions(lines)
	pc := startLine
	for pc < len(lines) {
		line := strings.TrimSpace(lines[pc])
		if line == "" || strings.HasPrefix(line, "//") {
			pc++
			continue
		}
		if strings.HasPrefix(strings.ToUpper(line), "ENDFUNC") {
			return nil
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
		case "CALL":
			funcName := operands[0]
			if lineNum, ok := functions[funcName]; ok {
				if err := runProgramFrom(lineNum); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("unknown function: %s", funcName)
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

// --- Handle INCLUDE, IMPORT, IMPORTVAR ---
func preprocessIncludes(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %v", filename, err)
	}
	lines := strings.Split(string(content), "\n")
	expanded := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(line), "INCLUDE") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				return nil, fmt.Errorf("invalid INCLUDE: %s", line)
			}
			includeFile := strings.Trim(parts[1], `"`)
			incLines, err := preprocessIncludes(includeFile)
			if err != nil {
				return nil, err
			}
			expanded = append(expanded, incLines...)
		} else if strings.HasPrefix(strings.ToUpper(line), "IMPORT") {
			parts := strings.Fields(line)
			if len(parts) != 4 || strings.ToUpper(parts[2]) != "FROM" {
				return nil, fmt.Errorf("invalid IMPORT: %s", line)
			}
			funcName := parts[1]
			importFile := strings.Trim(parts[3], `"`)
			incLines, err := preprocessIncludes(importFile)
			if err != nil {
				return nil, err
			}
			funcStart := -1
			funcEnd := -1
			for i, l := range incLines {
				if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(l)), "FUNC "+strings.ToUpper(funcName)) {
					funcStart = i
				}
				if funcStart >= 0 && strings.HasPrefix(strings.ToUpper(strings.TrimSpace(l)), "ENDFUNC") {
					funcEnd = i
					break
				}
			}
			if funcStart >= 0 && funcEnd >= 0 {
				expanded = append(expanded, incLines[funcStart:funcEnd+1]...)
			}
		} else if strings.HasPrefix(strings.ToUpper(line), "IMPORTVAR") {
			parts := strings.Fields(line)
			if len(parts) != 4 || strings.ToUpper(parts[2]) != "FROM" {
				return nil, fmt.Errorf("invalid IMPORTVAR: %s", line)
			}
			varName := parts[1]
			importFile := strings.Trim(parts[3], `"`)
			incLines, err := preprocessIncludes(importFile)
			if err != nil {
				return nil, err
			}
			for _, l := range incLines {
				l = strings.TrimSpace(l)
				if strings.HasPrefix(strings.ToUpper(l), "MOV") {
					tokens := strings.Fields(strings.ReplaceAll(l, ",", " "))
					dest := strings.ToLower(tokens[1])
					if dest == strings.ToLower(varName) {
						expanded = append(expanded, l)
						break
					}
				}
			}
		} else {
			expanded = append(expanded, line)
		}
	}
	return expanded, nil
}

// --- Run full program with includes ---
func runProgramWithIncludes(filename string) error {
	lines, err := preprocessIncludes(filename)
	if err != nil {
		return err
	}
	globalLines = lines
	_, functions = findLabelsAndFunctions(globalLines)
	return runProgramFrom(0)
}

// --- Main ---
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: assemplex.exe <file.asp>")
		return
	}
	filename := os.Args[1]
	initRegisters()
	err := runProgramWithIncludes(filename)
	if err != nil {
		fmt.Println("Execution error:", err)
	}
	fmt.Println("Total cycles:", cycleCount)
}
