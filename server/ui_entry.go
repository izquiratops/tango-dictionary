package server

import "tango/model"

// TODO: Add a test for ProcessEntries
// UIEntry represents a simplified and processed version of a JMdictWord.
// Its used for populating HTML templates.
type UIEntry struct {
	MainWord   Furigana   `json:"mainWord"`
	OtherForms []Furigana `json:"otherForms"`
	IsCommon   bool       `json:"isCommon"`
	Meanings   []string   `json:"meanings"`
}

type Furigana struct {
	Word    string `json:"word"` // Kanji or kana as fallback
	Reading string `json:"reading"`
}

func ProcessEntries(words []model.JMdictWord) []UIEntry {
	processed := make([]UIEntry, 0, len(words))

	for _, word := range words {
		entry := UIEntry{
			OtherForms: make([]Furigana, 0),
			Meanings:   make([]string, 0),
		}

		// Handle main word and other forms
		if len(word.Kanji) > 0 {
			for i := 0; i < len(word.Kanji); i++ {
				reading := findMatchingKana(word.Kanji[i], word.Kana)
				entry.OtherForms = append(entry.OtherForms, Furigana{
					Word:    word.Kanji[i].Text,
					Reading: reading,
				})

				if i == 0 {
					entry.IsCommon = word.Kanji[0].Common
				}
			}
		} else if len(word.Kana) > 0 {
			// If no kanji, use first kana as main word
			for i := 0; i < len(word.Kana); i++ {
				entry.OtherForms = append(entry.OtherForms, Furigana{
					Word:    word.Kana[i].Text,
					Reading: "",
				})

				if i == 0 {
					entry.IsCommon = word.Kana[0].Common
				}
			}
		}

		// Process meanings from sense
		for _, sense := range word.Sense {
			// Check if this sense applies to the main word
			if isValidSense(entry.MainWord.Word, sense, word.Kanji) {
				for _, gloss := range sense.Gloss {
					if gloss.Lang == "eng" { // Assuming we want English meanings
						entry.Meanings = append(entry.Meanings, gloss.Text)
					}
				}
			}
		}

		processed = append(processed, entry)
	}

	return processed
}

// Returns the appropriate reading for a kanji based on AppliesToKanji
func findMatchingKana(kanji model.JMdictKanji, kanaList []model.JMdictKana) string {
	for _, kana := range kanaList {
		// TODO: Include also '*'
		// If appliesToKanji is empty, this kana applies to all kanji
		if len(kana.AppliesToKanji) == 0 {
			return kana.Text
		}
		// Check if this kana explicitly applies to this kanji
		for _, applies := range kana.AppliesToKanji {
			if applies == kanji.Text {
				return kana.Text
			}
		}
	}

	// If no match found, return first kana as fallback
	if len(kanaList) > 0 {
		return kanaList[0].Text
	}

	return ""
}

// Checks if a sense applies to the given word
func isValidSense(word string, sense model.JMdictSense, kanjiList []model.JMdictKanji) bool {
	// If both applies lists are empty, the sense applies to all forms
	if len(sense.AppliesToKanji) == 0 && len(sense.AppliesToKana) == 0 {
		return true
	}

	// Check if word is in kanji list
	isKanji := false
	for _, k := range kanjiList {
		if k.Text == word {
			isKanji = true
			break
		}
	}

	if isKanji {
		// Check if sense applies to this kanji
		if len(sense.AppliesToKanji) == 0 {
			return true
		}
		for _, applies := range sense.AppliesToKanji {
			if applies == word {
				return true
			}
		}
	} else {
		// Check if sense applies to this kana
		if len(sense.AppliesToKana) == 0 {
			return true
		}
		for _, applies := range sense.AppliesToKana {
			if applies == word {
				return true
			}
		}
	}

	return false
}
