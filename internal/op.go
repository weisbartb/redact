package internal

type opCode int

const (
	opCodeNone opCode = iota
	opCodeSet
	opCodeParams
	opCodeString
	opCodeFloat
	opCodeInt
	opCodeBool
	opCodeRun
	opCodeNil
	opCodeChain
)

func (o opCode) String() string {
	switch o {
	case opCodeNone:
		return "nil"
	case opCodeSet:
		return "set"
	case opCodeParams:
		return "params"
	case opCodeString:
		return "string"
	case opCodeFloat:
		return "float"
	case opCodeInt:
		return "int"
	case opCodeBool:
		return "bool"
	case opCodeRun:
		return "run"
	case opCodeNil:
		return "nil"
	case opCodeChain:
		return "chain"
	}
	return "unknown"
}

// op contains instruction sets for the basic parsing of a redaction string
type op struct {
	opCode   opCode
	inverse  bool
	children []*op
	next     *op
	// string, float64, bool, int
	value any
}

// Next transverses the linked list of ops till a matching opcode is encountered.
func (o *op) Next(codes ...opCode) (*op, bool) {
	n := o
	for {
		if n.next == nil {
			return nil, false
		}
		for _, code := range codes {
			if n.opCode == code {
				return n, true
			}
		}
		n = n.next
	}
}
