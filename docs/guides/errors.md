# Error Handling

omniimage provides structured error types for handling API errors gracefully.

## Error Types

### Sentinel Errors

```go
import "github.com/plexusone/omniimage"

// Check for specific error types
if errors.Is(err, omniimage.ErrNoProviders) {
    // No providers configured
}

if errors.Is(err, omniimage.ErrUnknownProvider) {
    // Invalid provider name
}

if errors.Is(err, omniimage.ErrNotSupported) {
    // Operation not supported by provider
}

if errors.Is(err, omniimage.ErrRateLimited) {
    // API rate limit exceeded
}

if errors.Is(err, omniimage.ErrContentPolicy) {
    // Content violates provider policy
}
```

### API Errors

```go
var apiErr *omniimage.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("Provider: %s\n", apiErr.Provider)
    fmt.Printf("Status: %d\n", apiErr.StatusCode)
    fmt.Printf("Code: %s\n", apiErr.Code)
    fmt.Printf("Message: %s\n", apiErr.Message)
}
```

## Common Error Scenarios

### Invalid API Key

```go
resp, err := client.Generate(ctx, req)
if err != nil {
    var apiErr *omniimage.APIError
    if errors.As(err, &apiErr) && apiErr.StatusCode == 401 {
        log.Fatal("Invalid API key")
    }
}
```

### Rate Limiting

```go
func generateWithRetry(ctx context.Context, client *omniimage.Client, req *provider.GenerateRequest) (*provider.GenerateResponse, error) {
    maxRetries := 3

    for i := 0; i < maxRetries; i++ {
        resp, err := client.Generate(ctx, req)
        if err == nil {
            return resp, nil
        }

        var apiErr *omniimage.APIError
        if errors.As(err, &apiErr) && apiErr.StatusCode == 429 {
            // Exponential backoff
            time.Sleep(time.Duration(1<<i) * time.Second)
            continue
        }

        return nil, err
    }

    return nil, fmt.Errorf("max retries exceeded")
}
```

### Content Policy Violations

```go
resp, err := client.Generate(ctx, req)
if err != nil {
    var apiErr *omniimage.APIError
    if errors.As(err, &apiErr) {
        if apiErr.Code == "content_policy_violation" {
            log.Printf("Prompt violates content policy: %s", apiErr.Message)
            // Modify prompt and retry
            return
        }
    }
}
```

### Unsupported Operations

```go
// Check before calling
if !client.SupportsEdit() {
    return fmt.Errorf("provider does not support editing")
}

resp, err := client.Edit(ctx, editReq)
if err != nil {
    if errors.Is(err, omniimage.ErrNotSupported) {
        // Fallback to different approach
    }
}
```

## Provider-Specific Errors

### OpenAI

| Status | Code | Description |
|--------|------|-------------|
| 400 | `invalid_request_error` | Bad request parameters |
| 401 | `invalid_api_key` | Invalid API key |
| 403 | `content_policy_violation` | Prompt violates policy |
| 429 | `rate_limit_exceeded` | Too many requests |
| 500 | `server_error` | OpenAI server error |

### Fal AI

| Status | Description |
|--------|-------------|
| 400 | Invalid request |
| 401 | Invalid API key |
| 429 | Rate limited |
| 500 | Server error |

## Best Practices

### Wrap Errors with Context

```go
resp, err := client.Generate(ctx, req)
if err != nil {
    return fmt.Errorf("generate image for %s: %w", userID, err)
}
```

### Log Errors Appropriately

```go
resp, err := client.Generate(ctx, req)
if err != nil {
    var apiErr *omniimage.APIError
    if errors.As(err, &apiErr) {
        slog.Error("image generation failed",
            "provider", apiErr.Provider,
            "status", apiErr.StatusCode,
            "code", apiErr.Code,
            "message", apiErr.Message,
        )
    } else {
        slog.Error("image generation failed", "error", err)
    }
    return err
}
```

### Graceful Degradation

```go
func generateImage(ctx context.Context, prompt string) (string, error) {
    // Try primary provider
    resp, err := primaryClient.Generate(ctx, &provider.GenerateRequest{
        Model:  omniimage.ModelGPTImage2,
        Prompt: prompt,
    })
    if err == nil {
        return resp.Images[0].URL, nil
    }

    // Log and try fallback
    slog.Warn("primary provider failed, trying fallback", "error", err)

    resp, err = fallbackClient.Generate(ctx, &provider.GenerateRequest{
        Model:  omniimage.ModelFluxPro,
        Prompt: prompt,
    })
    if err != nil {
        return "", fmt.Errorf("all providers failed: %w", err)
    }

    return resp.Images[0].URL, nil
}
```
