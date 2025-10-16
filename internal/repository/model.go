package repository

import "encoding/json"

type Issue struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	MarkdownBody string `json:"markdownBody"`
}

func (i *Issue) toJson() ([]byte, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
