// Package provider defines the core interfaces that image generation providers must implement.
// External provider packages should import this package to implement the Provider interface.
package provider

import "context"

// Provider defines the interface that all image generation providers must implement.
// External packages can implement this interface and inject via omniimage.ClientConfig.CustomProvider.
type Provider interface {
	// Generate creates images from a text prompt.
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

	// Close closes the provider and cleans up resources.
	Close() error

	// Name returns the provider name.
	Name() string
}

// EditProvider extends Provider with image editing capabilities.
// Not all providers support editing (e.g., DALL-E supports it, Fal does not).
type EditProvider interface {
	Provider

	// Edit modifies an existing image based on a prompt and optional mask.
	Edit(ctx context.Context, req *EditRequest) (*EditResponse, error)
}

// VariationsProvider extends Provider with image variation capabilities.
// Not all providers support variations (e.g., DALL-E 2 supports it).
type VariationsProvider interface {
	Provider

	// Variations creates variations of an existing image.
	Variations(ctx context.Context, req *VariationsRequest) (*VariationsResponse, error)
}

// UpscaleProvider extends Provider with image upscaling capabilities.
type UpscaleProvider interface {
	Provider

	// Upscale increases the resolution of an image.
	Upscale(ctx context.Context, req *UpscaleRequest) (*UpscaleResponse, error)
}
