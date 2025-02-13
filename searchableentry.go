package main

import "encoding/json"

func (se *SearchableEntry) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to unmarshal the JSON data
	type Alias SearchableEntry
	temp := &struct {
		Kanji    interface{} `json:"kanji"`
		Kana     interface{} `json:"kana"`
		Meanings interface{} `json:"meanings"`
		*Alias
	}{
		Alias: (*Alias)(se),
	}

	// Unmarshal the JSON data into the temporary struct
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	// Convert Kanji, Kana, and Meanings to slices of strings
	se.Kanji = EnsureSlice(temp.Kanji)
	se.Kana = EnsureSlice(temp.Kana)
	se.Meanings = EnsureSlice(temp.Meanings)

	return nil
}

func EnsureSlice(value interface{}) []string {
	switch v := value.(type) {
	case []interface{}: // string[] comes as interface{}
		result := make([]string, len(v))
		for i, val := range v {
			if str, ok := val.(string); ok {
				result[i] = str
			}
		}
		return result
	case string:
		return []string{v}
	default:
		return []string{}
	}
}
