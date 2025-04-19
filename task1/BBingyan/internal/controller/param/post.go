package param

type PostRequest struct {
	Title   string `json:"title"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}
