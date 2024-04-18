package internal

import (
	"bytes"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type scanner struct {
	instruction []byte
	pos         int
}

var ErrIllegalOpCode = errors.New("illegal opcode")
var ErrNoMatchingTransformer = errors.New("no matching transformer found")
var ErrInvalidArgument = errors.New("invalid argument")

func (s *scanner) nextChar() (c byte, stop bool) {
	if s.pos >= len(s.instruction) {
		return 0, true
	}
	c = s.instruction[s.pos]
	s.pos++
	switch c {
	case '[', ']', '"', ',', '(', ')', '=', '~', '|':
		return c, true
	}
	return c, false
}

// InstructionScanner is a public struct that should be created with NewInstructionScanner.
// This struct provides the logic for parsing the instruction into a series of chained op codes and values.
type InstructionScanner struct {
	firstOp   *op
	currentOp *op
	*scanner
	inverseNextOp bool
}

func (ris *InstructionScanner) setOp(op *op) {
	if ris.firstOp == nil {
		ris.firstOp = op
	} else {
		ris.currentOp.next = op
	}
	ris.currentOp = op
}

func (ris *InstructionScanner) nextInverse() bool {
	if ris.inverseNextOp {
		ris.inverseNextOp = false
		return true
	}
	return false
}

func (ris *InstructionScanner) genericTokenToOpCode(token string) op {
	if len(token) == 0 {
		return op{
			opCode: opCodeNil,
		}
	}
	var isString bool
	var isFloat bool
	if token[0] == '"' {
		isString = true
		token = token[1:]
	} else {
		for _, c := range token {
			if !((c >= '0' && c <= '9') || c == '.' || c == '-') {
				isString = true
				break
			}
			if c == '.' {
				isFloat = true
			}

		}
	}
	var inverse bool
	var value any
	var code opCode
	var err error
	if isString {
		normalizedToken := strings.ToLower(token)
		if len(token) > 0 && token[0] == '~' {
			inverse = true
			token = token[1:]
		}
		if normalizedToken == "true" || normalizedToken == "false" {
			code = opCodeBool
			if token == "true" {
				value = true
			} else {
				value = false
			}
		} else {
			code = opCodeString
			value = token
		}
	} else if isFloat {
		code = opCodeFloat
		value, err = strconv.ParseFloat(token, 64)
		if err != nil {
			code = opCodeString
			value = token
		}
	} else {
		code = opCodeInt
		value, err = strconv.Atoi(token)
		if err != nil {
			code = opCodeString
			value = token
		}
	}
	if ris.inverseNextOp {
		inverse = true
		ris.inverseNextOp = false
	}

	return op{
		opCode:  code,
		inverse: inverse,
		value:   value,
	}
}

var ErrInvalidOpChain = errors.New("invalid operation chain")

type group struct {
	identifier string
	inverse    bool
}

type rule struct {
	groups       []group
	runOnNoMatch bool
	MemoizedMethod
}

// GetEvaluator gets a memoized rule chain evaluator that can be called
// Note: if Scan has not been called first, it will be called by this method
func (ris *InstructionScanner) GetEvaluator(methodTable map[string]RawMethod) (Evaluator, error) {
	var parsedRules []rule
	if ris.firstOp == nil {
		ris.Scan()
	}
	var activeOp = ris.firstOp
	for activeOp != nil {
		var groups []group
		groupOp, found := activeOp.Next(opCodeSet, opCodeString)
		if !found {
			if len(parsedRules) == 0 {
				return nil, ErrInvalidOpChain
			}
			break
		}
		switch groupOp.opCode {
		case opCodeString:
			groups = append(groups, group{
				identifier: groupOp.value.(string),
				inverse:    groupOp.inverse,
			})
		case opCodeSet:
			for _, v := range groupOp.children {
				var inverse = groupOp.inverse
				if v.inverse {
					inverse = true
				}
				groups = append(groups, group{
					identifier: v.value.(string),
					inverse:    inverse,
				})
			}
		default:
			return nil, ErrInvalidOpChain
		}
		runOp, found := groupOp.Next(opCodeRun)
		if !found {
			return nil, ErrInvalidOpChain
		}
		methodIdentifier := runOp.next
		t, ok := methodTable[methodIdentifier.value.(string)]
		if !ok {
			return nil, errors.Wrapf(ErrNoMatchingTransformer, "transformer for %s not found", runOp.value)
		}
		var args []Arg
		if methodIdentifier.next != nil && methodIdentifier.next.opCode == opCodeParams {
			params := methodIdentifier.next
			for _, arg := range params.children {
				switch arg.opCode {
				case opCodeNil, opCodeString, opCodeFloat, opCodeBool, opCodeInt:
					args = append(args, Arg{OpCode: arg.opCode, Value: arg.value})
				default:
					return nil, errors.Wrapf(ErrInvalidArgument, "%v is not a valid opcode for an argument", arg.opCode.String())
				}
			}
			activeOp = params.next
		} else {
			activeOp = methodIdentifier.next
		}
		memoedFunction, err := t(args...)
		if err != nil {
			return nil, errors.Wrapf(err, "could not compile method for rule %v", string(ris.instruction))
		}
		parsedRules = append(parsedRules, rule{
			groups:         groups,
			runOnNoMatch:   groupOp.inverse,
			MemoizedMethod: memoedFunction,
		})
	}
	// Memoize the instruction into an evaluator
	return func(value any, targetGroups ...string) (bool, error) {
		if len(targetGroups) == 0 {
			targetGroups = []string{"none"}
		}
		for _, targetGroup := range targetGroups {
			targetGroup = strings.ToLower(targetGroup)
			for _, v := range parsedRules {
				runOnNoMatch := v.runOnNoMatch
				for _, group := range v.groups {
					var err error
					if group.identifier == targetGroup {
						if !group.inverse {
							err = v.MemoizedMethod(value)
						}
						return true, err
					} else if group.inverse || group.identifier == "all" {
						runOnNoMatch = true
					}
				}
				if runOnNoMatch {
					err := v.MemoizedMethod(value)
					return true, err
				}
			}
		}
		return false, nil
	}, nil
}

// Scan parses the constructor string and turns it into an opcode chain
func (ris *InstructionScanner) Scan() {
	for {
		c, stop := ris.nextChar()
		if stop {
			switch c {
			case 0:
				return
			// Done parsing
			case '(':
				ris.setOp(&op{
					opCode:  opCodeParams,
					inverse: ris.nextInverse(),
					children: mapFilter(
						setDecoder{
							scanner:      ris.scanner,
							activeBuffer: &bytes.Buffer{},
						}.decode(),
						func(v string) (*op, bool) {
							o := ris.genericTokenToOpCode(v)
							return &o, true
						}),
				})
			case '[':
				ris.setOp(&op{
					opCode:  opCodeSet,
					inverse: ris.nextInverse(),
					children: mapFilter(
						setDecoder{
							scanner:      ris.scanner,
							activeBuffer: &bytes.Buffer{},
						}.decode(),
						func(v string) (*op, bool) {
							o := ris.genericTokenToOpCode(v)
							return &o, true
						}),
				})
			case '"':
				o := ris.genericTokenToOpCode(stringDecoder{
					scanner:      ris.scanner,
					activeBuffer: &bytes.Buffer{},
				}.decode())
				ris.setOp(&o)
			case '~':
				ris.inverseNextOp = true
			case '=':
				ris.setOp(&op{
					opCode: opCodeRun,
				})
			case '|':
				ris.setOp(&op{
					opCode: opCodeChain,
				})
			default:
				o := ris.genericTokenToOpCode(tokenDecoder{
					scanner:      ris.scanner,
					activeBuffer: &bytes.Buffer{},
				}.decode())
				ris.setOp(&o)
			}
		} else {
			// unread the byte
			ris.scanner.pos--
			o := ris.genericTokenToOpCode(tokenDecoder{
				scanner:      ris.scanner,
				activeBuffer: &bytes.Buffer{},
			}.decode())
			ris.setOp(&o)
		}
	}
}

// NewInstructionScanner creates a new scanner and sets the interior scanner buffer to the tag provided.
func NewInstructionScanner(tag string) *InstructionScanner {
	return &InstructionScanner{
		firstOp:   nil,
		currentOp: nil,
		scanner: &scanner{
			instruction: []byte(tag),
		},
		inverseNextOp: false,
	}
}
