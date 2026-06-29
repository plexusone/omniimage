# omniimage

Unified Go library for image generation across multiple providers.

## Features

- **Multiple Providers** - OpenAI (GPT Image, DALL-E) and Fal AI (FLUX, Stable Diffusion)
- **Unified Interface** - Single API for all providers
- **Type Safety** - Strongly typed requests and responses
- **Extensible** - Easy to add custom providers

## Provider Comparison

| Provider | Generate | Edit | Variations | Upscale |
|----------|:--------:|:----:|:----------:|:-------:|
| OpenAI   | Yes      | Yes  | Yes        | No      |
| Fal AI   | Yes      | No   | No         | Yes     |

## Supported Models

### OpenAI

| Model | Description | Max Images |
|-------|-------------|------------|
| `gpt-image-2` | Latest GPT Image model | 10 |
| `gpt-image-1` | First generation GPT Image | 10 |
| `dall-e-3` | DALL-E 3 (legacy) | 1 |
| `dall-e-2` | DALL-E 2 (legacy) | 10 |

### Fal AI

| Model | Description | Max Images |
|-------|-------------|------------|
| `fal-ai/flux-pro` | FLUX Pro - highest quality | 4 |
| `fal-ai/flux/dev` | FLUX Dev - balanced | 4 |
| `fal-ai/flux/schnell` | FLUX Schnell - fastest | 4 |
| `fal-ai/fast-sdxl` | Stable Diffusion XL | 4 |

## Quick Example

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
    client, err := omniimage.NewClient(omniimage.ClientConfig{
        Providers: []omniimage.ProviderConfig{
            {Provider: omniimage.ProviderNameOpenAI},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    resp, err := client.Generate(context.Background(), &provider.GenerateRequest{
        Model:  omniimage.ModelGPTImage2,
        Prompt: "A serene mountain landscape at sunset",
        Size:   provider.Size1024x1024,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated: %s\n", resp.Images[0].URL)
}
```

## Installation

```bash
go get github.com/plexusone/omniimage
```

## Next Steps

- [Getting Started](getting-started.md) - Step-by-step setup guide
- [Providers](providers/index.md) - Provider-specific documentation
- [API Reference](api/client.md) - Complete API documentation
