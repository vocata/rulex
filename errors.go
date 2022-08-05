package rulex

import "errors"

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
	ErrCondNotMatch  = errors.New("condition not match")
)
