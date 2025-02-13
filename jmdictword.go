package main

import "fmt"

func (d *JMdictWord) ToSearchable() (SearchableEntry, error) {
	searchable := SearchableEntry{
		ID:       d.ID,
		Kanji:    make([]string, 0),
		Kana:     make([]string, 0),
		Meanings: make([]string, 0),
	}

	for _, k := range d.Kanji {
		if k.Text == "" {
			return searchable, fmt.Errorf("emtpy field at %v", d.ID)
		}

		searchable.Kanji = append(searchable.Kanji, k.Text)
	}

	for _, k := range d.Kana {
		if k.Text == "" {
			return searchable, fmt.Errorf("emtpy field at %v", d.ID)
		}

		searchable.Kana = append(searchable.Kana, k.Text)
	}

	for _, s := range d.Sense {
		for _, g := range s.Gloss {
			if g.Lang == "eng" {
				if g.Text == "" {
					return searchable, fmt.Errorf("emtpy field at %v", d.ID)
				}

				searchable.Meanings = append(searchable.Meanings, g.Text)
			}
		}
	}

	return searchable, nil
}
