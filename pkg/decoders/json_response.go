package decoders

// Helper for decoding json bodies from responses

import (
	"encoding/json"
	"io"
	"net/http"
)

// ReadJSONResponse reads a json blob from a
// http response and decodes it into a map
func ReadJSONResponse(res *http.Response) (map[string]any, error) {
	// Read body
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	payload := make(map[string]any)
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}
