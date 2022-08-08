package rulex

import "errors"

var (
	ErrInvalidSyntax  = errors.New("invalid syntax")
	ErrInvalidOperand = errors.New("invalid operand")
	ErrCondNotMatch   = errors.New("condition not match")
)
