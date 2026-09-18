package help

func AsSlice(v any) ([]any, bool) {
	a, ok := v.([]any)
	return a, ok
}

func AsInt(v any) (int, bool) {
	switch x := v.(type) {
	case float64:
		return int(x), true
	case int:
		return x, true
	default:
		return 0, false
	}
}
