// Package fal implements the image generation provider for Fal AI.
// Fal AI provides access to FLUX, Stable Diffusion, and other models.
package fal

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
	defaultBaseURL = "https://fal.run"
	providerName   = "fal"
)

// Config configures the Fal AI provider.
type Config struct {
	// APIKey is the Fal AI API key.
	// If empty, uses FAL_KEY environment variable.
	APIKey string

	// BaseURL is the API base URL.
	// If empty, uses the default Fal AI URL.
	BaseURL string

	// HTTPClient is an optional custom HTTP client.
	HTTPClient *http.Client

	// Timeout is the request timeout.
	Timeout time.Duration
}

// Provider implements image generation using Fal AI.
type Provider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// New creates a new Fal AI provider.
func New(cfg Config) (*Provider, error) {
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("FAL_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("fal AI API key is required")
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

// Generate creates images from a prompt using Fal AI models.
func (p *Provider) Generate(ctx context.Context, req *provider.GenerateRequest) (*provider.GenerateResponse, error) {
	// Build Fal AI request
	falReq := map[string]any{
		"prompt": req.Prompt,
	}

	if req.NegativePrompt != "" {
		falReq["negative_prompt"] = req.NegativePrompt
	}
	if req.N > 0 {
		falReq["num_images"] = req.N
	}
	if req.Size != "" {
		// Parse size into width/height
		width, height := parseSize(req.Size)
		falReq["image_size"] = map[string]int{
			"width":  width,
			"height": height,
		}
	}
	if req.Steps != nil {
		falReq["num_inference_steps"] = *req.Steps
	}
	if req.GuidanceScale != nil {
		falReq["guidance_scale"] = *req.GuidanceScale
	}
	if req.Seed != nil {
		falReq["seed"] = *req.Seed
	}

	// Add provider-specific options
	for k, v := range req.ProviderOptions {
		falReq[k] = v
	}

	body, err := json.Marshal(falReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Fal AI uses model ID as the endpoint path
	endpoint := p.baseURL + "/" + req.Model
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Key "+p.apiKey)

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

	var falResp struct {
		Images []struct {
			URL         string `json:"url"`
			Width       int    `json:"width"`
			Height      int    `json:"height"`
			ContentType string `json:"content_type"`
		} `json:"images"`
		Seed    int64 `json:"seed"`
		Timings struct {
			Inference float64 `json:"inference"`
		} `json:"timings"`
	}

	if err := json.Unmarshal(respBody, &falResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	images := make([]provider.Image, len(falResp.Images))
	for i, img := range falResp.Images {
		seed := falResp.Seed
		images[i] = provider.Image{
			URL:         img.URL,
			ContentType: img.ContentType,
			Width:       img.Width,
			Height:      img.Height,
			Seed:        &seed,
		}
	}

	return &provider.GenerateResponse{
		Created: time.Now(),
		Images:  images,
		Model:   req.Model,
		ProviderMetadata: map[string]any{
			"seed":           falResp.Seed,
			"inference_time": falResp.Timings.Inference,
		},
	}, nil
}

// Upscale increases the resolution of an image using Fal AI upscalers.
func (p *Provider) Upscale(ctx context.Context, req *provider.UpscaleRequest) (*provider.UpscaleResponse, error) {
	falReq := map[string]any{
		"image_url": req.Image,
	}

	if req.Scale > 0 {
		falReq["scale"] = req.Scale
	}

	body, err := json.Marshal(falReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := p.baseURL + "/" + req.Model
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Key "+p.apiKey)

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

	var falResp struct {
		Image struct {
			URL         string `json:"url"`
			Width       int    `json:"width"`
			Height      int    `json:"height"`
			ContentType string `json:"content_type"`
		} `json:"image"`
	}

	if err := json.Unmarshal(respBody, &falResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &provider.UpscaleResponse{
		Created: time.Now(),
		Image: provider.Image{
			URL:         falResp.Image.URL,
			ContentType: falResp.Image.ContentType,
			Width:       falResp.Image.Width,
			Height:      falResp.Image.Height,
		},
	}, nil
}

// parseSize converts an ImageSize to width and height.
func parseSize(size provider.ImageSize) (width, height int) {
	switch size {
	case provider.Size256x256:
		return 256, 256
	case provider.Size512x512:
		return 512, 512
	case provider.Size768x1024:
		return 768, 1024
	case provider.Size1024x768:
		return 1024, 768
	case provider.Size1024x1024:
		return 1024, 1024
	case provider.Size1024x1792:
		return 1024, 1792
	case provider.Size1792x1024:
		return 1792, 1024
	default:
		return 1024, 1024
	}
}

// parseError parses a Fal AI API error response.
func (p *Provider) parseError(statusCode int, body []byte) error {
	var errResp struct {
		Detail string `json:"detail"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("API error (status %d): %s", statusCode, string(body))
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    errResp.Detail,
		Provider:   providerName,
	}
}

// APIError represents a Fal AI API error.
type APIError struct {
	StatusCode int
	Message    string
	Provider   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Provider, e.Message)
}

// Ensure Provider implements the required interfaces.
var (
	_ provider.Provider        = (*Provider)(nil)
	_ provider.UpscaleProvider = (*Provider)(nil)
)
