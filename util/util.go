package util

func EnsureSlice(value interface{}) []string {
	switch v := value.(type) {
	case []interface{}:
		result := make([]string, len(v))
		for i, val := range v {
			if str, ok := val.(string); ok {
				result[i] = str
			}
		}
		return result
	default:
		return []string{}
	}
}

func ToStringPtr(v interface{}) *string {
	if str, ok := v.(string); ok {
		return &str
	}
	return nil
}

func ToIntPtr(v interface{}) *int {
	if floatNum, ok := v.(float64); ok {
		num := int(floatNum)
		return &num
	}
	return nil
}

func EqualPointers[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
