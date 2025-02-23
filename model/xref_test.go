package model

import (
	"encoding/json"
	"tango/util"
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
				Kanji:      util.ToStringPtr("日本"),
				Kana:       util.ToStringPtr("にほん"),
				SenseIndex: util.ToIntPtr(1.0),
			},
		},
		{
			"[kanji, kana]",
			`["世界", "せかい"]`,
			Xref{
				Kanji: util.ToStringPtr("世界"),
				Kana:  util.ToStringPtr("せかい"),
			},
		},
		{
			"[kanjiOrKana, senseIndex]",
			`["日本", 2]`,
			Xref{
				KanjiOrKana: util.ToStringPtr("日本"),
				SenseIndex:  util.ToIntPtr(2.0),
			},
		},
		{
			"[kanjiOrKana]",
			`["こんにちは"]`,
			Xref{
				KanjiOrKana: util.ToStringPtr("こんにちは"),
			},
		},
	}

	for _, testCase := range testCases {
		var xRef Xref
		err := json.Unmarshal([]byte(testCase.jsonData), &xRef)
		if err != nil {
			t.Errorf("Error unmarshalling JSON: %v", err)
		}

		if !util.EqualPointers(xRef.Kanji, testCase.expected.Kanji) {
			t.Errorf(
				"Unexpected Kanji value. Expected: %v, Got: %v",
				testCase.expected.Kanji,
				xRef.Kanji,
			)
		}

		if !util.EqualPointers(xRef.Kana, testCase.expected.Kana) {
			t.Errorf(
				"Unexpected Kana value. Expected: %v, Got: %v",
				testCase.expected.Kana,
				xRef.Kana,
			)
		}

		if !util.EqualPointers(xRef.KanjiOrKana, testCase.expected.KanjiOrKana) {
			t.Errorf("Unexpected KanjiOrKana value. Expected: %v, Got: %v",
				testCase.expected.KanjiOrKana,
				xRef.KanjiOrKana,
			)
		}

		if !util.EqualPointers(xRef.SenseIndex, testCase.expected.SenseIndex) {
			t.Errorf(
				"Unexpected SenseIndex value. Expected: %d, Got: %d",
				testCase.expected.SenseIndex,
				xRef.SenseIndex,
			)
		}
	}
}
