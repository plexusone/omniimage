# Providers API

The provider interface defines how image generation backends are implemented.

## Provider Interface

All providers must implement the base `Provider` interface:

```go
type Provider interface {
    // Generate creates images from a text prompt.
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

    // Close releases any resources held by the provider.
    Close() error

    // Name returns the provider name.
    Name() string
}
```

## Optional Interfaces

Providers can implement additional interfaces for extended functionality:

### EditProvider

```go
type EditProvider interface {
    Provider
    Edit(ctx context.Context, req *EditRequest) (*EditResponse, error)
}
```

Implemented by: OpenAI (DALL-E 2)

### VariationsProvider

```go
type VariationsProvider interface {
    Provider
    Variations(ctx context.Context, req *VariationsRequest) (*VariationsResponse, error)
}
```

Implemented by: OpenAI (DALL-E 2)

### UpscaleProvider

```go
type UpscaleProvider interface {
    Provider
    Upscale(ctx context.Context, req *UpscaleRequest) (*UpscaleResponse, error)
}
```

Implemented by: Fal AI (Clarity Upscaler)

## Built-in Providers

### OpenAI Provider

```go
import "github.com/plexusone/omniimage/providers/openai"

p, err := openai.New(openai.Config{
    APIKey:  "sk-...",
    BaseURL: "", // Optional, defaults to api.openai.com
})
```

**Supported Operations:**

| Operation | Supported | Models |
|-----------|-----------|--------|
| Generate | Yes | All |
| Edit | Yes | DALL-E 2 |
| Variations | Yes | DALL-E 2 |
| Upscale | No | - |

### Fal AI Provider

```go
import "github.com/plexusone/omniimage/providers/fal"

p, err := fal.New(fal.Config{
    APIKey: "your-fal-api-key",
})
```

**Supported Operations:**

| Operation | Supported | Models |
|-----------|-----------|--------|
| Generate | Yes | All |
| Edit | No | - |
| Variations | No | - |
| Upscale | Yes | Clarity Upscaler |

## Custom Providers

You can implement custom providers by implementing the `Provider` interface:

```go
type MyProvider struct {
    // ... your fields
}

func (p *MyProvider) Generate(ctx context.Context, req *provider.GenerateRequest) (*provider.GenerateResponse, error) {
    // Your implementation
}

func (p *MyProvider) Close() error {
    // Cleanup resources
    return nil
}

func (p *MyProvider) Name() string {
    return "my-provider"
}
```

### Registering Custom Providers

Use `CustomProvider` in the config to inject your implementation:

```go
client, err := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {
            Provider:       "custom",
            CustomProvider: &MyProvider{},
        },
    },
})
```

## Provider Factory

The `NewProvider` function creates providers from configuration:

```go
func NewProvider(config ProviderConfig) (provider.Provider, error)
```

**Parameters:**

| Field | Type | Description |
|-------|------|-------------|
| `Provider` | `ProviderName` | Provider name (`openai`, `fal`) |
| `APIKey` | `string` | API key (or use env var) |
| `BaseURL` | `string` | Custom base URL (optional) |
| `CustomProvider` | `provider.Provider` | Custom implementation |

**Environment Variables:**

If `APIKey` is not set, providers check these environment variables:

| Provider | Environment Variable |
|----------|---------------------|
| OpenAI | `OPENAI_API_KEY` |
| Fal AI | `FAL_KEY` |

## Capability Detection

Check provider capabilities before calling operations:

```go
// Check at runtime via client
if client.SupportsEdit() {
    resp, err := client.Edit(ctx, editReq)
}

// Check statically via capabilities
caps := omniimage.GetCapabilities(omniimage.ProviderNameOpenAI)
fmt.Printf("Supports Edit: %v\n", caps.Edit)
fmt.Printf("Supports Upscale: %v\n", caps.Upscale)
fmt.Printf("Max Images: %d\n", caps.MaxImages)
```

### Capabilities Struct

```go
type Capabilities struct {
    Generate       bool
    Edit           bool
    Variations     bool
    Upscale        bool
    Models         []string
    MaxImages      int
    SupportedSizes []provider.ImageSize
}
```

## Type Assertions

For direct access to optional interfaces:

```go
p := client.Provider()

// Check for edit support
if editor, ok := p.(provider.EditProvider); ok {
    resp, err := editor.Edit(ctx, editReq)
}

// Check for upscale support
if upscaler, ok := p.(provider.UpscaleProvider); ok {
    resp, err := upscaler.Upscale(ctx, upscaleReq)
}
```
