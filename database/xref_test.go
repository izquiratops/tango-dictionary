package main

import (
	"encoding/json"
	JishoUtil "jisho-clone-database/util"
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
			Xref{
				Kanji:      JishoUtil.ToStringPtr("日本"),
				Kana:       JishoUtil.ToStringPtr("にほん"),
				SenseIndex: JishoUtil.ToIntPtr(1.0),
			},
		},
		{
			"[kanji, kana]",
			`["世界", "せかい"]`,
			Xref{
				Kanji: JishoUtil.ToStringPtr("世界"),
				Kana:  JishoUtil.ToStringPtr("せかい"),
			},
		},
		{
			"[kanjiOrKana, senseIndex]",
			`["日本", 2]`,
			Xref{
				KanjiOrKana: JishoUtil.ToStringPtr("日本"),
				SenseIndex:  JishoUtil.ToIntPtr(2.0),
			},
		},
		{
			"[kanjiOrKana]",
			`["こんにちは"]`,
			Xref{
				KanjiOrKana: JishoUtil.ToStringPtr("こんにちは"),
			},
		},
	}

	for _, testCase := range testCases {
		var xRef Xref
		err := json.Unmarshal([]byte(testCase.jsonData), &xRef)
		if err != nil {
			t.Errorf("Error unmarshalling JSON: %v", err)
		}

		if !EqualPointers(xRef.Kanji, testCase.expected.Kanji) {
			t.Errorf(
				"Unexpected Kanji value. Expected: %v, Got: %v",
				testCase.expected.Kanji,
				xRef.Kanji,
			)
		}

		if !EqualPointers(xRef.Kana, testCase.expected.Kana) {
			t.Errorf(
				"Unexpected Kana value. Expected: %v, Got: %v",
				testCase.expected.Kana,
				xRef.Kana,
			)
		}

		if !EqualPointers(xRef.KanjiOrKana, testCase.expected.KanjiOrKana) {
			t.Errorf("Unexpected KanjiOrKana value. Expected: %v, Got: %v",
				testCase.expected.KanjiOrKana,
				xRef.KanjiOrKana,
			)
		}

		if !EqualPointers(xRef.SenseIndex, testCase.expected.SenseIndex) {
			t.Errorf(
				"Unexpected SenseIndex value. Expected: %d, Got: %d",
				testCase.expected.SenseIndex,
				xRef.SenseIndex,
			)
		}
	}
}
