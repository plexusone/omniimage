package omniimage

// Model constants for image generation.
const (
	// OpenAI GPT Image models (current)
	ModelGPTImage2 = "gpt-image-2"
	ModelGPTImage1 = "gpt-image-1"

	// OpenAI DALL-E models (legacy, dall-e-3 retired March 2026)
	ModelDALLE3 = "dall-e-3"
	ModelDALLE2 = "dall-e-2"

	// Fal AI FLUX models
	ModelFluxPro     = "fal-ai/flux-pro"
	ModelFluxDev     = "fal-ai/flux/dev"
	ModelFluxSchnell = "fal-ai/flux/schnell"

	// Fal AI Stable Diffusion models
	ModelSDXL       = "fal-ai/fast-sdxl"
	ModelSD3        = "fal-ai/stable-diffusion-v3-medium"
	ModelSDXLTurbo  = "fal-ai/fast-turbo-diffusion"
	ModelPlayground = "fal-ai/playground-v25"

	// Fal AI upscaling models
	ModelClarityUpscaler = "fal-ai/clarity-upscaler"
	ModelCreativeUpscale = "fal-ai/creative-upscaler"
)

// ModelInfo contains information about a model.
type ModelInfo struct {
	ID           string
	Provider     ProviderName
	Name         string
	MaxImages    int
	SupportsHD   bool
	SupportsEdit bool
}

// GetModelInfo returns information about a model.
func GetModelInfo(modelID string) *ModelInfo {
	models := map[string]ModelInfo{
		ModelGPTImage2: {
			ID:           ModelGPTImage2,
			Provider:     ProviderNameOpenAI,
			Name:         "GPT Image 2",
			MaxImages:    10,
			SupportsHD:   true,
			SupportsEdit: false,
		},
		ModelGPTImage1: {
			ID:           ModelGPTImage1,
			Provider:     ProviderNameOpenAI,
			Name:         "GPT Image 1",
			MaxImages:    10,
			SupportsHD:   true,
			SupportsEdit: false,
		},
		ModelDALLE3: {
			ID:           ModelDALLE3,
			Provider:     ProviderNameOpenAI,
			Name:         "DALL-E 3 (legacy)",
			MaxImages:    1,
			SupportsHD:   true,
			SupportsEdit: false,
		},
		ModelDALLE2: {
			ID:           ModelDALLE2,
			Provider:     ProviderNameOpenAI,
			Name:         "DALL-E 2 (legacy)",
			MaxImages:    10,
			SupportsHD:   false,
			SupportsEdit: true,
		},
		ModelFluxPro: {
			ID:           ModelFluxPro,
			Provider:     ProviderNameFal,
			Name:         "FLUX Pro",
			MaxImages:    4,
			SupportsHD:   false,
			SupportsEdit: false,
		},
		ModelFluxDev: {
			ID:           ModelFluxDev,
			Provider:     ProviderNameFal,
			Name:         "FLUX Dev",
			MaxImages:    4,
			SupportsHD:   false,
			SupportsEdit: false,
		},
		ModelFluxSchnell: {
			ID:           ModelFluxSchnell,
			Provider:     ProviderNameFal,
			Name:         "FLUX Schnell",
			MaxImages:    4,
			SupportsHD:   false,
			SupportsEdit: false,
		},
		ModelSDXL: {
			ID:           ModelSDXL,
			Provider:     ProviderNameFal,
			Name:         "Stable Diffusion XL",
			MaxImages:    4,
			SupportsHD:   false,
			SupportsEdit: false,
		},
	}

	if info, ok := models[modelID]; ok {
		return &info
	}
	return nil
}
