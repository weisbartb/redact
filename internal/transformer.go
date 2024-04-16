package internal

import "strconv"

// Arg is an argument parsed from an opcode.
type Arg struct {
	OpCode opCode
	Value  any
}

// String returns the argument as a string, this will coerce based on the opcode.
func (a Arg) String() string {
	switch a.OpCode {
	case opCodeInt:
		return strconv.Itoa(a.Value.(int))
	case opCodeBool:
		return strconv.FormatBool(a.Value.(bool))
	case opCodeString:
		return a.Value.(string)
	case opCodeFloat:
		return strconv.FormatFloat(a.Value.(float64), 'f', -1, 64)
	default:
		return ""
	}
}

// Int returns the argument as an int, this will coerce based on the opcode.
func (a Arg) Int() int {
	switch a.OpCode {
	case opCodeInt:
		return a.Value.(int)
	case opCodeBool:
		if a.Value.(bool) {
			return 1
		}
		return 0
	case opCodeString:
		v, _ := strconv.Atoi(a.Value.(string))
		return v
	case opCodeFloat:
		return int(a.Value.(float64))
	default:
		return 0
	}
}

// Float returns the argument as a float, this will coerce based on the opcode.
func (a Arg) Float() float64 {
	switch a.OpCode {
	case opCodeInt:
		return float64(a.Value.(int))
	case opCodeBool:
		if a.Value.(bool) {
			return 1
		}
		return 0
	case opCodeString:
		v, _ := strconv.ParseFloat(a.Value.(string), 64)
		return v
	case opCodeFloat:
		return a.Value.(float64)
	default:
		return 0
	}
}

// Bool returns the argument as a Bool, this will coerce based on the opcode.
func (a Arg) Bool() bool {
	switch a.OpCode {
	case opCodeInt:
		return a.Value.(int) != 0
	case opCodeBool:
		return a.Value.(bool)
	case opCodeString:
		v, _ := strconv.ParseBool(a.Value.(string))
		return v
	case opCodeFloat:
		return a.Value.(float64) != 0
	default:
		return false
	}
}

type Transformer func(value any)

type TransformerGenerator func(arguments ...Arg) Transformer

type Evaluator func(value any, groups ...string) (bool, error)
