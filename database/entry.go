package main

import (
	"encoding/json"
	"fmt"
)

func (d *JMdictWord) ToBleveEntry() (BleveEntry, error) {
	searchable := BleveEntry{
		ID:       d.ID,
		Kanji:    make([]string, 0),
		Kana:     make([]string, 0),
		Meanings: make([]string, 0),
	}

	for _, k := range d.Kanji {
		if k.Text == "" {
			return searchable, fmt.Errorf("emtpy field at %v", d.ID)
		}

		searchable.Kanji = append(searchable.Kanji, k.Text)
	}

	for _, k := range d.Kana {
		if k.Text == "" {
			return searchable, fmt.Errorf("emtpy field at %v", d.ID)
		}

		searchable.Kana = append(searchable.Kana, k.Text)
	}

	for _, s := range d.Sense {
		for _, g := range s.Gloss {
			if g.Lang == "eng" {
				if g.Text == "" {
					return searchable, fmt.Errorf("emtpy field at %v", d.ID)
				}

				searchable.Meanings = append(searchable.Meanings, g.Text)
			}
		}
	}

	return searchable, nil
}

func (be *BleveEntry) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to unmarshal the JSON data
	type Alias BleveEntry
	temp := &struct {
		Kanji    interface{} `json:"kanji"`
		Kana     interface{} `json:"kana"`
		Meanings interface{} `json:"meanings"`
		*Alias
	}{
		Alias: (*Alias)(be),
	}

	// Unmarshal the JSON data into the temporary struct
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	// Convert Kanji, Kana, and Meanings to slices of strings
	be.Kanji = EnsureSlice(temp.Kanji)
	be.Kana = EnsureSlice(temp.Kana)
	be.Meanings = EnsureSlice(temp.Meanings)

	return nil
}

func EnsureSlice(value interface{}) []string {
	// TODO: Improve this function with reflect
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
