package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var errStatusNotOK = errors.New("status not ok")

func DoHTTPRequest(ctx context.Context, method, path string, headers, params map[string]string, body, result interface{}) error {
	b := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(b).Encode(body)
		if err != nil {
			return fmt.Errorf("failed to encode json: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, path, b)
	if err != nil {
		return fmt.Errorf("new request failure: %w", err)
	}

	values := req.URL.Query()
	for k, v := range params {
		values.Add(k, v)
	}
	req.URL.RawQuery = values.Encode()

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpclient request failure: %w", err)
	}
	defer resp.Body.Close()
	bo, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return errStatusNotOK
	}

	err = json.Unmarshal(bo, &result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return nil
}
