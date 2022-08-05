package rulex

import (
	"strings"
	"testing"
)

var ExprCase = [][2]string{
	{"a|b", "a b |"},
	{"a|b&c", "a b | c &"},
	{"a|b&(c)", "a b | c &"},
	{"a|b&(c|d)", "a b | c d | &"},
	{"a|b&!(c|d)", "a b | c d | ! &"},
	{"a|!b&!(c|d)", "a b ! | c d | ! &"},
	{"a|!!!!!b&!(c|d)", "a b ! ! ! ! ! | c d | ! &"},
	{"你|我&!(他|她)", "你 我 | 他 她 | ! &"}, // support utf-8
	{"a|(b&(c&!(d|e))&(f|g)|(h&!i))&!(j|!k)", "a b c d e | ! & & f g | & h i ! & | | j k ! | ! &"},
}

func TestRPN(t *testing.T) {
	for i, c := range ExprCase {
		rpn, err := RPN(c[0], nil)
		if err != nil {
			t.Fatal(err)
		}
		actual, expected := strings.Join(rpn, " "), c[1]
		t.Logf("rpn %d - result: %s", i, actual)
		if actual != expected {
			t.Fatalf("test failed, actual: %s, expected: %s", actual, expected)
		}
	}
}
