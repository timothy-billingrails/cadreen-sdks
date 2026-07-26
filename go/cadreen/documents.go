package cadreen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func (c *Client) ListDocuments(ctx context.Context, opts ...RequestOption) (*ListDocumentsResponse, error) {
	var result ListDocumentsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/documents", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	return &result, nil
}

func (c *Client) GetDocument(ctx context.Context, id string, opts ...RequestOption) (*Document, error) {
	var result Document
	if err := c.do(ctx, "GET", "/api/v1/cadreen/documents/"+url.PathEscape(id), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}
	return &result, nil
}

func (c *Client) UploadDocument(ctx context.Context, filePath string) (*UploadDocumentResponse, error) {
	if c.config.Sandbox {
		var result UploadDocumentResponse
		if v, ok := c.config.Fixtures["POST /api/v1/cadreen/documents/upload"]; ok {
			if err := populateResult(v, &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
		return nil, &APIError{
			StatusCode: 404,
			Type:       "not_found",
			Code:       "not_found",
			Message:    "No fixture for POST /api/v1/cadreen/documents/upload",
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("document", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("copy file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close writer: %w", err)
	}

	reqURL := c.config.BaseURL + "/api/v1/cadreen/documents/upload"
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, &body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Accept", fmt.Sprintf(`application/json; profile="%s"`, c.config.Profile))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	c.setLastTrace(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, parseError(resp.StatusCode, respBody)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result UploadDocumentResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &result, nil
}
