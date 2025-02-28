package server

import (
	"strings"
	"tango/model"
	"tango/util"
)

// UIEntry represents a simplified and processed version of a JMdictWord
// optimized for HTML template rendering.
type UIEntry struct {
	MainWord   Furigana   `json:"mainWord"`   // Primary word representation
	OtherForms []Furigana `json:"otherForms"` // Alternative forms of the word
	Common     bool       `json:"isCommon"`   // Indicates if word is frequently used
	Meanings   []string   `json:"meanings"`   // Word definitions/translations
}

type Furigana struct {
	Word    string `json:"word"`    // Kanji representation (or kana if no kanji exists)
	Reading string `json:"reading"` // Kana reading (empty for kana-only words)
}

func ProcessEntries(words []model.JMdictWord) []UIEntry {
	entries := make([]UIEntry, 0, len(words))

	for _, word := range words {
		entry := UIEntry{}

		if hasKanji := len(word.Kanji) > 0; hasKanji {
			processKanjiWord(&entry, word)
		} else {
			processKanaOnlyWord(&entry, word)
		}

		processSenseWord(&entry, word)

		entries = append(entries, entry)
	}

	return entries
}

func processKanjiWord(entry *UIEntry, word model.JMdictWord) {
	for _, kana := range word.Kana {
		if util.ContainsString(kana.Tags, "sK") {
			continue
		}

		for _, kanji := range word.Kanji {
			if util.ContainsString(kanji.Tags, "sK") {
				continue
			}

			for _, kanjiApplied := range kana.AppliesToKanji {
				if kanjiApplied == kanji.Text || kanjiApplied == "*" {
					furigana := Furigana{
						Word:    kanji.Text,
						Reading: kana.Text,
					}

					if entry.MainWord.Word == "" {
						entry.Common = word.Kanji[0].Common
						entry.MainWord = furigana
					} else {
						entry.OtherForms = append(entry.OtherForms, furigana)
					}
				}
			}
		}
	}
}

func processKanaOnlyWord(entry *UIEntry, word model.JMdictWord) {
	for i, kana := range word.Kana {
		furigana := Furigana{
			Word:    kana.Text,
			Reading: "",
		}

		if i == 0 {
			entry.Common = kana.Common
			entry.MainWord = furigana
		} else {
			entry.OtherForms = append(entry.OtherForms, furigana)
		}
	}
}

func processSenseWord(entry *UIEntry, word model.JMdictWord) {
	for _, sense := range word.Sense {
		var glossList []string
		for _, g := range sense.Gloss {
			glossList = append(glossList, g.Text)
		}

		entry.Meanings = append(entry.Meanings, strings.Join(glossList, ", "))
	}
}
