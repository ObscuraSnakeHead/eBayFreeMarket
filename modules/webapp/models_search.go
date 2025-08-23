package webapp

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
	for _, val := range searchResults.Hits {
		id := val.ID
		doc, err := BleveIndex.Document(id)
		if err != nil {
			fmt.Printf("Error retrieving document %s: %v\n", id, err)
			continue
		}

		rv := struct {
			ID     string                 `json:"id"`
			Fields map[string]interface{} `json:"fields"`
		}{
			ID:     id,
			Fields: map[string]interface{}{},
		}

		for _, field := range doc.Fields {
			var newval interface{}
			switch field := field.(type) {
			case *document.TextField:
				newval = string(field.Value())
			case *document.NumericField:
				n, err := field.Number()
				if err == nil {
					newval = n
				}
			case *document.DateTimeField:
				d, err := field.DateTime()
				if err == nil {
					newval = d.Format(time.RFC3339Nano)
				}
			case *document.StringField:
				newval = string(field.Value())
			default:
				newval = nil
			}

			if newval == nil {
				continue
			}

			existing, existed := rv.Fields[field.Name()]
			if existed {
				switch existing := existing.(type) {
				case []interface{}:
					rv.Fields[field.Name()] = append(existing, newval)
				case interface{}:
					arr := make([]interface{}, 2)
					arr[0] = existing
					arr[1] = newval
					rv.Fields[field.Name()] = arr
				}
			} else {
				rv.Fields[field.Name()] = newval
			}
		}

		js, err := json.MarshalIndent(rv, "", "    ")
		if err != nil {
			fmt.Printf("Error marshalling JSON for doc %s: %v\n", id, err)
			continue
		}
		fmt.Printf("%s\n", js)
	}
}

func searchWithQuery(query query.Query) *bleve.SearchResult {
	search := bleve.NewSearchRequest(query)
	search.Size = 2000

	searchResults, err := BleveIndex.Search(search)
	if err != nil {
		panic(err)
	}
	return searchResults
}

func SearchItems(text string) []string {
	query := bleve.NewMatchQuery(text)
	var searchResults *bleve.SearchResult
	searchResults = searchWithQuery(query)

	if len(searchResults.Hits) < 1 {
		query := bleve.NewQueryStringQuery(text + "~2")
		searchResults = searchWithQuery(query)
	}

	ids := []string{}
	for _, hit := range searchResults.Hits {
		ids = append(ids, hit.ID)
	}
	return ids
}
