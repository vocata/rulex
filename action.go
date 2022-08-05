//go:build go1.18 && !go1.19
// +build go1.18,!go1.19

package rulex

type Comparable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string
}

// GT means Greater Then, >
func GT[T Comparable](expected T) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			return value > expected
		}
		return false
	})
}

// GET means Greater than or Equal To, >=
func GET[T Comparable](expected T) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			return value >= expected
		}
		return false
	})
}

// LT means Less Than, <
func LT[T Comparable](expected T) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			return value < expected
		}
		return false
	})
}

// LET means Less than or Equal To, <=
func LET[T Comparable](expected T) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			return value <= expected
		}
		return false
	})
}

// EQ means EQual to, ==
func EQ[T Comparable](expected T) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			return value == expected
		}
		return false
	})
}

// NEQ means Not EQual to, !=
func NEQ[T Comparable](expected T) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			return value != expected
		}
		return false
	})
}

// Between means one is within the interval, lc and rc represent left closed and right closed interval respectively
func Between[T Comparable](left, right T, lc, rc bool) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			if lc && rc {
				return left <= value && value <= right
			} else if lc {
				return left <= value && value < right
			} else if rc {
				return left < value && value <= right
			} else {
				return left < value && value < right
			}
		}
		return false
	})
}

// In means one is a member of list
func In[T Comparable](list []T) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			for _, v := range list {
				if value == v {
					return true
				}
			}
			return false
		}
		return false
	})
}

// NotIn means one is not a member of list, as opposed to IN
func NotIn[T Comparable](list []T) Action {
	return ActionFunc(func(actual interface{}) bool {
		if value, ok := actual.(T); ok {
			for _, v := range list {
				if value == v {
					return false
				}
			}
			return true
		}
		return false
	})
}
