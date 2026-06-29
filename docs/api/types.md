# Types Reference

Request and response types for image generation.

## Request Types

### GenerateRequest

Request for image generation.

```go
type GenerateRequest struct {
    // Model is the model to use for generation (required).
    Model string

    // Prompt is the text description (required).
    Prompt string

    // NegativePrompt specifies what to avoid (Fal AI only).
    NegativePrompt string

    // N is the number of images to generate (default: 1).
    N int

    // Size is the image dimensions.
    Size ImageSize

    // Quality is the quality level (standard or hd).
    Quality ImageQuality

    // Style affects the visual style (vivid or natural).
    Style ImageStyle

    // ResponseFormat specifies URL or base64 response.
    ResponseFormat ResponseFormat

    // User is an optional end-user identifier.
    User string

    // Seed for reproducible generation (Fal AI).
    Seed *int64

    // Steps is the number of inference steps (Fal AI).
    Steps *int

    // GuidanceScale controls prompt adherence (Fal AI).
    GuidanceScale *float64

    // ProviderOptions for provider-specific parameters.
    ProviderOptions map[string]any
}
```

### EditRequest

Request for image editing.

```go
type EditRequest struct {
    // Model is the model to use (required).
    Model string

    // Image is the source image (base64 or URL, required).
    Image string

    // Mask indicates edit areas (transparent = edit).
    Mask string

    // Prompt describes the desired edit (required).
    Prompt string

    // N is the number of edited images.
    N int

    // Size is the output dimensions.
    Size ImageSize

    // ResponseFormat specifies URL or base64 response.
    ResponseFormat ResponseFormat

    // User is an optional end-user identifier.
    User string
}
```

### VariationsRequest

Request for image variations.

```go
type VariationsRequest struct {
    // Model is the model to use (required).
    Model string

    // Image is the source image (base64 or URL, required).
    Image string

    // N is the number of variations.
    N int

    // Size is the output dimensions.
    Size ImageSize

    // ResponseFormat specifies URL or base64 response.
    ResponseFormat ResponseFormat

    // User is an optional end-user identifier.
    User string
}
```

### UpscaleRequest

Request for image upscaling.

```go
type UpscaleRequest struct {
    // Model is the upscaling model (required).
    Model string

    // Image is the source image (base64 or URL, required).
    Image string

    // Scale is the upscaling factor (e.g., 2 for 2x).
    Scale int

    // ResponseFormat specifies URL or base64 response.
    ResponseFormat ResponseFormat
}
```

## Response Types

### GenerateResponse

```go
type GenerateResponse struct {
    // Created is the generation timestamp.
    Created time.Time

    // Images contains the generated images.
    Images []Image

    // Model is the model used.
    Model string

    // ProviderMetadata contains provider-specific data.
    ProviderMetadata map[string]any
}
```

### EditResponse

```go
type EditResponse struct {
    Created          time.Time
    Images           []Image
    ProviderMetadata map[string]any
}
```

### VariationsResponse

```go
type VariationsResponse struct {
    Created          time.Time
    Images           []Image
    ProviderMetadata map[string]any
}
```

### UpscaleResponse

```go
type UpscaleResponse struct {
    Created          time.Time
    Image            Image
    ProviderMetadata map[string]any
}
```

### Image

```go
type Image struct {
    // URL is the image URL (if FormatURL).
    URL string

    // B64JSON is base64 data (if FormatB64JSON).
    B64JSON string

    // RevisedPrompt is the modified prompt (if any).
    RevisedPrompt string

    // ContentType is the MIME type.
    ContentType string

    // Width in pixels.
    Width int

    // Height in pixels.
    Height int

    // Seed used for this image (Fal AI).
    Seed *int64
}
```

## Enums

### ImageSize

```go
type ImageSize string

const (
    Size256x256   ImageSize = "256x256"
    Size512x512   ImageSize = "512x512"
    Size768x1024  ImageSize = "768x1024"
    Size1024x768  ImageSize = "1024x768"
    Size1024x1024 ImageSize = "1024x1024"
    Size1024x1792 ImageSize = "1024x1792"
    Size1792x1024 ImageSize = "1792x1024"
)
```

### ImageQuality

```go
type ImageQuality string

const (
    QualityStandard ImageQuality = "standard"
    QualityHD       ImageQuality = "hd"
)
```

### ImageStyle

```go
type ImageStyle string

const (
    StyleVivid   ImageStyle = "vivid"
    StyleNatural ImageStyle = "natural"
)
```

### ResponseFormat

```go
type ResponseFormat string

const (
    FormatURL     ResponseFormat = "url"
    FormatB64JSON ResponseFormat = "b64_json"
)
```

## Constants

### Provider Names

```go
const (
    ProviderNameOpenAI ProviderName = "openai"
    ProviderNameFal    ProviderName = "fal"
)
```

### Model Constants

```go
// OpenAI
const (
    ModelGPTImage2 = "gpt-image-2"
    ModelGPTImage1 = "gpt-image-1"
    ModelDALLE3    = "dall-e-3"
    ModelDALLE2    = "dall-e-2"
)

// Fal AI
const (
    ModelFluxPro         = "fal-ai/flux-pro"
    ModelFluxDev         = "fal-ai/flux/dev"
    ModelFluxSchnell     = "fal-ai/flux/schnell"
    ModelSDXL            = "fal-ai/fast-sdxl"
    ModelClarityUpscaler = "fal-ai/clarity-upscaler"
)
```

## Errors

```go
var (
    ErrNoProviders     = errors.New("no providers configured")
    ErrUnknownProvider = errors.New("unknown provider")
    ErrNotSupported    = errors.New("operation not supported")
    ErrInvalidRequest  = errors.New("invalid request")
    ErrRateLimited     = errors.New("rate limited")
    ErrContentPolicy   = errors.New("content policy violation")
    ErrModelNotFound   = errors.New("model not found")
)
```

### APIError

```go
type APIError struct {
    StatusCode int
    Code       string
    Message    string
    Provider   string
}

func (e *APIError) Error() string
```
