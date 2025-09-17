package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const version = "Assemplex v2.1"

type VarType int

const (
	INT16 VarType = iota
	INT32
	INT64
	FLOAT32
	FLOAT64
	CHAR
)

type Variable struct {
	Type  VarType
	Const bool
	I16   int16
	I32   int32
	I64   int64
	F32   float32
	F64   float64
	Char  string
}

type Instruction struct {
	Op    string
	A     int
	B     int
	Value string
	IsVar bool
}

type State struct {
	Variables []*Variable
	VarIndex  map[string]int
	Program   []Instruction
	Labels    map[string]int
	Funcs     map[string]int
	CallStack []int
	CmpFlag   bool
}

func NewState() *State {
	return &State{
		Variables: []*Variable{},
		VarIndex:  make(map[string]int),
		Program:   []Instruction{},
		Labels:    make(map[string]int),
		Funcs:     make(map[string]int),
		CallStack: []int{},
		CmpFlag:   false,
	}
}

func addVar(st *State, name string, t VarType, isConst bool) int {
	idx := len(st.Variables)
	st.Variables = append(st.Variables, &Variable{Type: t, Const: isConst})
	st.VarIndex[name] = idx
	return idx
}

func parseVarType(s string) VarType {
	switch strings.ToUpper(s) {
	case "INT16":
		return INT16
	case "INT32":
		return INT32
	case "INT64":
		return INT64
	case "FLOAT32":
		return FLOAT32
	case "FLOAT64":
		return FLOAT64
	case "CHAR":
		return CHAR
	default:
		panic("Unknown type: " + s)
	}
}

// Runtime setVar for INPUT and initialization
func setVar(v *Variable, val string, isVar bool, st *State) {
	if v == nil || v.Const {
		return
	}
	if isVar {
		if idx, ok := st.VarIndex[val]; ok {
			src := st.Variables[idx]
			if src == nil {
				return
			}
			switch v.Type {
			case INT16:
				v.I16 = src.I16
			case INT32:
				v.I32 = src.I32
			case INT64:
				v.I64 = src.I64
			case FLOAT32:
				v.F32 = src.F32
			case FLOAT64:
				v.F64 = src.F64
			case CHAR:
				v.Char = src.Char
			}
		}
		return
	}
	switch v.Type {
	case INT16:
		n, _ := strconv.ParseInt(val, 10, 16)
		v.I16 = int16(n)
	case INT32:
		n, _ := strconv.ParseInt(val, 10, 32)
		v.I32 = int32(n)
	case INT64:
		n, _ := strconv.ParseInt(val, 10, 64)
		v.I64 = n
	case FLOAT32:
		f, _ := strconv.ParseFloat(val, 32)
		v.F32 = float32(f)
	case FLOAT64:
		f, _ := strconv.ParseFloat(val, 64)
		v.F64 = f
	case CHAR:
		if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
			val = val[1 : len(val)-1]
		}
		v.Char = val
	}
}

func printVar(v *Variable) {
	if v == nil {
		return
	}
	switch v.Type {
	case INT16:
		fmt.Println(v.I16)
	case INT32:
		fmt.Println(v.I32)
	case INT64:
		fmt.Println(v.I64)
	case FLOAT32:
		fmt.Println(v.F32)
	case FLOAT64:
		fmt.Println(v.F64)
	case CHAR:
		fmt.Println(v.Char)
	}
}

