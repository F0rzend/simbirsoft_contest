package types

import "net/http"

type Response struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

func (r *Response) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
