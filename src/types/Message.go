package types

type Message struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}
