package rulex

import (
	"fmt"
	"unicode"
)

var priorityTable = [6][6]byte{
	//        &    |    !	 (    )    0
	/* & */ {'>', '>', '<', '<', '>', '>'},
	/* | */ {'>', '>', '<', '<', '>', '>'},
	/* ! */ {'>', '>', '<', '<', '>', '>'},
	/* ( */ {'<', '<', '<', '<', '=', ')'},
	/* ) */ {' ', ' ', ' ', ' ', ' ', ' '},
	/* 0 */ {'<', '<', '<', '<', '(', '='},
}

func getIndex(op rune) int {
	switch op {
	case '&':
		return 0
	case '|':
		return 1
	case '!':
		return 2
	case '(':
		return 3
	case ')':
		return 4
	case 0:
		return 5
	}
	return -1
}

// ValidateOperandFunc is type of func to validate the name of operand
type ValidateOperandFunc func(string) bool

// RPN converts infix expression to Reverse Polish Notation expression
func RPN(expr string, fn ValidateOperandFunc) ([]string, error) {
	// compatible with utf-8
	exprUTF8 := []rune(expr)

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
		if len(item) == 1 && getIndex(item[0]) != -1 {
			operator := item[0]

			switch order := orderBetween(operatorStack.Top().(rune), operator); order {
			case '<':
				operatorStack.Push(operator)
				begin = end
			case '>':
				operator := operatorStack.Pop().(rune)
				if unary(operator) {
					if count < 1 {
						return nil, fmt.Errorf("%w, '%c' requires one operand, expr: '%s'", ErrInvalidSyntax, operator, expr)
					}
				}
				if binary(operator) {
					if count < 2 {
						return nil, fmt.Errorf("%w, '%c' requires tow operand, expr: '%s'", ErrInvalidSyntax, operator, expr)
					}
					count--
				}
				exprRPN = append(exprRPN, string(operator))
			case '=':
				operatorStack.Pop()
				begin = end
			default:
				return nil, fmt.Errorf("%w, missing '%c', expr: '%s'", ErrInvalidSyntax, order, expr)
			}
		} else {
			operand := string(item)

			if fn != nil && !fn(operand) {
				return nil, fmt.Errorf("%w, missing '%s' in conditions, expr: '%s'", ErrCondNotMatch, operand, expr)
			}

			exprRPN = append(exprRPN, operand)
			begin = end
			count++
		}
	}
	if count == 0 {
		return nil, fmt.Errorf("%w, empty expr", ErrInvalidSyntax)
	}
	if count > 1 {
		return nil, fmt.Errorf("%w, too many operand", ErrInvalidSyntax)
	}

	return exprRPN, nil
}

func orderBetween(left, right rune) byte {
	top_idx := getIndex(left)
	cur_idx := getIndex(right)

	return priorityTable[top_idx][cur_idx]
}

func unary(op rune) bool {
	return op == '!'
}

func binary(op rune) bool {
	return op == '&' || op == '|'
}

func getNext(expr []rune, begin int) (int, int) {
	// trim left space
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
		if getIndex(r) != -1 {
			if i == 0 {
				return begin, begin + 1
			}
			return begin, begin + i
		}
	}
	return begin, len(expr)
}
