package rulex

import (
	"fmt"
	"unicode"
)

var priorityTable = [6][6]byte{
	//        0    &    |    !    (    )
	/* 0 */ {'=', '<', '<', '<', '<', ' '},
	/* & */ {'>', '>', '>', '<', '<', '>'},
	/* | */ {'>', '>', '>', '<', '<', '>'},
	/* ! */ {'>', '>', '>', '<', '<', '>'},
	/* ( */ {' ', '<', '<', '<', '<', '='},
	/* ) */ {' ', ' ', ' ', ' ', ' ', ' '},
}

var positionTable = [7][7]byte{
	//       0  &  |  !  (  )  o
	/* 0 */ {0, 0, 0, 1, 1, 0, 1},
	/* & */ {0, 0, 0, 1, 1, 0, 1},
	/* | */ {0, 0, 0, 1, 1, 0, 1},
	/* ! */ {0, 0, 0, 1, 1, 0, 1},
	/* ( */ {0, 0, 0, 1, 1, 0, 1},
	/* ) */ {1, 1, 1, 0, 0, 1, 0},
	/* o */ {1, 1, 1, 0, 0, 1, 0},
}

// getpriorityTableIndex returns index of operator in priority table
func getpriorityTableIndex(op rune) int {
	switch op {
	case 0:
		return 0
	case '&':
		return 1
	case '|':
		return 2
	case '!':
		return 3
	case '(':
		return 4
	case ')':
		return 5
	}
	return -1
}

func orderBetween(left, right rune) byte {
	top_idx := getpriorityTableIndex(left)
	cur_idx := getpriorityTableIndex(right)

	return priorityTable[top_idx][cur_idx]
}

func getPositionTableIndex(item rune) int {
	switch item {
	case 0:
		return 0
	case '&':
		return 1
	case '|':
		return 2
	case '!':
		return 3
	case '(':
		return 4
	case ')':
		return 5
	}
	return 6
}

func validPosition(former, latter rune) bool {
	former_idx := getPositionTableIndex(former)
	latter_idx := getPositionTableIndex(latter)

	return positionTable[former_idx][latter_idx] == 1
}

// ValidateOperandFunc is type of func to validate the name of operand
type ValidateOperandFunc func(string) bool

// RPN converts infix expression to Reverse Polish Notation expression
func RPN(expr string, fn ValidateOperandFunc) ([]string, error) {
	// validate expression
	exprUTF8, err := validate(expr)
	if err != nil {
		return nil, err
	}

	// sentinel 0
	exprUTF8 = append(exprUTF8, 0)
	operatorStack := NewStack(rune(0))

	var exprRPN []string
	var begin, end, count int
	for !operatorStack.Empty() {
		if begin, end = getNext(exprUTF8, begin); begin == end {
			break
		}

		item := exprUTF8[begin:end]
		if len(item) == 1 && getpriorityTableIndex(item[0]) != -1 {
			operator := item[0]

			switch orderBetween(operatorStack.Top().(rune), operator) {
			case '<':
				operatorStack.Push(operator)
				begin = end
			case '>':
				operator := operatorStack.Pop().(rune)
				exprRPN = append(exprRPN, string(operator))
			case '=':
				operatorStack.Pop()
				begin = end
			default:
				return nil, fmt.Errorf("%w, no matching operator '%c', expr: '%s', idx: %d", ErrInvalidSyntax, operator, expr, begin)
			}
		} else {
			operand := string(item)

			if fn != nil && !fn(operand) {
				return nil, fmt.Errorf("%w, missing '%s' in conditions, expr: '%s', idx: %d", ErrCondNotMatch, operand, expr, begin)
			}

			exprRPN = append(exprRPN, operand)
			begin = end
			count++
		}
	}

	return exprRPN, nil
}

func validate(expr string) ([]rune, error) {
	// compatible with utf-8
	exprUTF8 := []rune(expr)
	if err := validateParentheses(expr, exprUTF8); err != nil {
		return nil, err
	}
	if err := validatePosition(expr, exprUTF8); err != nil {
		return nil, err
	}
	return exprUTF8, nil
}

func validateParentheses(expr string, exprUTF8 []rune) error {
	stk := NewStack()
	for i, r := range exprUTF8 {
		if r == '(' {
			stk.Push(i)
		}
		if r == ')' {
			if stk.Empty() {
				return fmt.Errorf("%w, expecting '(', expr: '%s', idx: %d", ErrInvalidSyntax, expr, i)
			}
			stk.Pop()
		}
	}
	if !stk.Empty() {
		i := stk.Pop().(int)
		return fmt.Errorf("%w, expecting ')', expr: '%s', idx: %d", ErrInvalidSyntax, expr, i)
	}
	return nil
}

func validatePosition(expr string, exprUTF8 []rune) error {
	exprUTF8 = append(exprUTF8, 0) // sentinel 0
	var former, latter = []rune{0}, []rune{}
	var last int
	for begin, end := getNext(exprUTF8, 0); begin != end; begin, end = getNext(exprUTF8, end) {
		latter = exprUTF8[begin:end]
		if !validPosition(former[0], latter[0]) {
			if former[0] == 0 && latter[0] == 0 {
				return fmt.Errorf("%w, empty expr", ErrInvalidSyntax)
			} else if former[0] != 0 {
				return fmt.Errorf("%w, incorrect item: '%s', expr: '%s', idx: %d", ErrInvalidSyntax, string(former), expr, last)
			} else {
				return fmt.Errorf("%w, incorrect item: '%s', expr: '%s', idx: %d", ErrInvalidSyntax, string(latter), expr, begin)
			}
		}
		former = latter
		last = begin
	}
	return nil
}

func getNext(expr []rune, begin int) (int, int) {
	// trim left spaces
	for begin < len(expr) {
		if !unicode.IsSpace(expr[begin]) {
			break
		}
		begin++
	}

	for i, r := range expr[begin:] {
		if unicode.IsSpace(r) {
			return begin, begin + i
		}
		if getpriorityTableIndex(r) != -1 {
			if i == 0 {
				return begin, begin + 1
			}
			return begin, begin + i
		}
	}
	return begin, len(expr)
}
