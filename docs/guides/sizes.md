# Image Sizes

omniimage provides standard size constants for common image dimensions.

## Size Constants

| Constant | Dimensions | Aspect Ratio | Use Case |
|----------|------------|--------------|----------|
| `Size256x256` | 256x256 | 1:1 | Thumbnails, icons |
| `Size512x512` | 512x512 | 1:1 | Small images, previews |
| `Size768x1024` | 768x1024 | 3:4 | Portrait photos |
| `Size1024x768` | 1024x768 | 4:3 | Landscape photos |
| `Size1024x1024` | 1024x1024 | 1:1 | Standard, social media |
| `Size1024x1792` | 1024x1792 | 9:16 | Tall portrait, mobile |
| `Size1792x1024` | 1792x1024 | 16:9 | Wide landscape, banners |

## Provider Support

| Size | OpenAI | Fal AI |
|------|--------|--------|
| 256x256 | DALL-E 2 | No |
| 512x512 | DALL-E 2 | Yes |
| 768x1024 | No | Yes |
| 1024x768 | No | Yes |
| 1024x1024 | All | Yes |
| 1024x1792 | GPT Image, DALL-E 3 | No |
| 1792x1024 | GPT Image, DALL-E 3 | No |

## Usage

```go
import "github.com/plexusone/omniimage/provider"

// Square image
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A profile picture",
    Size:   provider.Size1024x1024,
})

// Portrait
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A full-body portrait",
    Size:   provider.Size1024x1792,
})

// Landscape
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A panoramic view",
    Size:   provider.Size1792x1024,
})
```

## Custom Sizes (Fal AI)

Fal AI supports arbitrary dimensions via ProviderOptions:

```go
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxPro,
    Prompt: "A website banner",
    ProviderOptions: map[string]any{
        "image_size": map[string]int{
            "width":  1920,
            "height": 480,
        },
    },
})
```

## Choosing the Right Size

| Content Type | Recommended Size |
|--------------|------------------|
| Profile pictures | 1024x1024 |
| Social media posts | 1024x1024 |
| Mobile wallpapers | 1024x1792 |
| Desktop wallpapers | 1792x1024 |
| Thumbnails | 256x256 or 512x512 |
| Product images | 1024x1024 |
| Banners | Custom via Fal AI |

## Size and Quality

Larger sizes generally provide:

- More detail
- Better for printing
- Larger file sizes
- Longer generation time

!!! tip "Performance"
    For rapid prototyping, use smaller sizes (512x512). Scale up to
    production sizes once you've refined your prompt.
