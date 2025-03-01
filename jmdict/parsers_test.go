package jmdict

import (
	"encoding/json"
	"tango/utils"
	"testing"
)

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
				Kanji:      utils.ToStringPtr("日本"),
				Kana:       utils.ToStringPtr("にほん"),
				SenseIndex: utils.ToIntPtr(1.0),
			},
		},
		{
			"[kanji, kana]",
			`["世界", "せかい"]`,
			Xref{
				Kanji: utils.ToStringPtr("世界"),
				Kana:  utils.ToStringPtr("せかい"),
			},
		},
		{
			"[kanjiOrKana, senseIndex]",
			`["日本", 2]`,
			Xref{
				KanjiOrKana: utils.ToStringPtr("日本"),
				SenseIndex:  utils.ToIntPtr(2.0),
			},
		},
		{
			"[kanjiOrKana]",
			`["こんにちは"]`,
			Xref{
				KanjiOrKana: utils.ToStringPtr("こんにちは"),
			},
		},
	}

	for _, testCase := range testCases {
		var xref Xref
		err := json.Unmarshal([]byte(testCase.jsonData), &xref)
		if err != nil {
			t.Errorf("Error unmarshalling JSON: %v", err)
		}

		if !utils.EqualPointers(xref.Kanji, testCase.expected.Kanji) {
			t.Errorf(
				"Unexpected Kanji value. Expected: %v, Got: %v",
				testCase.expected.Kanji,
				xref.Kanji,
			)
		}

		if !utils.EqualPointers(xref.Kana, testCase.expected.Kana) {
			t.Errorf(
				"Unexpected Kana value. Expected: %v, Got: %v",
				testCase.expected.Kana,
				xref.Kana,
			)
		}

		if !utils.EqualPointers(xref.KanjiOrKana, testCase.expected.KanjiOrKana) {
			t.Errorf("Unexpected KanjiOrKana value. Expected: %v, Got: %v",
				testCase.expected.KanjiOrKana,
				xref.KanjiOrKana,
			)
		}

		if !utils.EqualPointers(xref.SenseIndex, testCase.expected.SenseIndex) {
			t.Errorf(
				"Unexpected SenseIndex value. Expected: %d, Got: %d",
				testCase.expected.SenseIndex,
				xref.SenseIndex,
			)
		}
	}
}
