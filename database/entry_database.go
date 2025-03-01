package database

import (
	"strings"
	"tango/jmdict"
	"tango/utils"
)

func ToEntryDatabase(word *jmdict.JMdictWord) EntryDatabase {
	entry := EntryDatabase{
		ID: word.ID,
	}

	if hasKanji := len(word.Kanji) > 0; hasKanji {
		processKanjiWord(&entry, word)
	} else {
		processKanaOnlyWord(&entry, word)
	}

	processSenseWord(&entry, word)

	return entry
}

func processKanjiWord(entry *EntryDatabase, word *jmdict.JMdictWord) {
	for _, kana := range word.Kana {
		if utils.ContainsString(kana.Tags, "sK") {
			continue
		}

		for _, kanji := range word.Kanji {
			if utils.ContainsString(kanji.Tags, "sK") {
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

func processKanaOnlyWord(entry *EntryDatabase, word *jmdict.JMdictWord) {
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

func processSenseWord(entry *EntryDatabase, word *jmdict.JMdictWord) {
	for _, sense := range word.Sense {
		var glossList []string
		for _, g := range sense.Gloss {
			glossList = append(glossList, g.Text)
		}

		entry.Meanings = append(entry.Meanings, strings.Join(glossList, "; "))
	}
}
