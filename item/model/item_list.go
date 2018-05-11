package item

import (
	"strconv"

	"github.com/blevesearch/bleve"
)

// List searches for items which maches text in items and comments
func (ds *Datastore) List(queryText string, offset, limit int) ([]*Item, error) {
	matchQuery := bleve.NewQueryStringQuery(queryText)
	searchRequest := bleve.NewSearchRequest(matchQuery)
	bleve.NewConjunctionQuery(matchQuery)
	searchRequest.Fields = []string{"*"}
	searchRequest.From = offset
	searchRequest.Size = limit

	projectFacet := bleve.NewFacetRequest("project", 7)
	searchRequest.AddFacet("projects", projectFacet)

	labelFacet := bleve.NewFacetRequest("label", 7)
	searchRequest.AddFacet("labels", labelFacet)

	milestoneFacet := bleve.NewFacetRequest("milestone", 7)
	searchRequest.AddFacet("milestones", milestoneFacet)

	authorFacet := bleve.NewFacetRequest("author", 7)
	searchRequest.AddFacet("authors", authorFacet)

	editorFacet := bleve.NewFacetRequest("editor", 7)
	searchRequest.AddFacet("editors", editorFacet)

	openStateFacet := bleve.NewFacetRequest("state", 2)
	searchRequest.AddFacet("states", openStateFacet)

	mentionFacet := bleve.NewFacetRequest("mention", 7)
	searchRequest.AddFacet("mentions", mentionFacet)

	//searchRequest.SortBy([]string{"-_score"}) // Best match
	/*
		searchRequest.SortBy([]string{"-updated"}) // Recently updated
		searchRequest.SortBy([]string{"-_id"}) // Newest
		searchRequest.SortBy([]string{"_id"}) // Oldest
		searchRequest.SortBy([]string{"updated"}) // Least recently updated
	*/

	searchResults, err := ds.bi.Idx.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	var items []*Item
	for _, res := range searchResults.Hits {
		id, err := strconv.Atoi(res.ID)
		if err != nil {
			return nil, err
		}
		itm := Item{ID: id, Number: strconv.Itoa(int(res.Fields["num"].(float64))), ProjectID: "1", Title: res.Fields["title"].(string), Description: res.Fields["description"].(string)}
		i, ok := res.Fields["label"]
		if ok {
			for _, j := range i.([]interface{}) {
				itm.Labels = append(itm.Labels, j.(string))
			}
		}
		items = append(items, &itm)
	}
	return items, nil
}
