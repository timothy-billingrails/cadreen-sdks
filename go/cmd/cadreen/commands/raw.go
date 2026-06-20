package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ListToolsResponse = ListToolsResponseData

type ListToolsResponseData struct {
	Object string      `json:"object"`
	Data   []ToolEntry `json:"data"`
}

type ToolEntry struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

func rawDo(method, path string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	baseURL := cfg.BaseURLResolved()
	if !strings.HasPrefix(baseURL, "https://") && !strings.HasPrefix(baseURL, "http://localhost:") && !strings.HasPrefix(baseURL, "http://127.0.0.1:") {
		return nil, fmt.Errorf("base URL must use HTTPS (got %s). Set CADREEN_BASE_URL or run 'cadreen config set base_url <url>'", baseURL)
	}

	url := baseURL + path
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKeyResolved())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return resp, nil
}
