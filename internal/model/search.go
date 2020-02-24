package model

type SearchParameters struct {
	Searches []string
}

type SearchResult struct {
	Type    string        `json:"type"`
	Results []interface{} `json:"results"`
}
