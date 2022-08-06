package rulex

import (
	"strings"
	"testing"
)

func TestPredefinedAction(t *testing.T) {
	cond := NewCondition()
	cond.Add("a", GT(10))
	cond.Add("b", LT(20))
	cond.Add("c", In([]int{1, 2, 3}))

	actual := map[string]interface{}{
		"a": 15,
		"b": 15,
		"c": 5,
	}

	ruleX, err := NewRuleX("a & b & !c", cond)
	if err != nil {
		t.Fatal(err)
	}
	result, err := ruleX.Eval(actual)
	if err != nil {
		t.Fatal(err)
	}
	expected := true
	t.Logf("expression - result: %t", result)
	if result != expected {
		t.Fatalf("test failed, actual: %t, expected: %t", result, expected)
	}
}

func TestSelfDefinedAction(t *testing.T) {
	cond := NewCondition()
	cond.Add("a", ActionFunc(func(actual interface{}) bool {
		return strings.HasPrefix(actual.(string), "http://")
	}))
	cond.Add("b", ActionFunc(func(actual interface{}) bool {
		return strings.HasPrefix(actual.(string), "https://")
	}))
	cond.Add("c", ActionFunc(func(actual interface{}) bool {
		return strings.HasSuffix(actual.(string), ".git")
	}))
	actual := map[string]interface{}{
		"a": "https://github.com/chromium/chromium.git",
		"b": "https://github.com/chromium/chromium.git",
		"c": "https://github.com/chromium/chromium.git",
	}

	var expected bool
	// expression 1
	expr1 := "(a | b) & c"
	ruleX1, err := NewRuleX(expr1, cond)
	if err != nil {
		t.Fatal(err)
	}
	result, err := ruleX1.Eval(actual)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("expression 1 - result: %t", result)
	expected = true
	if result != expected {
		t.Fatalf("test failed, actual: %t, expected: %t", result, expected)
	}

	// expression 2
	expr2 := "(a | b) & !!c"
	ruleX2, err := NewRuleX(expr2, cond)
	if err != nil {
		t.Fatal(err)
	}
	result, err = ruleX2.Eval(actual)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("expression 2 - result: %t", result)
	expected = true
	if result != expected {
		t.Fatalf("test failed, actual: %t, expected: %t", result, expected)
	}

	// expression 3
	expr3 := "!!(a | b) & !!c"
	ruleX3, err := NewRuleX(expr3, cond)
	if err != nil {
		t.Fatal(err)
	}
	result, err = ruleX3.Eval(actual)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("expression 3 - result: %t", result)
	expected = true
	if result != expected {
		t.Fatalf("test failed, actual: %t, expected: %t", result, expected)
	}

	// expression 4
	expr4 := "!(!a & !b) & !!c"
	ruleX4, err := NewRuleX(expr4, cond)
	if err != nil {
		t.Fatal(err)
	}
	result, err = ruleX4.Eval(actual)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("expression 4 - result: %t", result)
	expected = true
	if result != expected {
		t.Fatalf("test failed, actual: %t, expected: %t", result, expected)
	}

	// expression 5
	expr5 := "a & b & c"
	ruleX5, err := NewRuleX(expr5, cond)
	if err != nil {
		t.Fatal(err)
	}
	result, err = ruleX5.Eval(actual)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("expression 5 - result: %t", result)
	expected = false
	if result != expected {
		t.Fatalf("test failed, actual: %t, expected: %t", result, expected)
	}

	// expression 6
	expr6 := "!a & b & c"
	ruleX6, err := NewRuleX(expr6, cond)
	if err != nil {
		t.Fatal(err)
	}
	result, err = ruleX6.Eval(actual)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("expression 5 - result: %t", result)
	expected = true
	if result != expected {
		t.Fatalf("test failed, actual: %t, expected: %t", result, expected)
	}
}
