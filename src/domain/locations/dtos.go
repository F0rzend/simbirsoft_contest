package locations

import (
	"net/http"
)

type Response struct {
	ID        int64   `json:"id"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (*Response) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
