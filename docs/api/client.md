# Client API

The `Client` struct provides the main interface for image generation.

## Creating a Client

```go
func NewClient(config ClientConfig) (*Client, error)
```

### ClientConfig

```go
type ClientConfig struct {
    // Providers is a list of provider configurations.
    // Index 0 is the primary provider.
    Providers []ProviderConfig

    // Logger for internal logging (optional).
    Logger *slog.Logger
}
```

### ProviderConfig

```go
type ProviderConfig struct {
    // Provider is the provider name.
    Provider ProviderName

    // APIKey is the API key for authentication.
    APIKey string

    // BaseURL is an optional custom base URL for the API.
    BaseURL string

    // CustomProvider allows injecting a custom provider implementation.
    CustomProvider provider.Provider
}
```

### Example

```go
client, err := omniimage.NewClient(omniimage.ClientConfig{
    Providers: []omniimage.ProviderConfig{
        {
            Provider: omniimage.ProviderNameOpenAI,
            APIKey:   "sk-...",
        },
    },
})
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## Methods

### Generate

Generate images from a text prompt.

```go
func (c *Client) Generate(ctx context.Context, req *provider.GenerateRequest) (*provider.GenerateResponse, error)
```

**Example:**

```go
resp, err := client.Generate(ctx, &provider.GenerateRequest{
    Model:  omniimage.ModelGPTImage2,
    Prompt: "A sunset over the ocean",
    Size:   provider.Size1024x1024,
    N:      1,
})
```

### Edit

Edit an existing image with a mask.

```go
func (c *Client) Edit(ctx context.Context, req *provider.EditRequest) (*provider.EditResponse, error)
```

Returns `ErrNotSupported` if the provider doesn't support editing.

**Example:**

```go
resp, err := client.Edit(ctx, &provider.EditRequest{
    Model:  omniimage.ModelDALLE2,
    Image:  "data:image/png;base64,...",
    Mask:   "data:image/png;base64,...",
    Prompt: "Add a red car",
})
```

### Variations

Create variations of an existing image.

```go
func (c *Client) Variations(ctx context.Context, req *provider.VariationsRequest) (*provider.VariationsResponse, error)
```

Returns `ErrNotSupported` if the provider doesn't support variations.

**Example:**

```go
resp, err := client.Variations(ctx, &provider.VariationsRequest{
    Model: omniimage.ModelDALLE2,
    Image: "data:image/png;base64,...",
    N:     4,
})
```

### Upscale

Upscale an image to a higher resolution.

```go
func (c *Client) Upscale(ctx context.Context, req *provider.UpscaleRequest) (*provider.UpscaleResponse, error)
```

Returns `ErrNotSupported` if the provider doesn't support upscaling.

**Example:**

```go
resp, err := client.Upscale(ctx, &provider.UpscaleRequest{
    Model: omniimage.ModelClarityUpscaler,
    Image: "https://example.com/image.jpg",
    Scale: 2,
})
```

### Capability Checks

```go
// SupportsEdit returns true if the provider supports image editing.
func (c *Client) SupportsEdit() bool

// SupportsVariations returns true if the provider supports variations.
func (c *Client) SupportsVariations() bool

// SupportsUpscale returns true if the provider supports upscaling.
func (c *Client) SupportsUpscale() bool
```

**Example:**

```go
if client.SupportsEdit() {
    resp, err := client.Edit(ctx, editReq)
}
```

### Close

Release resources held by the client.

```go
func (c *Client) Close() error
```

**Example:**

```go
defer client.Close()
```

### Provider

Get the underlying provider.

```go
func (c *Client) Provider() provider.Provider
```

## Helper Functions

### GetCapabilities

Get provider capabilities without creating a client.

```go
func GetCapabilities(name ProviderName) *Capabilities
```

**Capabilities struct:**

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

**Example:**

```go
caps := omniimage.GetCapabilities(omniimage.ProviderNameOpenAI)
if caps.Edit {
    fmt.Println("Provider supports editing")
}
```

### GetModelInfo

Get information about a specific model.

```go
func GetModelInfo(modelID string) *ModelInfo
```

**ModelInfo struct:**

```go
type ModelInfo struct {
    ID           string
    Provider     ProviderName
    Name         string
    MaxImages    int
    SupportsHD   bool
    SupportsEdit bool
}
```

**Example:**

```go
info := omniimage.GetModelInfo(omniimage.ModelGPTImage2)
fmt.Printf("Model: %s, Max Images: %d\n", info.Name, info.MaxImages)
```
