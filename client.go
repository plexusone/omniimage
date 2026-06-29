// Package omniimage provides a unified interface for image generation across multiple providers.
//
// Supported providers:
//   - OpenAI (DALL-E 2, DALL-E 3)
//   - Fal AI (FLUX, Stable Diffusion)
//
// Example usage:
//
//	client, err := omniimage.NewClient(omniimage.ClientConfig{
//	    Providers: []omniimage.ProviderConfig{
//	        {Provider: omniimage.ProviderNameOpenAI, APIKey: "sk-..."},
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	resp, err := client.Generate(ctx, &provider.GenerateRequest{
//	    Model:  "dall-e-3",
//	    Prompt: "A serene mountain landscape at sunset",
//	    Size:   provider.Size1024x1024,
//	})
package omniimage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/plexusone/omniimage/provider"
)

// Client is the main client for image generation.
type Client struct {
	provider provider.Provider
	logger   *slog.Logger
}

// ClientConfig holds configuration for creating a client.
type ClientConfig struct {
	// Providers is a list of provider configurations.
	// Index 0 is the primary provider.
	Providers []ProviderConfig

	// Logger for internal logging (optional).
	Logger *slog.Logger
}

// NewClient creates a new image generation client.
func NewClient(config ClientConfig) (*Client, error) {
	if len(config.Providers) == 0 {
		return nil, ErrNoProviders
	}

	// Build the primary provider
	prov, err := buildProviderFromConfig(config.Providers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to create provider (%s): %w",
			config.Providers[0].Provider, err)
	}

	logger := config.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	return &Client{
		provider: prov,
		logger:   logger,
	}, nil
}

// Generate creates images from a text prompt.
func (c *Client) Generate(ctx context.Context, req *provider.GenerateRequest) (*provider.GenerateResponse, error) {
	return c.provider.Generate(ctx, req)
}

// Edit modifies an existing image based on a prompt.
// Returns ErrNotSupported if the provider doesn't support editing.
func (c *Client) Edit(ctx context.Context, req *provider.EditRequest) (*provider.EditResponse, error) {
	ep, ok := c.provider.(provider.EditProvider)
	if !ok {
		return nil, fmt.Errorf("%w: edit not supported by %s", ErrNotSupported, c.provider.Name())
	}
	return ep.Edit(ctx, req)
}

// Variations creates variations of an existing image.
// Returns ErrNotSupported if the provider doesn't support variations.
func (c *Client) Variations(ctx context.Context, req *provider.VariationsRequest) (*provider.VariationsResponse, error) {
	vp, ok := c.provider.(provider.VariationsProvider)
	if !ok {
		return nil, fmt.Errorf("%w: variations not supported by %s", ErrNotSupported, c.provider.Name())
	}
	return vp.Variations(ctx, req)
}

// Upscale increases the resolution of an image.
// Returns ErrNotSupported if the provider doesn't support upscaling.
func (c *Client) Upscale(ctx context.Context, req *provider.UpscaleRequest) (*provider.UpscaleResponse, error) {
	up, ok := c.provider.(provider.UpscaleProvider)
	if !ok {
		return nil, fmt.Errorf("%w: upscale not supported by %s", ErrNotSupported, c.provider.Name())
	}
	return up.Upscale(ctx, req)
}

// SupportsEdit returns true if the provider supports image editing.
func (c *Client) SupportsEdit() bool {
	_, ok := c.provider.(provider.EditProvider)
	return ok
}

// SupportsVariations returns true if the provider supports variations.
func (c *Client) SupportsVariations() bool {
	_, ok := c.provider.(provider.VariationsProvider)
	return ok
}

// SupportsUpscale returns true if the provider supports upscaling.
func (c *Client) SupportsUpscale() bool {
	_, ok := c.provider.(provider.UpscaleProvider)
	return ok
}

// Provider returns the underlying provider.
func (c *Client) Provider() provider.Provider {
	return c.provider
}

// Close closes the client and releases resources.
func (c *Client) Close() error {
	return c.provider.Close()
}
