package omniimage

import (
	"fmt"

	"github.com/plexusone/omniimage/provider"
	"github.com/plexusone/omniimage/providers/fal"
	"github.com/plexusone/omniimage/providers/openai"
)

// ProviderName identifies an image generation provider.
type ProviderName string

const (
	// ProviderNameOpenAI is the OpenAI provider (DALL-E).
	ProviderNameOpenAI ProviderName = "openai"

	// ProviderNameFal is the Fal AI provider.
	ProviderNameFal ProviderName = "fal"
)

// ProviderConfig configures a provider.
type ProviderConfig struct {
	// Provider is the provider name.
	Provider ProviderName

	// APIKey is the API key for authentication.
	APIKey string

	// BaseURL is an optional custom base URL for the API.
	BaseURL string

	// CustomProvider allows injecting a custom provider implementation.
	CustomProvider provider.Provider
}

// buildProviderFromConfig creates a provider from configuration.
func buildProviderFromConfig(cfg ProviderConfig) (provider.Provider, error) {
	// Use custom provider if provided
	if cfg.CustomProvider != nil {
		return cfg.CustomProvider, nil
	}

	switch cfg.Provider {
	case ProviderNameOpenAI:
		return openai.New(openai.Config{
			APIKey:  cfg.APIKey,
			BaseURL: cfg.BaseURL,
		})

	case ProviderNameFal:
		return fal.New(fal.Config{
			APIKey:  cfg.APIKey,
			BaseURL: cfg.BaseURL,
		})

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownProvider, cfg.Provider)
	}
}

// Capabilities describes the features supported by a provider.
type Capabilities struct {
	// Generate indicates support for image generation.
	Generate bool

	// Edit indicates support for image editing.
	Edit bool

	// Variations indicates support for image variations.
	Variations bool

	// Upscale indicates support for image upscaling.
	Upscale bool

	// Models lists supported model IDs.
	Models []string

	// MaxImages is the maximum number of images per request.
	MaxImages int

	// SupportedSizes lists supported image sizes.
	SupportedSizes []provider.ImageSize
}

// GetCapabilities returns the capabilities for a provider.
func GetCapabilities(name ProviderName) *Capabilities {
	switch name {
	case ProviderNameOpenAI:
		return &Capabilities{
			Generate:   true,
			Edit:       true,
			Variations: true,
			Upscale:    false,
			Models:     []string{ModelDALLE3, ModelDALLE2},
			MaxImages:  10, // DALL-E 2; DALL-E 3 is limited to 1
			SupportedSizes: []provider.ImageSize{
				provider.Size256x256,
				provider.Size512x512,
				provider.Size1024x1024,
				provider.Size1024x1792,
				provider.Size1792x1024,
			},
		}

	case ProviderNameFal:
		return &Capabilities{
			Generate:   true,
			Edit:       false,
			Variations: false,
			Upscale:    true,
			Models:     []string{ModelFluxPro, ModelFluxDev, ModelFluxSchnell, ModelSDXL},
			MaxImages:  4,
			SupportedSizes: []provider.ImageSize{
				provider.Size512x512,
				provider.Size768x1024,
				provider.Size1024x768,
				provider.Size1024x1024,
			},
		}

	default:
		return nil
	}
}
