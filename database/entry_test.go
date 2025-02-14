package main

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalEntry(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectedErr bool
		expected    BleveEntry
	}{
		{
			name: "Valid JSON with all fields present",
			jsonData: `{
				"id": "12345",
				"kanji": ["日本語", "日本"],
				"kana": ["にほんご", "にほん"],
				"meanings": ["Japanese language", "Japan"]
			}`,
			expectedErr: false,
			expected: BleveEntry{
				ID:       "12345",
				Kanji:    []string{"日本語", "日本"},
				Kana:     []string{"にほんご", "にほん"},
				Meanings: []string{"Japanese language", "Japan"},
			},
		},
		{
			name: "JSON with missing optional fields",
			jsonData: `{
				"id": "67890"
			}`,
			expectedErr: false,
			expected: BleveEntry{
				ID:       "67890",
				Kanji:    []string{},
				Kana:     []string{},
				Meanings: []string{},
			},
		},
		{
			name: "JSON with single string for Kanji, Kana, Meanings",
			jsonData: `{
				"id": "54321",
				"kanji": "日本語",
				"kana": "にほんご",
				"meanings": "Japanese language"
			}`,
			expectedErr: false,
			expected: BleveEntry{
				ID:       "54321",
				Kanji:    []string{"日本語"},
				Kana:     []string{"にほんご"},
				Meanings: []string{"Japanese language"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var be BleveEntry
			err := json.Unmarshal([]byte(testCase.jsonData), &be)

			if (err != nil) != testCase.expectedErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, testCase.expectedErr)
				return
			}

			// TODO: Compare structs with string[] inside
			// if !testCase.expectedErr && be != testCase.expected {
			// 	t.Errorf("UnmarshalJSON() = %v, want %v", be, testCase.expected)
			// }
		})
	}
}
