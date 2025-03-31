package main

import (
	"fmt"
	"strings"

	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/jmdict"
	"github.com/izquiratops/tango/common/utils"
)

func ToWord(word *jmdict.JMdictWord) database.Word {
	entry := database.Word{
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

func ToWordSearchable(d *jmdict.JMdictWord) (database.WordSearchable, error) {
	entry := database.WordSearchable{
		ID:         d.ID,
		KanjiExact: make([]string, 0),
		KanjiChar:  make([]string, 0),
		KanaExact:  make([]string, 0),
		KanaChar:   make([]string, 0),
		Meanings:   make([]string, 0),
	}

	for _, k := range d.Kanji {
		if k.Text == "" {
			return entry, fmt.Errorf("emtpy field at %v", d.ID)
		}

		entry.KanjiExact = append(entry.KanjiExact, k.Text)
		entry.KanjiChar = append(entry.KanjiChar, k.Text)
	}

	for _, k := range d.Kana {
		if k.Text == "" {
			return entry, fmt.Errorf("emtpy field at %v", d.ID)
		}

		entry.KanaExact = append(entry.KanaExact, k.Text)
		entry.KanaChar = append(entry.KanaChar, k.Text)
	}

	for _, s := range d.Sense {
		for _, g := range s.Gloss {
			if g.Lang == "eng" {
				if g.Text == "" {
					return entry, fmt.Errorf("emtpy field at %v", d.ID)
				}

				entry.Meanings = append(entry.Meanings, g.Text)
			}
		}
	}

	return entry, nil
}

func processKanjiWord(entry *database.Word, word *jmdict.JMdictWord) {
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
					furigana := database.Furigana{
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

func processKanaOnlyWord(entry *database.Word, word *jmdict.JMdictWord) {
	for i, kana := range word.Kana {
		furigana := database.Furigana{
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

func processSenseWord(entry *database.Word, word *jmdict.JMdictWord) {
	for _, sense := range word.Sense {
		var glossList []string
		for _, g := range sense.Gloss {
			glossList = append(glossList, g.Text)
		}

		entry.Meanings = append(entry.Meanings, strings.Join(glossList, "; "))
	}
}