func arith(st *State, op string, a, b int) {
	va := st.Variables[a]
	vb := st.Variables[b]
	if va == nil || vb == nil {
		return
	}
	if va.Type == CHAR && vb.Type == CHAR && op == "ADD" {
		va.Char += vb.Char
		return
	}
	switch va.Type {
	case INT16:
		switch op {
		case "ADD":
			va.I16 += vb.I16
		case "SUB":
			va.I16 -= vb.I16
		case "MUL":
			va.I16 *= vb.I16
		case "DIV":
			if vb.I16 != 0 {
				va.I16 /= vb.I16
			}
		case "LT":
			st.CmpFlag = va.I16 < vb.I16
		case "LE":
			st.CmpFlag = va.I16 <= vb.I16
		case "GT":
			st.CmpFlag = va.I16 > vb.I16
		case "GE":
			st.CmpFlag = va.I16 >= vb.I16
		case "EQ":
			st.CmpFlag = va.I16 == vb.I16
		case "NE":
			st.CmpFlag = va.I16 != vb.I16
		}
	case INT32:
		switch op {
		case "ADD":
			va.I32 += vb.I32
		case "SUB":
			va.I32 -= vb.I32
		case "MUL":
			va.I32 *= vb.I32
		case "DIV":
			if vb.I32 != 0 {
				va.I32 /= vb.I32
			}
		case "LT":
			st.CmpFlag = va.I32 < vb.I32
		case "LE":
			st.CmpFlag = va.I32 <= vb.I32
		case "GT":
			st.CmpFlag = va.I32 > vb.I32
		case "GE":
			st.CmpFlag = va.I32 >= vb.I32
		case "EQ":
			st.CmpFlag = va.I32 == vb.I32
		case "NE":
			st.CmpFlag = va.I32 != vb.I32
		}
	case INT64:
		switch op {
		case "ADD":
			va.I64 += vb.I64
		case "SUB":
			va.I64 -= vb.I64
		case "MUL":
			va.I64 *= vb.I64
		case "DIV":
			if vb.I64 != 0 {
				va.I64 /= vb.I64
			}
		case "LT":
			st.CmpFlag = va.I64 < vb.I64
		case "LE":
			st.CmpFlag = va.I64 <= vb.I64
		case "GT":
			st.CmpFlag = va.I64 > vb.I64
		case "GE":
			st.CmpFlag = va.I64 >= vb.I64
		case "EQ":
			st.CmpFlag = va.I64 == vb.I64
		case "NE":
			st.CmpFlag = va.I64 != vb.I64
		}
	case FLOAT32:
		switch op {
		case "ADD":
			va.F32 += vb.F32
		case "SUB":
			va.F32 -= vb.F32
		case "MUL":
			va.F32 *= vb.F32
		case "DIV":
			if vb.F32 != 0 {
				va.F32 /= vb.F32
			}
		case "LT":
			st.CmpFlag = va.F32 < vb.F32
		case "LE":
			st.CmpFlag = va.F32 <= vb.F32
		case "GT":
			st.CmpFlag = va.F32 > vb.F32
		case "GE":
			st.CmpFlag = va.F32 >= vb.F32
		case "EQ":
			st.CmpFlag = va.F32 == vb.F32
		case "NE":
			st.CmpFlag = va.F32 != vb.F32
		}
	case FLOAT64:
		switch op {
		case "ADD":
			va.F64 += vb.F64
		case "SUB":
			va.F64 -= vb.F64
		case "MUL":
			va.F64 *= vb.F64
		case "DIV":
			if vb.F64 != 0 {
				va.F64 /= vb.F64
			}
		case "LT":
			st.CmpFlag = va.F64 < vb.F64
		case "LE":
			st.CmpFlag = va.F64 <= vb.F64
		case "GT":
			st.CmpFlag = va.F64 > vb.F64
		case "GE":
			st.CmpFlag = va.F64 >= vb.F64
		case "EQ":
			st.CmpFlag = va.F64 == vb.F64
		case "NE":
			st.CmpFlag = va.F64 != vb.F64
		}
	}
}

