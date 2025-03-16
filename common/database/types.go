package database

type WordSearchable struct {
	ID         string   `json:"id"`
	KanjiExact []string `json:"kanji_exact"`
	KanjiChar  []string `json:"kanji_char"`
	KanaExact  []string `json:"kana_exact"`
	KanaChar   []string `json:"kana_char"`
	Meanings   []string `json:"meanings"`
	Romaji     []string `json:"romaji"`
}

type Word struct {
	ID         string     `json:"id" bson:"_id"`
	MainWord   Furigana   `json:"mainWord" bson:"main_word"`     // Primary word representation
	OtherForms []Furigana `json:"otherForms" bson:"other_forms"` // Alternative forms of the word
	Common     bool       `json:"isCommon" bson:"is_common"`     // Indicates if word is frequently used
	Meanings   []string   `json:"meanings" bson:"meanings"`      // Word definitions/translations
}

type Furigana struct {
	Word    string `json:"word" bson:"word"`       // Kanji representation (or kana if no kanji exists)
	Reading string `json:"reading" bson:"reading"` // Kana reading (empty for kana-only words)
}
