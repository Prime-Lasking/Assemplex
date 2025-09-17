package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// --- Version ---
const version = "Assemplex v2.0"

// --- Variable Definition ---
type Var struct {
	Type  string
	Value interface{}
}

var (
	globalMemory = make(map[string]*Var)
	scopeStack   []map[string]*Var // for function local scopes
	labels       = map[string]int{}
	functions    = map[string]int{}
	cmpFlag      = false
	pc           = 0
	cycleCount   = 0
	program      []string
	callStack    []int
	returnValue  *Var
)

// --- Helpers ---
func maskInt(v int64, typ string) int64 {
	switch typ {
	case "INT16":
		return v & 0xFFFF
	case "INT32":
		return v & 0xFFFFFFFF
	default:
		return v
	}
}

func defaultValue(typ string) interface{} {
	switch typ {
	case "INT16", "INT32", "INT64":
		return int64(0)
	case "FLOAT16", "FLOAT32", "FLOAT64":
		return float64(0)
	case "CHAR":
		return ""
	default:
		return nil
	}
}

func getMemory() map[string]*Var {
	if len(scopeStack) == 0 {
		return globalMemory
	}
	return scopeStack[len(scopeStack)-1]
}

func getVarVal(s string) interface{} {
	mem := getMemory()
	if v, ok := mem[s]; ok {
		return v.Value
	}
	if v, ok := globalMemory[s]; ok {
		return v.Value
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
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
		if len(parts) >= 2 && strings.ToUpper(parts[0]) == "FUNC" {
			functions[parts[1]] = len(program)
		}
		program = append(program, line)
		if strings.ToUpper(line) == "ENDFUNC" {
			program = append(program, "RET") // implicit return
		}
	}
}

// --- Comparison Helper ---
func compare(a, b interface{}) int {
	switch va := a.(type) {
	case int64:
		vb := b.(int64)
		switch {
		case va < vb:
			return -1
		case va > vb:
			return 1
		default:
			return 0
		}
	case float64:
		vb := b.(float64)
		switch {
		case va < vb:
			return -1
		case va > vb:
			return 1
		default:
			return 0
		}
	case string:
		vb := b.(string)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	}
	return 0
}

// --- Execute Program ---
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
		case "VAR":
			handleVAR(parts)
		case "FREE":
			handleFREE(parts)
		case "ADD", "SUB", "MUL", "DIV":
			handleArithmetic(parts, op)
		case "PRINT":
			handlePRINT(parts)
		case "INPUT":
			handleINPUT(parts)
		case "LT", "LE", "GT", "GE", "EQ", "NE":
			handleCompare(parts, op)
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
			callFunction(parts[1], parts[2:])
		case "RET", "HALT":
			return
		}
	}
}

// --- Line Handlers ---
func handleVAR(parts []string) {
	if len(parts) < 3 {
		fmt.Println("Error: VAR requires type and name")
		return
	}
	typ := strings.ToUpper(parts[1])
	name := parts[2]
	var val interface{} = defaultValue(typ)
	if len(parts) > 3 {
		// function call as value
		if strings.Contains(parts[3], "(") {
			funcName := strings.Split(parts[3], "(")[0]
			argStr := strings.TrimRight(parts[3][len(funcName):], ")")
			args := strings.Fields(argStr)
			val = callFunction(funcName, args)
		} else {
			lit := parts[3]
			switch typ {
			case "INT16", "INT32", "INT64":
				n, _ := strconv.ParseInt(lit, 10, 64)
				val = maskInt(n, typ)
			case "FLOAT16", "FLOAT32", "FLOAT64":
				f, _ := strconv.ParseFloat(lit, 64)
				val = f
			case "CHAR":
				val = lit
			}
		}
	}
	getMemory()[name] = &Var{Type: typ, Value: val}
}

func handleFREE(parts []string) {
	delete(getMemory(), parts[1])
}

func handleArithmetic(parts []string, op string) {
	mem := getMemory()
	dst := mem[parts[1]]
	srcVal := getVarVal(parts[2])
	switch dst.Type {
	case "INT16", "INT32", "INT64":
		n := dst.Value.(int64)
		s := srcVal.(int64)
		switch op {
		case "ADD":
			dst.Value = maskInt(n+s, dst.Type)
		case "SUB":
			dst.Value = maskInt(n-s, dst.Type)
		case "MUL":
			dst.Value = maskInt(n*s, dst.Type)
		case "DIV":
			if s == 0 {
				fmt.Println("Error: division by zero")
				return
			}
			dst.Value = maskInt(n/s, dst.Type)
		}
	case "FLOAT16", "FLOAT32", "FLOAT64":
		n := dst.Value.(float64)
		s := srcVal.(float64)
		switch op {
		case "ADD":
			dst.Value = n + s
		case "SUB":
			dst.Value = n - s
		case "MUL":
			dst.Value = n * s
		case "DIV":
			dst.Value = n / s
		}
	}
}

