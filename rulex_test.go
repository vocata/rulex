package rulex

import (
	"strings"
	"testing"
)

func TestPredefinedAction(t *testing.T) {
	cond := NewCondition()
	cond.Add("a", "height", GT(165))
	cond.Add("b", "height", LT(180))
	cond.Add("c", "gender", In([]string{"male", "female"}))

	actual := map[string]interface{}{
		"height": 175,
		"gender": "male",
	}

	ruleX, err := NewRuleX("a & b & c", cond)
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
	cond.Add("a", "chromium", ActionFunc(func(actual interface{}) bool {
		return strings.HasPrefix(actual.(string), "http://")
	})).Add("b", "chromium", ActionFunc(func(actual interface{}) bool {
		return strings.HasPrefix(actual.(string), "https://")
	})).Add("c", "chromium", ActionFunc(func(actual interface{}) bool {
		return strings.HasSuffix(actual.(string), ".git")
	})).Add("d", "is_public", ActionFunc(func(actual interface{}) bool {
		return actual.(bool)
	}))
	actual := map[string]interface{}{
		"chromium":  "https://github.com/chromium/chromium.git",
		"is_public": true,
	}

	var expected bool
	// expression 1
	expr1 := "(a | b) & c & d"
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
	expr2 := "(a | b) & !!c & !!d"
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
	expr3 := "!!(a | b) & !!c & d"
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
	expr4 := "!(!a & !b) & !!c & d"
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
	expr5 := "a & b & c & d"
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
	expr6 := "!a & b & c & !!d"
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
