# Getting Started

This guide walks you through setting up omniimage and generating your first image.

## Installation

```bash
go get github.com/plexusone/omniimage
```

## Configuration

### Environment Variables

Set your API keys as environment variables:

```bash
# OpenAI
export OPENAI_API_KEY="sk-..."

# Fal AI
export FAL_KEY="..."
```

### Programmatic Configuration

Alternatively, pass API keys directly:

```go
client, err := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {
            Provider: omniimage.ProviderNameOpenAI,
            APIKey:   "sk-...",
        },
    },
})
```

## Basic Usage

### Generate an Image

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
    // Create client
    client, err := omniimage.NewClient(omniimage.ClientConfig{
        Providers: []omniimage.ProviderConfig{
            {Provider: omniimage.ProviderNameOpenAI},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Generate image
    resp, err := client.Generate(context.Background(), &provider.GenerateRequest{
        Model:  omniimage.ModelGPTImage2,
        Prompt: "A cute robot learning to paint",
        Size:   provider.Size1024x1024,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Print result
    fmt.Printf("Image URL: %s\n", resp.Images[0].URL)
}
```

### Generate Multiple Images

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "Abstract geometric patterns",
    Size:   provider.Size1024x1024,
    N:      4, // Generate 4 images
})
if err != nil {
    log.Fatal(err)
}

for i, img := range resp.Images {
    fmt.Printf("Image %d: %s\n", i+1, img.URL)
}
```

### High Quality Images

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:   omniimage.ModelGPTImage2,
    Prompt:  "Professional product photography of a watch",
    Size:    provider.Size1024x1024,
    Quality: provider.QualityHD,
})
```

### Different Styles

```go
// Vivid style (dramatic, hyper-real)
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelDALLE3,
    Prompt: "A fantasy castle",
    Style:  provider.StyleVivid,
})

// Natural style (more realistic)
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelDALLE3,
    Prompt: "A forest path",
    Style:  provider.StyleNatural,
})
```

## Using Different Providers

### OpenAI

```go
client, _ := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {Provider: omniimage.ProviderNameOpenAI},
    },
})

resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A futuristic city",
})
```

### Fal AI

```go
client, _ := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {Provider: omniimage.ProviderNameFal},
    },
})

resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxPro,
    Prompt: "A photorealistic portrait",
})
```

## Getting Base64 Data

Instead of URLs, get images as base64-encoded data:

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:          omniimage.ModelGPTImage2,
    Prompt:         "A logo design",
    ResponseFormat: provider.FormatB64JSON,
})
if err != nil {
    log.Fatal(err)
}

// Decode base64 data
imageData, err := base64.StdEncoding.DecodeString(resp.Images[0].B64JSON)
if err != nil {
    log.Fatal(err)
}

// Save to file
err = os.WriteFile("image.png", imageData, 0644)
```

## Error Handling

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A test image",
})

if err != nil {
    // Check for specific error types
    var apiErr *omniimage.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: %s (code: %s)\n", apiErr.Message, apiErr.Code)
        return
    }

    // Check for common errors
    if errors.Is(err, omniimage.ErrRateLimited) {
        fmt.Println("Rate limited, try again later")
        return
    }

    log.Fatal(err)
}
```

## Next Steps

- [OpenAI Provider](providers/openai.md) - OpenAI-specific features
- [Fal AI Provider](providers/fal.md) - Fal AI-specific features
- [Image Sizes](guides/sizes.md) - Available image dimensions
- [API Reference](api/client.md) - Complete API documentation
