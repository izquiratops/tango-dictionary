/*
 * Implements the same interfaces as jmdict-simplified 💕
 * https://scriptin.github.io/jmdict-simplified/interfaces/JMdict.html
 */

package main

import (
	"github.com/blevesearch/bleve/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// JMDICT

type Language string

type DictionaryMetadata[L Language] struct {
	Version   string `json:"version"`
	Languages []L    `json:"languages"`
	DictDate  string `json:"dictDate"`
}

type JMdictDictionaryMetadata struct {
	DictionaryMetadata[Language]
	CommonOnly    bool              `json:"commonOnly"`
	DictRevisions []string          `json:"dictRevisions"`
	Tags          map[string]string `json:"tags"`
}

type JMdict struct {
	JMdictDictionaryMetadata
	Words []JMdictWord `json:"words"`
}

type JMdictWord struct {
	ID    string        `json:"id"`
	Kanji []JMdictKanji `json:"kanji"`
	Kana  []JMdictKana  `json:"kana"`
	Sense []JMdictSense `json:"sense"`
}

type JMdictKanji struct {
	Common bool     `json:"common"`
	Text   string   `json:"text"`
	Tags   []string `json:"tags"`
}

type JMdictKana struct {
	Common         bool     `json:"common"`
	Text           string   `json:"text"`
	Tags           []string `json:"tags"`
	AppliesToKanji []string `json:"appliesToKanji"`
}

type JMdictSense struct {
	PartOfSpeech   []string               `json:"partOfSpeech"`
	AppliesToKanji []string               `json:"appliesToKanji"`
	AppliesToKana  []string               `json:"appliesToKana"`
	Related        []Xref                 `json:"related"`
	Antonym        []Xref                 `json:"antonym"`
	Field          []string               `json:"field"`
	Dialect        []string               `json:"dialect"`
	Misc           []string               `json:"misc"`
	Info           []string               `json:"info"`
	LanguageSource []JMdictLanguageSource `json:"languageSource"`
	Gloss          []JMdictGloss          `json:"gloss"`
}

type Xref struct {
	Kanji      string `json:"kanji,omitempty"`
	Kana       string `json:"kana,omitempty"`
	SenseIndex int    `json:"senseIndex,omitempty"`
}

type JMdictLanguageSource struct {
	Lang  Language `json:"lang"`
	Full  bool     `json:"full"`
	Wasei bool     `json:"wasei"`
	Text  *string  `json:"text"`
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
	Lang   Language         `json:"lang"`
	Gender *JMdictGender    `json:"gender,omitempty"`
	Type   *JMdictGlossType `json:"type,omitempty"`
	Text   string           `json:"text"`
}

// Used in JishoClone

type SearchableEntry struct {
	ID       string   `json:"id"`
	Kanji    []string `json:"kanji"`
	Kana     []string `json:"kana"`
	Meanings []string `json:"meanings"`
}

type DictionaryImporter struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
	bleveIndex  bleve.Index
	batchSize   int
}
