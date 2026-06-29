# Best Practices

Guidelines for getting the best results with omniimage.

## Prompt Engineering

### Be Specific

```go
// Too vague
"A cat"

// Better
"A fluffy orange tabby cat sitting on a windowsill, soft natural lighting"

// Best
"A fluffy orange tabby cat with green eyes sitting on a wooden windowsill,
morning sunlight streaming through the window, cozy home interior,
photorealistic, shallow depth of field"
```

### Include Style Keywords

```go
// Specify artistic style
"oil painting of a mountain landscape"
"digital illustration of a robot"
"photorealistic portrait of a woman"
"watercolor sketch of a city street"
"3D render of a futuristic car"
```

### Describe Composition

```go
// Include composition details
"close-up portrait, centered framing"
"wide-angle landscape, rule of thirds"
"bird's eye view of a city"
"low angle shot of a skyscraper"
```

## Provider Selection

### Choose the Right Provider

| Scenario | Provider | Model |
|----------|----------|-------|
| Production images | OpenAI | `gpt-image-2` |
| Fast prototyping | Fal AI | `flux-schnell` |
| Reproducible results | Fal AI | Any (with seed) |
| Image editing | OpenAI | `dall-e-2` |
| Upscaling | Fal AI | `clarity-upscaler` |

### Use Seeds for Consistency

```go
// Fal AI: same seed = same output
seed := int64(42)
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelFluxPro,
    Prompt: "A logo design",
    Seed:   &seed,
})

// Store the seed for reproducibility
fmt.Printf("Seed used: %d\n", *resp.Images[0].Seed)
```

## Performance

### Batch When Possible

```go
// Instead of multiple single requests
for i := 0; i < 4; i++ {
    client.Generate(ctx, &provider.GenerateRequest{N: 1, ...})
}

// Use N parameter for batch generation
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "Product variations",
    N:      4,
})
```

### Use Appropriate Sizes

```go
// For prototyping - smaller, faster
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Size: provider.Size512x512,
    ...
})

// For production - full quality
resp, _ := client.Generate(ctx, &provider.GenerateRequest{
    Size:    provider.Size1024x1024,
    Quality: provider.QualityHD,
    ...
})
```

### Handle Timeouts

```go
// Set appropriate context timeout
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

resp, err := client.Generate(ctx, req)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Handle timeout
    }
}
```

## Error Handling

### Implement Retries

```go
func generateWithRetry(ctx context.Context, req *provider.GenerateRequest) (*provider.GenerateResponse, error) {
    var lastErr error

    for attempt := 0; attempt < 3; attempt++ {
        if attempt > 0 {
            time.Sleep(time.Duration(attempt) * time.Second)
        }

        resp, err := client.Generate(ctx, req)
        if err == nil {
            return resp, nil
        }

        lastErr = err

        // Don't retry non-retryable errors
        var apiErr *omniimage.APIError
        if errors.As(err, &apiErr) {
            if apiErr.StatusCode == 400 || apiErr.StatusCode == 403 {
                return nil, err // Bad request or policy violation
            }
        }
    }

    return nil, fmt.Errorf("after 3 attempts: %w", lastErr)
}
```

### Validate Inputs

```go
func validateRequest(req *provider.GenerateRequest) error {
    if req.Prompt == "" {
        return fmt.Errorf("prompt is required")
    }

    if len(req.Prompt) > 4000 {
        return fmt.Errorf("prompt too long (max 4000 chars)")
    }

    if req.N > 10 {
        return fmt.Errorf("max 10 images per request")
    }

    return nil
}
```

## Resource Management

### Close Clients

```go
client, err := omniimage.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close() // Always close

// Use client...
```

### Reuse Clients

```go
// Good: create once, reuse
var imageClient *omniimage.Client

func init() {
    var err error
    imageClient, err = omniimage.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    resp, err := imageClient.Generate(r.Context(), req)
    // ...
}

// Bad: create per request
func handleRequest(w http.ResponseWriter, r *http.Request) {
    client, _ := omniimage.NewClient(config) // Wasteful
    defer client.Close()
    resp, _ := client.Generate(r.Context(), req)
}
```

## Security

### Protect API Keys

```go
// Good: use environment variables
client, _ := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {Provider: omniimage.ProviderNameOpenAI}, // Uses OPENAI_API_KEY
    },
})

// Bad: hardcoded keys
{
    Provider: omniimage.ProviderNameOpenAI,
    APIKey:   "sk-abc123...", // Don't do this!
}
```

### Validate User Input

```go
func generateFromUserInput(userPrompt string) error {
    // Sanitize/validate user input
    prompt := sanitizePrompt(userPrompt)

    // Consider content moderation
    if containsProhibitedContent(prompt) {
        return fmt.Errorf("prohibited content")
    }

    resp, err := client.Generate(ctx, &provider.GenerateRequest{
        Model:  omniimage.ModelGPTImage2,
        Prompt: prompt,
    })
    // ...
}
```

## Cost Optimization

### Use Appropriate Quality

| Quality | Cost | Use Case |
|---------|------|----------|
| Standard | Lower | Previews, prototypes |
| HD | Higher | Final assets |

### Monitor Usage

```go
// Log generation requests for cost tracking
func generateImage(ctx context.Context, req *provider.GenerateRequest) (*provider.GenerateResponse, error) {
    start := time.Now()

    resp, err := client.Generate(ctx, req)

    slog.Info("image generated",
        "model", req.Model,
        "size", req.Size,
        "n", req.N,
        "duration", time.Since(start),
    )

    return resp, err
}
```
