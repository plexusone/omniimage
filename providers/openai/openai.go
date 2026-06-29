// Package openai implements the image generation provider for OpenAI's image models.
// Supports GPT Image models (gpt-image-2, gpt-image-1) and legacy DALL-E models.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/plexusone/omniimage/provider"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	providerName   = "openai"
)

// Config configures the OpenAI provider.
type Config struct {
	// APIKey is the OpenAI API key.
	// If empty, uses OPENAI_API_KEY environment variable.
	APIKey string

	// BaseURL is the API base URL.
	// If empty, uses the default OpenAI API URL.
	BaseURL string

	// HTTPClient is an optional custom HTTP client.
	HTTPClient *http.Client

	// Timeout is the request timeout.
	Timeout time.Duration
}

// Provider implements image generation using OpenAI's DALL-E API.
type Provider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// New creates a new OpenAI provider.
func New(cfg Config) (*Provider, error) {
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("openai API key is required")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = 120 * time.Second
		}
		httpClient = &http.Client{Timeout: timeout}
	}

	return &Provider{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return providerName
}

// Close closes the provider.
func (p *Provider) Close() error {
	return nil
}

// Generate creates images from a prompt using DALL-E.
func (p *Provider) Generate(ctx context.Context, req *provider.GenerateRequest) (*provider.GenerateResponse, error) {
	// Build OpenAI request
	openaiReq := map[string]any{
		"model":  req.Model,
		"prompt": req.Prompt,
	}

	if req.N > 0 {
		openaiReq["n"] = req.N
	}
	if req.Size != "" {
		openaiReq["size"] = string(req.Size)
	}
	if req.Quality != "" {
		openaiReq["quality"] = string(req.Quality)
	}
	if req.Style != "" {
		openaiReq["style"] = string(req.Style)
	}
	if req.ResponseFormat != "" {
		openaiReq["response_format"] = string(req.ResponseFormat)
	}
	if req.User != "" {
		openaiReq["user"] = req.User
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		p.baseURL+"/images/generations", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	//nolint:errcheck // Response body close errors are safe to ignore after successful read
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.parseError(resp.StatusCode, respBody)
	}

	var openaiResp struct {
		Created int64 `json:"created"`
		Data    []struct {
			URL           string `json:"url,omitempty"`
			B64JSON       string `json:"b64_json,omitempty"`
			RevisedPrompt string `json:"revised_prompt,omitempty"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	images := make([]provider.Image, len(openaiResp.Data))
	for i, d := range openaiResp.Data {
		images[i] = provider.Image{
			URL:           d.URL,
			B64JSON:       d.B64JSON,
			RevisedPrompt: d.RevisedPrompt,
			ContentType:   "image/png",
		}
	}

	return &provider.GenerateResponse{
		Created: time.Unix(openaiResp.Created, 0),
		Images:  images,
		Model:   req.Model,
	}, nil
}

// Edit modifies an existing image using DALL-E 2.
// Note: The OpenAI edit API traditionally requires multipart form data with image files.
// This implementation sends JSON with base64-encoded images, which may require
// the image to be provided as a data URL (data:image/png;base64,...).
func (p *Provider) Edit(ctx context.Context, req *provider.EditRequest) (*provider.EditResponse, error) {
	openaiReq := map[string]any{
		"model":  req.Model,
		"image":  req.Image,
		"prompt": req.Prompt,
	}

	if req.Mask != "" {
		openaiReq["mask"] = req.Mask
	}
	if req.N > 0 {
		openaiReq["n"] = req.N
	}
	if req.Size != "" {
		openaiReq["size"] = string(req.Size)
	}
	if req.ResponseFormat != "" {
		openaiReq["response_format"] = string(req.ResponseFormat)
	}
	if req.User != "" {
		openaiReq["user"] = req.User
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		p.baseURL+"/images/edits", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	//nolint:errcheck // Response body close errors are safe to ignore after successful read
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.parseError(resp.StatusCode, respBody)
	}

	var openaiResp struct {
		Created int64 `json:"created"`
		Data    []struct {
			URL     string `json:"url,omitempty"`
			B64JSON string `json:"b64_json,omitempty"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	images := make([]provider.Image, len(openaiResp.Data))
	for i, d := range openaiResp.Data {
		images[i] = provider.Image{
			URL:         d.URL,
			B64JSON:     d.B64JSON,
			ContentType: "image/png",
		}
	}

	return &provider.EditResponse{
		Created: time.Unix(openaiResp.Created, 0),
		Images:  images,
	}, nil
}

// Variations creates variations of an image using DALL-E 2.
func (p *Provider) Variations(ctx context.Context, req *provider.VariationsRequest) (*provider.VariationsResponse, error) {
	openaiReq := map[string]any{
		"model": req.Model,
		"image": req.Image,
	}

	if req.N > 0 {
		openaiReq["n"] = req.N
	}
	if req.Size != "" {
		openaiReq["size"] = string(req.Size)
	}
	if req.ResponseFormat != "" {
		openaiReq["response_format"] = string(req.ResponseFormat)
	}
	if req.User != "" {
		openaiReq["user"] = req.User
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		p.baseURL+"/images/variations", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	//nolint:errcheck // Response body close errors are safe to ignore after successful read
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.parseError(resp.StatusCode, respBody)
	}

	var openaiResp struct {
		Created int64 `json:"created"`
		Data    []struct {
			URL     string `json:"url,omitempty"`
			B64JSON string `json:"b64_json,omitempty"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	images := make([]provider.Image, len(openaiResp.Data))
	for i, d := range openaiResp.Data {
		images[i] = provider.Image{
			URL:         d.URL,
			B64JSON:     d.B64JSON,
			ContentType: "image/png",
		}
	}

	return &provider.VariationsResponse{
		Created: time.Unix(openaiResp.Created, 0),
		Images:  images,
	}, nil
}

// parseError parses an OpenAI API error response.
func (p *Provider) parseError(statusCode int, body []byte) error {
	var errResp struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("API error (status %d): %s", statusCode, string(body))
	}

	return &APIError{
		StatusCode: statusCode,
		Code:       errResp.Error.Code,
		Message:    errResp.Error.Message,
		Provider:   providerName,
	}
}

// APIError represents an OpenAI API error.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Provider   string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s: %s", e.Provider, e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Provider, e.Message)
}

// Ensure Provider implements the required interfaces.
var (
	_ provider.Provider           = (*Provider)(nil)
	_ provider.EditProvider       = (*Provider)(nil)
	_ provider.VariationsProvider = (*Provider)(nil)
)
