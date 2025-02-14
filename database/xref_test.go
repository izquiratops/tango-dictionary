package main

import (
	"encoding/json"
	"testing"
)

func EqualPointers[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func TestUnmarshalXRef(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		expected Xref
	}{
		{
			"[kanji, kana, senseIndex]",
			`["日本", "にほん", 1]`,
			Xref{Kanji: ToStringPtr("日本"), Kana: ToStringPtr("にほん"), SenseIndex: ToIntPtr(1.0)},
		},
		{
			"[kanji, kana]",
			`["世界", "せかい"]`,
			Xref{Kanji: ToStringPtr("世界"), Kana: ToStringPtr("せかい")},
		},
		{
			"[kanjiOrKana, senseIndex]",
			`["日本", 2]`,
			Xref{KanjiOrKana: ToStringPtr("日本"), SenseIndex: ToIntPtr(2.0)},
		},
		{
			"[kanjiOrKana]",
			`["こんにちは"]`,
			Xref{KanjiOrKana: ToStringPtr("こんにちは")},
		},
	}

	for _, testCase := range testCases {
		var xRef Xref
		err := json.Unmarshal([]byte(testCase.jsonData), &xRef)
		if err != nil {
			t.Errorf("Error unmarshalling JSON: %v", err)
		}

		if !EqualPointers(xRef.Kanji, testCase.expected.Kanji) {
			t.Errorf("Unexpected Kanji value. Expected: %v, Got: %v", testCase.expected.Kanji, xRef.Kanji)
		}

		if !EqualPointers(xRef.Kana, testCase.expected.Kana) {
			t.Errorf("Unexpected Kana value. Expected: %v, Got: %v", testCase.expected.Kana, xRef.Kana)
		}

		if !EqualPointers(xRef.KanjiOrKana, testCase.expected.KanjiOrKana) {
			t.Errorf("Unexpected KanjiOrKana value. Expected: %v, Got: %v", testCase.expected.KanjiOrKana, xRef.KanjiOrKana)
		}

		if !EqualPointers(xRef.SenseIndex, testCase.expected.SenseIndex) {
			t.Errorf("Unexpected SenseIndex value. Expected: %d, Got: %d", testCase.expected.SenseIndex, xRef.SenseIndex)
		}
	}
}
