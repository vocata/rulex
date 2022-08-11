package rulex

import "errors"

var (
	ErrInvalidSyntax  = errors.New("invalid syntax")
	ErrInvalidOperand = errors.New("invalid operand")
	ErrTagNotFound    = errors.New("tag not found")
)
