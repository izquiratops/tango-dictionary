package utils

func EnsureSlice(value any) []string {
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

func ContainsString(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}