func handlePRINT(parts []string) {
	mem := getMemory()
	if v, ok := mem[parts[1]]; ok {
		fmt.Println(v.Value)
	} else if v, ok := globalMemory[parts[1]]; ok {
		fmt.Println(v.Value)
	} else {
		fmt.Println(parts[1])
	}
}

func handleINPUT(parts []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("? ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	dst := getMemory()[parts[1]]
	switch dst.Type {
	case "INT16", "INT32", "INT64":
		n, _ := strconv.ParseInt(text, 10, 64)
		dst.Value = maskInt(n, dst.Type)
	case "FLOAT16", "FLOAT32", "FLOAT64":
		f, _ := strconv.ParseFloat(text, 64)
		dst.Value = f
	case "CHAR":
		dst.Value = text
	}
}

func handleCompare(parts []string, op string) {
	va := getVarVal(parts[1])
	vb := getVarVal(parts[2])
	switch op {
	case "LT":
		cmpFlag = compare(va, vb) < 0
	case "LE":
		cmpFlag = compare(va, vb) <= 0
	case "GT":
		cmpFlag = compare(va, vb) > 0
	case "GE":
		cmpFlag = compare(va, vb) >= 0
	case "EQ":
		cmpFlag = compare(va, vb) == 0
	case "NE":
		cmpFlag = compare(va, vb) != 0
	}
}

// --- Function Call ---
func callFunction(name string, args []string) interface{} {
	fnPC, ok := functions[name]
	if !ok {
		fmt.Println("Error: function not found", name)
		return nil
	}
	// push new local scope
	localScope := make(map[string]*Var)
	scopeStack = append(scopeStack, localScope)

	// assign parameters p0, p1, ...
	for i, arg := range args {
		localScope[fmt.Sprintf("p%d", i)] = &Var{Type: "INT64", Value: getVarVal(arg)}
	}

	// save return PC
	callStack = append(callStack, pc)
	pc = fnPC

	// execute until RETURN/RET
	for pc < len(program) {
		line := program[pc]
		pc++
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		op := strings.ToUpper(parts[0])
		if op == "RETURN" {
			returnValue = &Var{Type: "INT64", Value: getVarVal(parts[1])}
			pc = callStack[len(callStack)-1]
			callStack = callStack[:len(callStack)-1]
			val := returnValue.Value
			scopeStack = scopeStack[:len(scopeStack)-1]
			return val
		} else if op == "RET" {
			pc = callStack[len(callStack)-1]
			callStack = callStack[:len(callStack)-1]
			scopeStack = scopeStack[:len(scopeStack)-1]
			return nil
		} else {
			executeLine(parts)
		}
	}
	scopeStack = scopeStack[:len(scopeStack)-1]
	return nil
}

// --- Execute Single Line Helper ---
func executeLine(parts []string) {
	op := strings.ToUpper(parts[0])
	switch op {
	case "VAR":
		handleVAR(parts)
	case "FREE":
		handleFREE(parts)
	case "ADD", "SUB", "MUL", "DIV":
		handleArithmetic(parts, op)
	case "PRINT":
		handlePRINT(parts)
	case "INPUT":
		handleINPUT(parts)
	case "LT", "LE", "GT", "GE", "EQ", "NE":
		handleCompare(parts, op)
	}
}

// --- Help ---
func showHelp() {
	fmt.Println(version)
	fmt.Println("\nUsage: asp <file.asp>  (or asp --help, --version)")
	fmt.Println("Features:")
	fmt.Println("- VAR <TYPE> <name> [value] : declares variable (int16/32/64, float16/32/64, char)")
	fmt.Println("- FREE <name> : deletes variable")
	fmt.Println("- Arithmetic: ADD, SUB, MUL, DIV")
	fmt.Println("- Comparisons: LT, LE, GT, GE, EQ, NE")
	fmt.Println("- Branching: JMP, JZ, JNZ, HALT")
	fmt.Println("- Functions: FUNC, ENDFUNC, CALL, RETURN")
	fmt.Println("- Input/Output: INPUT, PRINT")
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
		fmt.Println("Usage: asp <program.asp>")
		return
	}

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

