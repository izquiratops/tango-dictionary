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