func execute(st *State) {
	pc := 0
	reader := bufio.NewReader(os.Stdin)
	for pc < len(st.Program) {
		inst := st.Program[pc]
		switch inst.Op {
		case "VAR", "CONST":
			if inst.Value != "" {
				setVar(st.Variables[inst.A], inst.Value, inst.IsVar, st)
			}
		case "FREE":
			st.Variables[inst.A] = nil
		case "INPUT":
			fmt.Print("? ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			setVar(st.Variables[inst.A], text, false, st)
		case "PRINT":
			printVar(st.Variables[inst.A])
		case "ADD", "SUB", "MUL", "DIV", "LT", "LE", "GT", "GE", "EQ", "NE":
			arith(st, inst.Op, inst.A, inst.B)
		case "JMP":
			if idx, ok := st.Labels[inst.Value]; ok {
				pc = idx
				continue
			}
		case "JZ":
			if !st.CmpFlag {
				if idx, ok := st.Labels[inst.Value]; ok {
					pc = idx
					continue
				}
			}
		case "JNZ":
			if st.CmpFlag {
				if idx, ok := st.Labels[inst.Value]; ok {
					pc = idx
					continue
				}
			}
		case "CALL":
			if fidx, ok := st.Funcs[inst.Value]; ok {
				st.CallStack = append(st.CallStack, pc+1)
				pc = fidx
				continue
			}
		case "RET":
			if len(st.CallStack) == 0 {
				return
			}
			pc = st.CallStack[len(st.CallStack)-1]
			st.CallStack = st.CallStack[:len(st.CallStack)-1]
			continue
		case "HALT":
			return
		}
		pc++
	}
}

func parseLine(st *State, line string, lineIndex int) *Instruction {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, ";") {
		return nil
	}
	if strings.HasSuffix(line, ":") {
		label := strings.TrimSuffix(line, ":")
		st.Labels[label] = lineIndex
		return nil
	}
	parts := strings.Fields(line)
	op := strings.ToUpper(parts[0])
	switch op {
	case "FUNC":
		if len(parts) >= 2 {
			fname := parts[1]
			st.Funcs[fname] = lineIndex + 1
		}
		return nil
	case "ENDFUNC":
		return nil
	case "VAR", "CONST":
		if len(parts) >= 3 {
			t := parseVarType(parts[1])
			isConst := (op == "CONST")
			idx := addVar(st, parts[2], t, isConst)
			val := ""
			isVar := false
			if len(parts) == 4 {
				val = parts[3]
				_, exists := st.VarIndex[val]
				isVar = exists
			}
			return &Instruction{Op: op, A: idx, Value: val, IsVar: isVar}
		}
	case "FREE", "INPUT", "PRINT":
		idx := st.VarIndex[parts[1]]
		return &Instruction{Op: op, A: idx}
	case "ADD", "SUB", "MUL", "DIV", "LT", "LE", "GT", "GE", "EQ", "NE":
		return &Instruction{Op: op, A: st.VarIndex[parts[1]], B: st.VarIndex[parts[2]]}
	case "JMP", "JZ", "JNZ", "CALL":
		return &Instruction{Op: op, Value: parts[1]}
	case "RET", "HALT":
		return &Instruction{Op: op}
	}
	return nil
}

func showHelp() {
	fmt.Println(version)
	fmt.Println("\nUsage:")
	fmt.Println("  asp <file.asp>       Run Assemplex program")
	fmt.Println("  asp --version        Show version")
	fmt.Println("  asp --help           Show this help\n")
	fmt.Println("Variable types: INT16, INT32, INT64, FLOAT32, FLOAT64, CHAR")
	fmt.Println("Instructions: VAR, CONST, FREE, ADD, SUB, MUL, DIV, LT, LE, GT, GE, EQ, NE, PRINT, INPUT, JMP, JZ, JNZ, CALL, RET, HALT")
	fmt.Println("Functions: FUNC <name> ... ENDFUNC, CALL <name>, RET")
	fmt.Println("Labels: use <label>: at start of line, JMP/JZ/JNZ <label>")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: asp <program.asp>")
		return
	}
	switch os.Args[1] {
	case "--help":
		showHelp()
		return
	case "--version":
		fmt.Println(version)
		return
	}

	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	st := NewState()
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if inst := parseLine(st, line, i); inst != nil {
			st.Program = append(st.Program, *inst)
		}
	}

	execute(st)
}
