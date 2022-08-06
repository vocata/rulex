package rulex

import (
	"errors"
	"strings"
	"testing"
)

var ExprCase = []struct {
	Expr string
	RPN  string
	Err  error
}{
	{"a|b", "a b |", nil},
	{"a|b&c", "a b | c &", nil},
	{"a|b&(c)", "a b | c &", nil},
	{"a|b&(c|d)", "a b | c d | &", nil},
	{"a|b&!(c|d)", "a b | c d | ! &", nil},
	{"a|!b&!(c|d)", "a b ! | c d | ! &", nil},
	{"a|!!!!!b&!(c|d)", "a b ! ! ! ! ! | c d | ! &", nil},
	{"你|我&!(他|她)", "你 我 | 他 她 | ! &", nil}, // support utf-8
	{"a|(b&(c&!(d|e))&(f|g)|(h&!i))&!(j|!k)", "a b c d e | ! & & f g | & h i ! & | | j k ! | ! &", nil},
	{"a|b|", "", ErrInvalidSyntax},
	{"a|b!c", "", ErrInvalidSyntax},
	{"|a&b!c", "", ErrInvalidSyntax},
	{"", "", ErrInvalidSyntax},
	{"(a|b", "", ErrInvalidSyntax},
	{"a", "a", nil},
	{"a | b", "a b |", nil},
	{"a | b & c", "a b | c &", nil},
	{"a | b & (c)", "a b | c &", nil},
	{"a | b & (c | d)", "a b | c d | &", nil},
	{"a | b & ! (c | d)", "a b | c d | ! &", nil},
	{"a | !b& ! (c | d)", "a b ! | c d | ! &", nil},
	{"a |! ! !!!b &!(c | d)", "a b ! ! ! ! ! | c d | ! &", nil},
	{"你 | 我&! (他| 她)", "你 我 | 他 她 | ! &", nil}, // support utf-8
	{"a|(b &(c & !(d| e))&(f| g)|(h&   ! i))&!  (  j|!k)", "a b c d e | ! & & f g | & h i ! & | | j k ! | ! &", nil},
	{"a| b|", "", ErrInvalidSyntax},
	{"a | b!c", "", ErrInvalidSyntax},
	{"| a& b !  c", "", ErrInvalidSyntax},
	{"  ", "", ErrInvalidSyntax},
	{"a  ", "a", nil},
	{" a  ", "a", nil},
	{" ( a |b", "", ErrInvalidSyntax},
	{"  a |b)", "", ErrInvalidSyntax},
	{"  a) |b)", "", ErrInvalidSyntax},
	{"  a |)b)", "", ErrInvalidSyntax},
	{"a b", "", ErrInvalidSyntax},
	{"!", "", ErrInvalidSyntax},
	{"(! )a", "", ErrInvalidSyntax},
	{"a(! )a", "", ErrInvalidSyntax},
}

func TestRPN(t *testing.T) {
	for i, c := range ExprCase {
		rpn, err := RPN(c.Expr, nil)
		t.Logf("rpn %d - error: %v", i, err)
		if !errors.Is(err, c.Err) {
			t.Fatalf("rpn %d - test failed, error occurs, actual: %v, expected: %v", i, errors.Unwrap(err), c.Err)
		}

		actual, expected := strings.Join(rpn, " "), c.RPN
		t.Logf("rpn %d - result: %s", i, actual)
		if actual != expected {
			t.Fatalf("rpn %d - test failed, actual: %s, expected: %s", i, actual, expected)
		}
	}
}
