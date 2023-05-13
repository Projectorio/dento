package response

type ExplainResponse struct {
	Keyword        string           `json:"keyword"`
	Explanation    string           `json:"explanation"`
	Disclaimer     string           `json:"disclaimer"`
	Example        string           `json:"example"`
	FurtherReading []FurtherReading `json:"furtherReading"`
}

type FurtherReading struct {
	Keyword     string `json:"keyword"`
	Description string `json:"description"`
	Link        string `json:"link"`
}
