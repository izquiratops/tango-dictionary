package model

/*
 * Implements the same interfaces as jmdict-simplified 💕
 * https://scriptin.github.io/jmdict-simplified/interfaces/JMdict.html
 */

type Language string

type DictionaryMetadata[L Language] struct {
	Version   string `json:"version" bson:"version"`
	Languages []L    `json:"languages" bson:"languages"`
	DictDate  string `json:"dictDate" bson:"dict_date"`
}

type JMdictDictionaryMetadata struct {
	DictionaryMetadata[Language]
	CommonOnly    bool              `json:"commonOnly" bson:"common_only"`
	DictRevisions []string          `json:"dictRevisions" bson:"dict_revisions"`
	Tags          map[string]string `json:"tags" bson:"tags"`
}

type JMdict struct {
	JMdictDictionaryMetadata
	Words []JMdictWord `json:"words" bson:"words"`
}

type JMdictWord struct {
	ID    string        `json:"id" bson:"_id"`
	Kanji []JMdictKanji `json:"kanji" bson:"kanji"`
	Kana  []JMdictKana  `json:"kana" bson:"kana"`
	Sense []JMdictSense `json:"sense" bson:"sense"`
}

type JMdictKanji struct {
	Common bool     `json:"common" bson:"common"`
	Text   string   `json:"text" bson:"text"`
	Tags   []string `json:"tags" bson:"tags"`
}

type JMdictKana struct {
	Common         bool     `json:"common" bson:"common"`
	Text           string   `json:"text" bson:"text"`
	Tags           []string `json:"tags" bson:"tags"`
	AppliesToKanji []string `json:"appliesToKanji" bson:"applies_to_kanji"`
}

type JMdictSense struct {
	PartOfSpeech   []string               `json:"partOfSpeech" bson:"part_of_speech"`
	AppliesToKanji []string               `json:"appliesToKanji" bson:"applies_to_kanji"`
	AppliesToKana  []string               `json:"appliesToKana" bson:"applies_to_kana"`
	Related        []Xref                 `json:"related" bson:"related"`
	Antonym        []Xref                 `json:"antonym" bson:"antonym"`
	Field          []string               `json:"field" bson:"field"`
	Dialect        []string               `json:"dialect" bson:"dialect"`
	Misc           []string               `json:"misc" bson:"misc"`
	Info           []string               `json:"info" bson:"info"`
	LanguageSource []JMdictLanguageSource `json:"languageSource" bson:"language_source"`
	Gloss          []JMdictGloss          `json:"gloss" bson:"gloss"`
}

type JMdictLanguageSource struct {
	Lang  Language `json:"lang" bson:"lang"`
	Full  bool     `json:"full" bson:"full"`
	Wasei bool     `json:"wasei" bson:"wasei"`
	Text  *string  `json:"text" bson:"text,omitempty"`
}

type JMdictGender string

const (
	Masculine JMdictGender = "masculine"
	Feminine  JMdictGender = "feminine"
	Neuter    JMdictGender = "neuter"
)

type JMdictGlossType string

const (
	Literal     JMdictGlossType = "literal"
	Figurative  JMdictGlossType = "figurative"
	Explanation JMdictGlossType = "explanation"
	Trademark   JMdictGlossType = "trademark"
)

type JMdictGloss struct {
	Lang   Language         `json:"lang" bson:"lang"`
	Gender *JMdictGender    `json:"gender,omitempty" bson:"gender,omitempty"`
	Type   *JMdictGlossType `json:"type,omitempty" bson:"type,omitempty"`
	Text   string           `json:"text" bson:"text"`
}
