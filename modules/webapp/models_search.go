package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/query"
)

const textFieldAnalyzer = "standard"

var BleveIndex bleve.Index

func buildIndexMapping() *mapping.IndexMappingImpl {
	nameTextFieldMapping := bleve.NewTextFieldMapping()
	nameTextFieldMapping.Analyzer = textFieldAnalyzer

	descriptionTextFieldMapping := bleve.NewTextFieldMapping()
	descriptionTextFieldMapping.Analyzer = textFieldAnalyzer

	categoryenTextFieldMapping := bleve.NewTextFieldMapping()
	categoryenTextFieldMapping.Analyzer = textFieldAnalyzer

	parentcategoryenTextFieldMapping := bleve.NewTextFieldMapping()
	parentcategoryenTextFieldMapping.Analyzer = textFieldAnalyzer

	itemMapping := bleve.NewDocumentMapping()
	itemMapping.AddFieldMappingsAt("Name", nameTextFieldMapping)
	itemMapping.AddFieldMappingsAt("Description", descriptionTextFieldMapping)
	itemMapping.AddFieldMappingsAt("CategoryEn", categoryenTextFieldMapping)
	itemMapping.AddFieldMappingsAt("ParentCategoryEn", parentcategoryenTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("item", itemMapping)
	indexMapping.DefaultAnalyzer = textFieldAnalyzer

	return indexMapping
}

func init() {
	var err error
	// First, try to open the index.
	BleveIndex, err = bleve.Open("./data/index.bleve")
	if err != nil {
		// If it doesn't exist, create it.
		BleveIndex, err = bleve.New("./data/index.bleve", buildIndexMapping())
		if err != nil {
			panic(err)
		}
	}
}

func printDocsFromSearchResults(searchResults *bleve.SearchResult) {
	for _, hit := range searchResults.Hits {
		rv := struct {
			ID     string                 `json:"id"`
			Fields map[string]interface{} `json:"fields"`
		}{
			ID:     hit.ID,
			Fields: hit.Fields, // Corrected: Use fields from the search result hit directly.
		}

		js, err := json.MarshalIndent(rv, "", "    ")
		if err != nil {
			fmt.Printf("Error marshalling JSON for doc %s: %v\n", hit.ID, err)
			continue
		}
		fmt.Printf("%s\n", js)
	}
}

func searchWithQuery(query query.Query) *bleve.SearchResult {
	search := bleve.NewSearchRequest(query)
	search.Size = 2000
	// Corrected: Request fields to be included in the search results
	search.Fields = []string{"*"} 

	searchResults, err := BleveIndex.Search(search)
	if err != nil {
		panic(err)
	}
	return searchResults
}

func SearchItems(text string) []string {
	query := bleve.NewMatchQuery(text)
	searchResults := searchWithQuery(query)

	// Corrected: Only perform the fuzzy search if the first one yields no results.
	if len(searchResults.Hits) < 1 {
		query = bleve.NewQueryStringQuery(text + "~2")
		searchResults = searchWithQuery(query)
	}

	ids := make([]string, len(searchResults.Hits))
	for i, hit := range searchResults.Hits {
		ids[i] = hit.ID
	}
	return ids
}

// Additional helper functions for a complete example.
// These are not part of the original code, but are useful for demonstrating the corrected functions.
type Item struct {
	Name             string `json:"Name"`
	Description      string `json:"Description"`
	CategoryEn       string `json:"CategoryEn"`
	ParentCategoryEn string `json:"ParentCategoryEn"`
}

func main() {
	// A basic example of how the corrected functions would be used.
	// Indexing some dummy data.
	items := []Item{
		{"Go Language", "A statically typed, compiled programming language designed at Google.", "Programming Languages", "Computer Science"},
		{"Python Programming", "An interpreted, high-level, general-purpose programming language.", "Programming Languages", "Computer Science"},
		{"Java Language", "A high-level, class-based, object-oriented programming language.", "Programming Languages", "Computer Science"},
	}

	for i, item := range items {
		_ = BleveIndex.Index(fmt.Sprintf("item-%d", i), item)
	}

	// Wait for the index to be ready for searching.
	time.Sleep(1 * time.Second)

	// Demonstrating the SearchItems function.
	fmt.Println("Searching for 'Go Language'")
	ids := SearchItems("Go Language")
	fmt.Printf("Found IDs: %v\n", ids)

	// Demonstrating the fuzzy search.
	fmt.Println("\nSearching for 'Pyton Prograaming' (with typos)")
	ids = SearchItems("Pyton Prograaming")
	fmt.Printf("Found IDs: %v\n", ids)

	// Retrieve and print documents from search results
	fmt.Println("\nPrinting documents for 'Go'")
	query := bleve.NewMatchQuery("Go")
	searchResults := searchWithQuery(query)
	printDocsFromSearchResults(searchResults)
}
