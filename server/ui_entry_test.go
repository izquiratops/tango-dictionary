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
			name: "word with kanji and reading",
			input: []model.JMdictWord{
				{
					Kanji: []model.JMdictKanji{
						{Text: "食べる", Common: true},
					},
					Kana: []model.JMdictKana{
						{Text: "たべる"},
					},
					Sense: []model.JMdictSense{
						{
							Gloss: []model.JMdictGloss{
								{Lang: "eng", Text: "to eat"},
							},
						},
					},
				},
			},
			expected: []UIEntry{
				{
					OtherForms: []Furigana{
						{Word: "食べる", Reading: "たべる"},
					},
					IsCommon: true,
					Meanings: []string{"to eat"},
				},
			},
		},
		{
			name: "kana only word",
			input: []model.JMdictWord{
				{
					Kana: []model.JMdictKana{
						{Text: "あそこ", Common: true},
					},
					Sense: []model.JMdictSense{
						{
							Gloss: []model.JMdictGloss{
								{Lang: "eng", Text: "there"},
								{Lang: "eng", Text: "that place"},
							},
						},
					},
				},
			},
			expected: []UIEntry{
				{
					OtherForms: []Furigana{
						{Word: "あそこ", Reading: ""},
					},
					IsCommon: true,
					Meanings: []string{"there", "that place"},
				},
			},
		},
		{
			name: "word with restricted senses",
			input: []model.JMdictWord{
				{
					Kanji: []model.JMdictKanji{
						{Text: "開く", Common: true},
						{Text: "明く"},
					},
					Kana: []model.JMdictKana{
						{Text: "あく", AppliesToKanji: []string{"開く"}},
						{Text: "ひらく", AppliesToKanji: []string{"開く"}},
					},
					Sense: []model.JMdictSense{
						{
							AppliesToKanji: []string{"開く"},
							Gloss: []model.JMdictGloss{
								{Lang: "eng", Text: "to open"},
							},
						},
					},
				},
			},
			expected: []UIEntry{
				{
					OtherForms: []Furigana{
						{Word: "開く", Reading: "あく"},
						{Word: "明く", Reading: "あく"},
					},
					IsCommon: true,
					Meanings: []string{"to open"},
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
