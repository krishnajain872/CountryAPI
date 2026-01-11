package types

type Country struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
	Population int64  `json:"population"`
}

type SearchRequest struct {
	Name string `json:"name"`
}
