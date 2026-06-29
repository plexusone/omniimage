# Fal AI Provider

The Fal AI provider supports FLUX, Stable Diffusion, and image upscaling models.

## Setup

```go
import (
    "github.com/plexusone/omniimage"
    "github.com/plexusone/omniimage/provider"
)

client, err := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {
            Provider: omniimage.ProviderNameFal,
            APIKey:   "...", // Or use FAL_KEY env var
        },
    },
})
```

## Models

### Generation Models

| Model | Constant | Quality | Speed | Best For |
|-------|----------|---------|-------|----------|
| FLUX Pro | `ModelFluxPro` | Highest | Slow | Production images |
| FLUX Dev | `ModelFluxDev` | High | Medium | Development/testing |
| FLUX Schnell | `ModelFluxSchnell` | Good | Fast | Rapid iteration |
| SDXL | `ModelSDXL` | High | Medium | Stable Diffusion style |

### Upscaling Models

| Model | Constant | Description |
|-------|----------|-------------|
| Clarity Upscaler | `ModelClarityUpscaler` | General upscaling |
| Creative Upscaler | `ModelCreativeUpscale` | Enhanced details |

## Image Generation

### Basic Generation

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxPro,
    Prompt: "A photorealistic portrait of a woman",
    Size:   provider.Size1024x1024,
})
```

### Multiple Images

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxDev,
    Prompt: "Abstract geometric art",
    Size:   provider.Size1024x1024,
    N:      4,
})

for _, img := range resp.Images {
    fmt.Println(img.URL)
}
```

### Negative Prompts

Specify what to avoid in the generated image:

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:          omniimage.ModelFluxPro,
    Prompt:         "A professional headshot",
    NegativePrompt: "blurry, distorted, low quality, cartoon",
    Size:           provider.Size1024x1024,
})
```

### Reproducible Generation

Use seeds for consistent results:

```go
seed := int64(12345)
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxPro,
    Prompt: "A mountain landscape",
    Seed:   &seed,
})

// Same seed + same prompt = same image
fmt.Printf("Used seed: %d\n", *resp.Images[0].Seed)
```

### Fine-Tuned Control

Control inference steps and guidance:

```go
steps := 30
guidance := 7.5

resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:         omniimage.ModelFluxPro,
    Prompt:        "A detailed fantasy scene",
    Steps:         &steps,         // More steps = more detail
    GuidanceScale: &guidance,      // Higher = closer to prompt
})
```

## Image Upscaling

Increase image resolution:

```go
resp, err := client.Upscale(ctx, &provider.UpscaleRequest{
    Model: omniimage.ModelClarityUpscaler,
    Image: "https://example.com/image.jpg",
    Scale: 2, // 2x upscale
})

fmt.Printf("Upscaled: %dx%d\n", resp.Image.Width, resp.Image.Height)
```

### Creative Upscaling

Add detail while upscaling:

```go
resp, err := client.Upscale(ctx, &provider.UpscaleRequest{
    Model: omniimage.ModelCreativeUpscale,
    Image: "https://example.com/low-res.jpg",
    Scale: 4,
})
```

## Provider-Specific Options

Pass additional Fal-specific parameters:

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxPro,
    Prompt: "A cyberpunk cityscape",
    ProviderOptions: map[string]any{
        "output_format":    "png",
        "safety_tolerance": 2,
        "enhance_prompt":   true,
    },
})
```

## Response Metadata

Fal returns additional metadata:

```go
resp, err := client.Generate(ctx, req)

// Access timing information
if timing, ok := resp.ProviderMetadata["inference_time"].(float64); ok {
    fmt.Printf("Generation took %.2f seconds\n", timing)
}

// Access the seed used
if seed, ok := resp.ProviderMetadata["seed"].(int64); ok {
    fmt.Printf("Seed: %d\n", seed)
}
```

## Supported Sizes

Fal supports custom dimensions. Standard sizes:

| Size | Constant | Aspect Ratio |
|------|----------|--------------|
| 512x512 | `Size512x512` | 1:1 |
| 768x1024 | `Size768x1024` | 3:4 |
| 1024x768 | `Size1024x768` | 4:3 |
| 1024x1024 | `Size1024x1024` | 1:1 |

Custom sizes via ProviderOptions:

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxPro,
    Prompt: "A banner image",
    ProviderOptions: map[string]any{
        "image_size": map[string]int{
            "width":  1920,
            "height": 480,
        },
    },
})
```

## Model Selection Guide

| Use Case | Recommended Model |
|----------|-------------------|
| Production images | `ModelFluxPro` |
| Prototyping | `ModelFluxSchnell` |
| Fine-tuning workflow | `ModelFluxDev` |
| SD-style images | `ModelSDXL` |
| Upscaling photos | `ModelClarityUpscaler` |
| Upscaling art | `ModelCreativeUpscale` |

## Error Handling

```go
resp, err := client.Generate(ctx, req)
if err != nil {
    var apiErr *fal.APIError
    if errors.As(err, &apiErr) {
        log.Printf("Fal API error: %s", apiErr.Message)
    }
    return
}
```

## Best Practices

1. **Use negative prompts** - Significantly improves quality
2. **Save seeds** - Store seeds for reproducible results
3. **Start with Schnell** - Iterate fast, then switch to Pro
4. **Tune guidance scale** - 3-7 for creative, 7-15 for precise
5. **Adjust steps wisely** - More isn't always better (20-30 typical)
