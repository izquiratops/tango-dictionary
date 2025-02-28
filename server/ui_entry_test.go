package server

import (
	"reflect"
	"tango/model"
	"testing"
)

func TestProcessEntries(t *testing.T) {
	tests := []struct {
		name     string
		input    []model.JMdictWord
		expected []UIEntry
	}{
		{
			name: "Test: 暖かい",
			input: []model.JMdictWord{
				{
					ID: "1586420",
					Kanji: []model.JMdictKanji{
						{Text: "暖かい", Common: true},
						{Text: "温かい", Common: true},
						{Text: "暖い", Common: false, Tags: []string{"sK"}},
					},
					Kana: []model.JMdictKana{
						{Text: "あたたかい", Common: true, AppliesToKanji: []string{"*"}},
						{Text: "あったかい", Common: false, AppliesToKanji: []string{"*"}},
						{Text: "あったけー", Common: false, Tags: []string{"sK"}, AppliesToKanji: []string{"*"}},
					},
					Sense: []model.JMdictSense{
						{
							PartOfSpeech: []string{"adj-i"},
							Gloss: []model.JMdictGloss{
								{Lang: "eng", Text: "warm"},
								{Lang: "eng", Text: "mild"},
								{Lang: "eng", Text: "(pleasantly) hot"},
							},
						},
					},
				},
			},
			expected: []UIEntry{
				{
					MainWord: Furigana{Word: "暖かい", Reading: "あたたかい"},
					OtherForms: []Furigana{
						{Word: "温かい", Reading: "あたたかい"},
						{Word: "暖い", Reading: "あったけー"},
					},
					IsCommon: true,
					Meanings: []string{
						"warm, mild, (pleasantly) hot",
						"considerate, kind, genial",
						"warm (of a colour), mellow",
						"having enough money",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessEntries(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ProcessEntries() = %v, want %v", got, tt.expected)
			}
		})
	}
}
