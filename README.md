# omniimage

Unified Go library for image generation across multiple providers.

## Supported Providers

| Provider | Models | Generate | Edit | Variations | Upscale |
|----------|--------|----------|------|------------|---------|
| **OpenAI** | DALL-E 2, DALL-E 3 | Yes | Yes | Yes | No |
| **Fal AI** | FLUX Pro/Dev/Schnell, SDXL | Yes | No | No | Yes |

## Installation

```bash
go get github.com/plexusone/omniimage
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/plexusone/omniimage"
    "github.com/plexusone/omniimage/provider"
)

func main() {
    // Create client with OpenAI provider
    client, err := omniimage.NewClient(omniimage.ClientConfig{
        Providers: []omniimage.ProviderConfig{
            {
                Provider: omniimage.ProviderNameOpenAI,
                APIKey:   "sk-...", // or use OPENAI_API_KEY env var
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Generate an image
    resp, err := client.Generate(context.Background(), &provider.GenerateRequest{
        Model:   omniimage.ModelDALLE3,
        Prompt:  "A serene mountain landscape at sunset",
        Size:    provider.Size1024x1024,
        Quality: provider.QualityHD,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated image: %s\n", resp.Images[0].URL)
}
```

## Providers

### OpenAI (DALL-E)

```go
client, _ := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {
            Provider: omniimage.ProviderNameOpenAI,
            APIKey:   "sk-...", // or OPENAI_API_KEY env var
        },
    },
})

// DALL-E 3 - high quality, single image
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:   omniimage.ModelDALLE3,
    Prompt:  "A futuristic cityscape",
    Size:    provider.Size1024x1024,
    Quality: provider.QualityHD,
    Style:   provider.StyleVivid,
})

// DALL-E 2 - multiple images, edit, variations
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelDALLE2,
    Prompt: "A cute robot",
    N:      4,
    Size:   provider.Size512x512,
})
```

### Fal AI (FLUX, Stable Diffusion)

```go
client, _ := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {
            Provider: omniimage.ProviderNameFal,
            APIKey:   "...", // or FAL_KEY env var
        },
    },
})

// FLUX Pro - high quality
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxPro,
    Prompt: "A photorealistic portrait",
    Size:   provider.Size1024x1024,
})

// FLUX Schnell - fast generation
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxSchnell,
    Prompt: "Abstract art",
    N:      4,
})

// Upscale an image
upscaled, _ := client.Upscale(ctx, &provider.UpscaleRequest{
    Model: omniimage.ModelClarityUpscaler,
    Image: "https://example.com/image.png",
    Scale: 2,
})
```

## API Reference

### GenerateRequest

| Field | Type | Description |
|-------|------|-------------|
| `Model` | string | Model ID (required) |
| `Prompt` | string | Text description (required) |
| `NegativePrompt` | string | What to avoid (Fal only) |
| `N` | int | Number of images |
| `Size` | ImageSize | Image dimensions |
| `Quality` | ImageQuality | "standard" or "hd" |
| `Style` | ImageStyle | "vivid" or "natural" |
| `Seed` | *int64 | For reproducibility |
| `Steps` | *int | Inference steps (Fal) |
| `GuidanceScale` | *float64 | Prompt adherence (Fal) |

### Image Sizes

- `Size256x256`, `Size512x512`, `Size1024x1024`
- `Size1024x1792` (portrait), `Size1792x1024` (landscape)
- `Size768x1024`, `Size1024x768`

### Models

**OpenAI:**
- `ModelDALLE3` - DALL-E 3 (highest quality, n=1 only)
- `ModelDALLE2` - DALL-E 2 (edit, variations support)

**Fal AI:**
- `ModelFluxPro` - FLUX Pro (high quality)
- `ModelFluxDev` - FLUX Dev (balanced)
- `ModelFluxSchnell` - FLUX Schnell (fast)
- `ModelSDXL` - Stable Diffusion XL
- `ModelClarityUpscaler` - Image upscaling

## Environment Variables

| Variable | Provider | Description |
|----------|----------|-------------|
| `OPENAI_API_KEY` | OpenAI | API key for DALL-E |
| `FAL_KEY` | Fal AI | API key for Fal |

## License

MIT
