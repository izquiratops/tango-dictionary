package database

import (
	"fmt"
	"tango/jmdict"
)

func ToEntrySearchable(d *jmdict.JMdictWord) (EntrySearchable, error) {
	entry := EntrySearchable{
		ID:         d.ID, // ID not indexed
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
