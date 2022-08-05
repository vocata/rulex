package rulex

import "testing"

func TestStackPushPop(t *testing.T) {
	var expected []int
	loop := 10
	stk := NewStack()
	for i := 0; i < loop; i++ {
		stk.Push(i)
		expected = append(expected, loop-i-1)
	}

	var actual []int
	for !stk.Empty() {
		v := stk.Pop()
		actual = append(actual, v.(int))
	}
	if len(actual) != loop {
		t.Fatalf("test failed, len(actual): %d, len(expected): %d", len(actual), len(expected))
	}

	equal := true
	for i := 0; i < loop; i++ {
		if actual[i] != 10-i-1 {
			equal = false
			break
		}
	}
	t.Logf("stack - result: %v", actual)
	if !equal {
		t.Fatalf("test failed, result: %v", actual)
	}
}
