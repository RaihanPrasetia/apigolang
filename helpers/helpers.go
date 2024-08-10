package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ParseJSONRequestBody membaca dan mendekode body request JSON
func ParseJSONRequestBody(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	fmt.Printf("Request body: %s\n", string(body))

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return nil
}
