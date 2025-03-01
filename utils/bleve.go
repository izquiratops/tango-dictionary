package utils

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

func newTermQueryWithBoost(field string, term string, boost float64) *query.TermQuery {
	query := bleve.NewTermQuery(term)
	query.SetField(field)
	query.SetBoost(boost)
	return query
}

func newMatchQueryWithBoost(field string, term string, boost float64) *query.MatchQuery {
	query := bleve.NewMatchQuery(term)
	query.SetField(field)
	query.SetBoost(boost)
	return query
}

func NewJapaneseFieldQuery(query string, exactField string, charField string) *query.BooleanQuery {
	booleanQuery := bleve.NewBooleanQuery()

	exactQuery := newTermQueryWithBoost(exactField, query, 10.0)
	booleanQuery.AddShould(exactQuery)

	charQuery := newMatchQueryWithBoost(charField, query, 1.0)
	disjunctionQuery := bleve.NewDisjunctionQuery(charQuery)
	disjunctionQuery.SetMin(1)
	booleanQuery.AddShould(disjunctionQuery)

	return booleanQuery
}
