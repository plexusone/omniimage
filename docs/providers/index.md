# Providers Overview

OmniImage supports multiple image generation providers through a unified interface.

## Available Providers

| Provider | Constant | Environment Variable |
|----------|----------|---------------------|
| OpenAI | `ProviderNameOpenAI` | `OPENAI_API_KEY` |
| Fal AI | `ProviderNameFal` | `FAL_KEY` |

## Capabilities Comparison

| Feature | OpenAI | Fal AI |
|---------|--------|--------|
| Generate | Yes | Yes |
| Edit | Yes | No |
| Variations | Yes | No |
| Upscale | No | Yes |
| Negative Prompts | No | Yes |
| Seed Control | No | Yes |
| Inference Steps | No | Yes |

## Choosing a Provider

### OpenAI

Best for:

- High-quality, general-purpose image generation
- Image editing with masks
- Creating variations of existing images
- Enterprise environments with existing OpenAI integration

### Fal AI

Best for:

- FLUX and Stable Diffusion models
- Fine-grained control (steps, guidance, seeds)
- Reproducible generation with seeds
- Image upscaling
- Cost-effective batch generation

## Provider Configuration

### Basic Setup

```go
client, err := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {
            Provider: omniimage.ProviderNameOpenAI,
            APIKey:   "sk-...", // Optional if env var set
        },
    },
})
```

### Custom Base URL

For proxies or self-hosted endpoints:

```go
{
    Provider: omniimage.ProviderNameOpenAI,
    APIKey:   "sk-...",
    BaseURL:  "https://my-proxy.example.com/v1",
}
```

### Custom Provider

Implement the `provider.Provider` interface:

```go
type MyProvider struct{}

func (p *MyProvider) Generate(ctx context.Context, req *provider.GenerateRequest) (*provider.GenerateResponse, error) {
    // Custom implementation
}

func (p *MyProvider) Close() error { return nil }
func (p *MyProvider) Name() string { return "my-provider" }

// Use custom provider
client, _ := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {CustomProvider: &MyProvider{}},
    },
})
```

## Checking Capabilities

```go
// Check provider capabilities at runtime
caps := omniimage.GetCapabilities(omniimage.ProviderNameOpenAI)
if caps.Edit {
    // Provider supports editing
}

// Check client capabilities
if client.SupportsEdit() {
    resp, _ := client.Edit(ctx, editReq)
}

if client.SupportsUpscale() {
    resp, _ := client.Upscale(ctx, upscaleReq)
}
```

## Model Selection

Each provider supports different models:

```go
// OpenAI models
omniimage.ModelGPTImage2    // Latest GPT Image
omniimage.ModelGPTImage1    // First gen GPT Image
omniimage.ModelDALLE3       // DALL-E 3 (legacy)
omniimage.ModelDALLE2       // DALL-E 2 (legacy)

// Fal AI models
omniimage.ModelFluxPro      // FLUX Pro
omniimage.ModelFluxDev      // FLUX Dev
omniimage.ModelFluxSchnell  // FLUX Schnell (fast)
omniimage.ModelSDXL         // Stable Diffusion XL
```

## Next Steps

- [OpenAI Provider](openai.md) - Detailed OpenAI documentation
- [Fal AI Provider](fal.md) - Detailed Fal AI documentation
