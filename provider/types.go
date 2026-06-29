package provider

import "time"

// ImageSize represents standard image dimensions.
type ImageSize string

const (
	// Square sizes
	Size256x256   ImageSize = "256x256"
	Size512x512   ImageSize = "512x512"
	Size1024x1024 ImageSize = "1024x1024"

	// Portrait sizes
	Size1024x1792 ImageSize = "1024x1792"
	Size768x1024  ImageSize = "768x1024"

	// Landscape sizes
	Size1792x1024 ImageSize = "1792x1024"
	Size1024x768  ImageSize = "1024x768"
)

// ImageQuality represents the quality level of generated images.
type ImageQuality string

const (
	QualityStandard ImageQuality = "standard"
	QualityHD       ImageQuality = "hd"
)

// ImageStyle represents the artistic style of generated images.
type ImageStyle string

const (
	StyleVivid   ImageStyle = "vivid"
	StyleNatural ImageStyle = "natural"
)

// ResponseFormat specifies how images are returned.
type ResponseFormat string

const (
	// FormatURL returns a URL to the generated image.
	FormatURL ResponseFormat = "url"
	// FormatB64JSON returns the image as base64-encoded JSON.
	FormatB64JSON ResponseFormat = "b64_json"
)

// GenerateRequest represents a request to generate images from a prompt.
type GenerateRequest struct {
	// Model is the model to use for generation.
	// Examples: "dall-e-3", "dall-e-2", "flux-pro", "flux-dev", "stable-diffusion-xl"
	Model string `json:"model"`

	// Prompt is the text description of the desired image(s).
	Prompt string `json:"prompt"`

	// NegativePrompt specifies what to avoid in the image (not all providers support this).
	NegativePrompt string `json:"negative_prompt,omitempty"`

	// N is the number of images to generate (default: 1).
	// DALL-E 3 only supports n=1.
	N int `json:"n,omitempty"`

	// Size is the dimensions of the generated image.
	Size ImageSize `json:"size,omitempty"`

	// Quality is the quality level (standard or hd).
	// Only supported by some providers (e.g., DALL-E 3).
	Quality ImageQuality `json:"quality,omitempty"`

	// Style affects the visual style of the image.
	// Only supported by some providers (e.g., DALL-E 3).
	Style ImageStyle `json:"style,omitempty"`

	// ResponseFormat specifies how the image is returned.
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`

	// User is an optional unique identifier for the end-user.
	User string `json:"user,omitempty"`

	// Seed for reproducible generation (not all providers support this).
	Seed *int64 `json:"seed,omitempty"`

	// Steps is the number of inference steps (for diffusion models).
	Steps *int `json:"steps,omitempty"`

	// GuidanceScale controls how closely the model follows the prompt.
	GuidanceScale *float64 `json:"guidance_scale,omitempty"`

	// ProviderOptions contains provider-specific options.
	ProviderOptions map[string]any `json:"provider_options,omitempty"`
}

// GenerateResponse represents the response from image generation.
type GenerateResponse struct {
	// Created is the timestamp when the response was created.
	Created time.Time `json:"created"`

	// Images contains the generated image(s).
	Images []Image `json:"images"`

	// Model is the model used for generation.
	Model string `json:"model"`

	// ProviderMetadata contains provider-specific metadata.
	ProviderMetadata map[string]any `json:"provider_metadata,omitempty"`
}

// Image represents a generated image.
type Image struct {
	// URL is the URL to the generated image (if ResponseFormat is FormatURL).
	URL string `json:"url,omitempty"`

	// B64JSON is the base64-encoded image data (if ResponseFormat is FormatB64JSON).
	B64JSON string `json:"b64_json,omitempty"`

	// RevisedPrompt is the prompt after any modifications by the model.
	// DALL-E 3 may rewrite prompts for safety or clarity.
	RevisedPrompt string `json:"revised_prompt,omitempty"`

	// ContentType is the MIME type of the image (e.g., "image/png").
	ContentType string `json:"content_type,omitempty"`

	// Width is the width of the image in pixels.
	Width int `json:"width,omitempty"`

	// Height is the height of the image in pixels.
	Height int `json:"height,omitempty"`

	// Seed is the seed used for this specific image (if available).
	Seed *int64 `json:"seed,omitempty"`
}

// EditRequest represents a request to edit an existing image.
type EditRequest struct {
	// Model is the model to use for editing.
	Model string `json:"model"`

	// Image is the original image to edit (base64 or URL).
	Image string `json:"image"`

	// Mask is an optional mask indicating areas to edit (base64 or URL).
	// Transparent areas indicate where the image should be edited.
	Mask string `json:"mask,omitempty"`

	// Prompt describes the desired edit.
	Prompt string `json:"prompt"`

	// N is the number of edited images to generate.
	N int `json:"n,omitempty"`

	// Size is the dimensions of the output image.
	Size ImageSize `json:"size,omitempty"`

	// ResponseFormat specifies how the image is returned.
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`

	// User is an optional unique identifier for the end-user.
	User string `json:"user,omitempty"`
}

// EditResponse represents the response from image editing.
type EditResponse struct {
	// Created is the timestamp when the response was created.
	Created time.Time `json:"created"`

	// Images contains the edited image(s).
	Images []Image `json:"images"`

	// ProviderMetadata contains provider-specific metadata.
	ProviderMetadata map[string]any `json:"provider_metadata,omitempty"`
}

// VariationsRequest represents a request to create variations of an image.
type VariationsRequest struct {
	// Model is the model to use for variations.
	Model string `json:"model"`

	// Image is the source image (base64 or URL).
	Image string `json:"image"`

	// N is the number of variations to generate.
	N int `json:"n,omitempty"`

	// Size is the dimensions of the output images.
	Size ImageSize `json:"size,omitempty"`

	// ResponseFormat specifies how images are returned.
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`

	// User is an optional unique identifier for the end-user.
	User string `json:"user,omitempty"`
}

// VariationsResponse represents the response from variation generation.
type VariationsResponse struct {
	// Created is the timestamp when the response was created.
	Created time.Time `json:"created"`

	// Images contains the variation images.
	Images []Image `json:"images"`

	// ProviderMetadata contains provider-specific metadata.
	ProviderMetadata map[string]any `json:"provider_metadata,omitempty"`
}

// UpscaleRequest represents a request to upscale an image.
type UpscaleRequest struct {
	// Model is the model to use for upscaling.
	Model string `json:"model"`

	// Image is the source image to upscale (base64 or URL).
	Image string `json:"image"`

	// Scale is the upscaling factor (e.g., 2 for 2x, 4 for 4x).
	Scale int `json:"scale,omitempty"`

	// ResponseFormat specifies how the image is returned.
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`
}

// UpscaleResponse represents the response from image upscaling.
type UpscaleResponse struct {
	// Created is the timestamp when the response was created.
	Created time.Time `json:"created"`

	// Image is the upscaled image.
	Image Image `json:"image"`

	// ProviderMetadata contains provider-specific metadata.
	ProviderMetadata map[string]any `json:"provider_metadata,omitempty"`
}
