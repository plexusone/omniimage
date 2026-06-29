# OpenAI Provider

The OpenAI provider supports GPT Image models and legacy DALL-E models for image generation, editing, and variations.

## Setup

```go
import (
    "github.com/plexusone/omniimage"
    "github.com/plexusone/omniimage/provider"
)

client, err := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {
            Provider: omniimage.ProviderNameOpenAI,
            APIKey:   "sk-...", // Or use OPENAI_API_KEY env var
        },
    },
})
```

## Models

| Model | Constant | Quality | Max N | Edit | Variations |
|-------|----------|---------|-------|------|------------|
| GPT Image 2 | `ModelGPTImage2` | HD | 10 | No | No |
| GPT Image 1 | `ModelGPTImage1` | HD | 10 | No | No |
| DALL-E 3 | `ModelDALLE3` | HD | 1 | No | No |
| DALL-E 2 | `ModelDALLE2` | Standard | 10 | Yes | Yes |

!!! note "DALL-E 3 Retirement"
    DALL-E 3 was retired on March 4, 2026. Use GPT Image models for new projects.

## Image Generation

### Basic Generation

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A serene Japanese garden with cherry blossoms",
    Size:   provider.Size1024x1024,
})
```

### High Quality

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:   omniimage.ModelGPTImage2,
    Prompt:  "Professional product photography",
    Size:    provider.Size1024x1024,
    Quality: provider.QualityHD,
})
```

### Style Options

```go
// Vivid - dramatic, hyper-real images
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelDALLE3,
    Prompt: "A fantasy dragon",
    Style:  provider.StyleVivid,
})

// Natural - realistic, less dramatic
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelDALLE3,
    Prompt: "A forest landscape",
    Style:  provider.StyleNatural,
})
```

### Portrait and Landscape

```go
// Portrait (vertical)
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A full-body portrait",
    Size:   provider.Size1024x1792,
})

// Landscape (horizontal)
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A panoramic mountain view",
    Size:   provider.Size1792x1024,
})
```

## Image Editing

Edit existing images with masks (DALL-E 2 only).

```go
// Edit an image - replace masked area
resp, err := client.Edit(ctx, &provider.EditRequest{
    Model:  omniimage.ModelDALLE2,
    Image:  "data:image/png;base64,...", // Base64 encoded image
    Mask:   "data:image/png;base64,...", // Transparent areas = edit region
    Prompt: "A red sports car",
    Size:   provider.Size1024x1024,
})
```

!!! warning "Multipart Upload"
    The Edit API traditionally requires multipart form data. For file-based
    uploads, you may need to use the OpenAI SDK directly or encode images
    as data URLs.

## Image Variations

Create variations of existing images (DALL-E 2 only).

```go
resp, err := client.Variations(ctx, &provider.VariationsRequest{
    Model: omniimage.ModelDALLE2,
    Image: "data:image/png;base64,...",
    N:     4, // Generate 4 variations
    Size:  provider.Size1024x1024,
})
```

## Response Format

### URL Response (Default)

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:          omniimage.ModelGPTImage2,
    Prompt:         "A test image",
    ResponseFormat: provider.FormatURL, // Default
})

imageURL := resp.Images[0].URL
```

### Base64 Response

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:          omniimage.ModelGPTImage2,
    Prompt:         "A test image",
    ResponseFormat: provider.FormatB64JSON,
})

b64Data := resp.Images[0].B64JSON
```

## Revised Prompts

GPT Image and DALL-E 3 may revise your prompt for safety or clarity:

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A person",
})

fmt.Println("Original:", "A person")
fmt.Println("Revised:", resp.Images[0].RevisedPrompt)
// Output might be: "A person standing in a park on a sunny day..."
```

## Supported Sizes

| Size | Constant | Models |
|------|----------|--------|
| 256x256 | `Size256x256` | DALL-E 2 |
| 512x512 | `Size512x512` | DALL-E 2 |
| 1024x1024 | `Size1024x1024` | All |
| 1024x1792 | `Size1024x1792` | GPT Image, DALL-E 3 |
| 1792x1024 | `Size1792x1024` | GPT Image, DALL-E 3 |

## Error Handling

```go
resp, err := client.Generate(ctx, req)
if err != nil {
    var apiErr *openai.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.Code {
        case "content_policy_violation":
            log.Println("Content violates policy")
        case "rate_limit_exceeded":
            log.Println("Rate limited, retry later")
        default:
            log.Printf("API error: %s", apiErr.Message)
        }
    }
    return
}
```

## Best Practices

1. **Use GPT Image models** - DALL-E 3 is deprecated
2. **Be specific** - Detailed prompts yield better results
3. **Handle revised prompts** - Check if your prompt was modified
4. **Use appropriate sizes** - Match size to use case
5. **Implement retries** - Handle rate limits gracefully
