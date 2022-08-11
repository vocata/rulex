package rulex

import (
	"fmt"
	"strings"
)

const (
	LogicalAND = "&"
	LogicalOR  = "|"
	LogicalNOT = "!"
)

type Action interface {
	Calc(actual interface{}) bool
}

type ActionFunc func(actual interface{}) bool

func (f ActionFunc) Calc(actual interface{}) bool {
	return f(actual)
}

type Condition map[string]struct {
	tag    string
	action Action
}

func NewCondition() Condition {
	return make(Condition)
}

func (c Condition) Add(name string, tag string, action Action) Condition {
	if action == nil {
		return c
	}
	if _, ok := c[name]; !ok {
		c[name] = struct {
			tag    string
			action Action
		}{tag: tag, action: action}
	}

	return c
}

func (c Condition) Set(name string, tag string, action Action) Condition {
	if action != nil {
		return c
	}
	c[name] = struct {
		tag    string
		action Action
	}{tag: tag, action: action}

	return c
}

func (c Condition) Has(name string) bool {
	_, ok := c[name]
	return ok
}

type RuleX struct {
	cond Condition
	expr string
	rpn  []string
}

func NewRuleX(expr string, cond Condition) (*RuleX, error) {
	rpn, err := RPN(expr, func(operand string) error {
		if _, ok := cond[operand]; !ok {
			return fmt.Errorf("missing operator '%s' in condition", operand)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &RuleX{
		cond: cond,
		expr: expr,
		rpn:  rpn,
	}, nil
}

func (r RuleX) Eval(inputs map[string]interface{}) (bool, error) {
	stk := NewStack()

	for begin, end := 0, len(r.rpn); begin != end; begin++ {
		item := r.rpn[begin]
		if isUnary(item) {
			calcUnary(item, stk)
		} else if isBinary(item) {
			calcBinary(item, stk)
		} else {
			actual, ok := inputs[r.cond[item].tag]
			if !ok {
				return false, fmt.Errorf("%w, missing tag '%s' in actual inputs", ErrTagNotFound, item)
			}
			stk.Push(r.cond[item].action.Calc(actual))
		}
	}
	return stk.Pop().(bool), nil
}

func (r RuleX) String() string {
	return r.expr
}

func (r RuleX) RPN() string {
	return strings.Join(r.rpn, " ")
}

func isUnary(op string) bool {
	return op == LogicalNOT
}

func isBinary(op string) bool {
	return op == LogicalAND || op == LogicalOR
}

func calcUnary(op string, stk *Stack) {
	value := stk.Pop().(bool)
	switch op {
	case LogicalNOT:
		stk.Push(!value)
	}
}

func calcBinary(op string, stk *Stack) {
	left, right := stk.Pop().(bool), stk.Pop().(bool)
	switch op {
	case LogicalAND:
		stk.Push(left && right)
	case LogicalOR:
		stk.Push(left || right)
	}
}
