# OmniImage

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/omniimage/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/omniimage/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/omniimage/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/omniimage/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/omniimage/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/omniimage/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/omniimage
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/omniimage
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/omniimage
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/omniimage
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/omniimage
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fomniimage
 [loc-svg]: https://tokei.rs/b1/github/plexusone/omniimage
 [repo-url]: https://github.com/plexusone/omniimage
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omniimage/blob/main/LICENSE

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
